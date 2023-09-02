# digger-config-generator

Automatically generates a [digger](https://github.com/diggerhq/digger) configuration, with some support for [terragrunt](https://terragrunt.gruntwork.io/) dependencies. Run this in CI so you don't need to manually update your digger.yaml!

Supports a basic terraform workflow, where there is one root directory with a structure like the following:

```
/terraform
  /project
    main.tf
    terragrunt.hcl
  /projects_with_common_module
    /_module
      main.tf
    /project
      terragrunt.hcl
  terragrunt.hcl
```

See [tests](/tests) for examples.
