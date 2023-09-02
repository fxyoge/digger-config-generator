package main

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/alecthomas/kong"
	"gopkg.in/yaml.v2"
)

type DiggerProject struct {
	Name            string   `yaml:"name"`
	Dir             string   `yaml:"dir"`
	Workflow        string   `yaml:"workflow"`
	Terragrunt      bool     `yaml:"terragrunt"`
	IncludePatterns []string `yaml:"include_patterns"`
	DependsOn       []string `yaml:"depends_on"`
}

type DiggerWorkflowStep struct {
	Run   string `yaml:"run"`
	Shell string `yaml:"shell"`
}

type DiggerWorkflowStage struct {
	Steps []interface{} `yaml:"steps"`
}

type DiggerWorkflowConfiguration struct {
	OnPullRequestClosed []string `yaml:"on_pull_request_closed"`
	OnPullRequestPushed []string `yaml:"on_pull_request_pushed"`
	OnCommitToDefault   []string `yaml:"on_commit_to_default"`
}

type DiggerWorkflow struct {
	Plan                  DiggerWorkflowStage         `yaml:"plan"`
	Apply                 DiggerWorkflowStage         `yaml:"apply"`
	WorkflowConfiguration DiggerWorkflowConfiguration `yaml:"workflow_configuration"`
}

type DiggerConfig struct {
	Projects         []DiggerProject           `yaml:"projects"`
	Workflows        map[string]DiggerWorkflow `yaml:"workflows"`
	CollectUsageData bool                      `yaml:"collect_usage_data"`
	AutoMerge        bool                      `yaml:"auto_merge"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func listDirectories(dir string) ([]string, error) {
	var directories []string

	d, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer d.Close()

	files, err := d.Readdir(-1)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			directories = append(directories, file.Name())
		}
	}

	return directories, nil
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func extractProjectDependencies(path string) ([]string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	regexPattern := `\$\{get_terragrunt_dir\(\)\}(/\.\.)(/\.\.)?(/[\w]+)(/[\w]+)?`
	re := regexp.MustCompile(regexPattern)
	matches := re.FindAllStringSubmatch(string(content), -1)

	var dependencies []string
	for _, match := range matches {
		filteredMatch := make([]string, 0, len(match)-1)
		for _, part := range match[1:] {
			if part != "/.." && part != "" {
				filteredMatch = append(filteredMatch, strings.Replace(part, "/", "", -1))
			}
		}

		dependencies = append(dependencies, strings.Join(filteredMatch, "_"))
	}

	return dependencies, nil
}

func generateDigger(w *bufio.Writer, terraformPath string) error {
	config := DiggerConfig{
		Projects: make([]DiggerProject, 0),
		Workflows: map[string]DiggerWorkflow{
			"prod": {
				Plan: DiggerWorkflowStage{
					Steps: []interface{}{
						"init",
						"plan",
						DiggerWorkflowStep{
							Run:   "terraform fmt -check -diff -recursive",
							Shell: "bash",
						},
					},
				},
				Apply: DiggerWorkflowStage{
					Steps: []interface{}{
						"init",
						DiggerWorkflowStep{
							Run:   "terraform fmt -check -diff -recursive",
							Shell: "bash",
						},
						"apply",
					},
				},
				WorkflowConfiguration: DiggerWorkflowConfiguration{
					OnPullRequestClosed: []string{"digger unlock"},
					OnPullRequestPushed: []string{"digger plan"},
					OnCommitToDefault:   []string{"digger unlock"},
				},
			},
		},
		CollectUsageData: false,
		AutoMerge:        false,
	}

	dirs, err := listDirectories(terraformPath)
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		isTerragrunt, err := pathExists(terraformPath + "/" + dir + "/terragrunt.hcl")
		if err != nil {
			return err
		}

		if isTerragrunt {
			config.Projects = append(config.Projects, DiggerProject{
				Name:       dir,
				Dir:        terraformPath + "/" + dir,
				Workflow:   "prod",
				Terragrunt: true,
				IncludePatterns: []string{
					terraformPath + "/terragrunt.hcl",
					terraformPath + "/" + dir + "/**",
				},
			})
		} else {
			hasModule, err := pathExists(terraformPath + "/" + dir + "/_module")
			if err != nil {
				return err
			}

			subdirs, err := listDirectories(terraformPath + "/" + dir)
			if err != nil {
				return err
			}

			for _, subdir := range subdirs {
				if subdir == "_module" {
					continue
				}

				deps, err := extractProjectDependencies(terraformPath + "/" + dir + "/" + subdir + "/terragrunt.hcl")
				if err != nil {
					return err
				}

				project := DiggerProject{
					Name:            dir + "_" + subdir,
					Dir:             terraformPath + "/" + dir + "/" + subdir,
					Workflow:        "prod",
					Terragrunt:      true,
					IncludePatterns: []string{},
					DependsOn:       deps,
				}

				project.IncludePatterns = append(project.IncludePatterns, terraformPath+"/terragrunt.hcl")
				if hasModule {
					project.IncludePatterns = append(project.IncludePatterns, terraformPath+"/"+dir+"/_module/**")
				}
				project.IncludePatterns = append(project.IncludePatterns, terraformPath+"/"+dir+"/"+subdir+"/**")

				config.Projects = append(config.Projects, project)
			}
		}
	}

	yamlData, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}

	w.Write(yamlData)

	return nil
}

var cli struct {
	InputPath  string `name:"input" type:"path" default:"./terraform" help:"Root path to terraform projects, e.g. './terraform'."`
	OutputPath string `name:"output" type:"path" help:"Path to output digger config, e.g. './digger.yaml'. If not specified, output will go to stdout."`
}

func main() {
	kong.Parse(&cli)

	var w *bufio.Writer
	if cli.OutputPath == "" {
		w = bufio.NewWriter(os.Stdout)
	} else {
		file, err := os.Create(cli.OutputPath)
		check(err)

		defer file.Close()

		w = bufio.NewWriter(file)
	}

	err := generateDigger(w, cli.InputPath)
	check(err)

	w.Flush()
}
