package main

import (
	"bufio"
	"bytes"
	"os"
	"testing"

	"gotest.tools/v3/assert"
)

func assertSnapshotEquals(t *testing.T, testName string) {
	var buffer bytes.Buffer
	w := bufio.NewWriter(&buffer)

	terraformPath := "tests/" + testName + "/tf"
	err := generateDigger(w, terraformPath)
	if err != nil {
		t.Fatalf("[%s] Failed to generate the digger config in test: %v", testName, err)
	}

	w.Flush()

	snapshotPath := "tests/" + testName + "/digger.yaml"
	snapshot, err := os.ReadFile(snapshotPath)
	if err != nil {
		t.Fatalf("[%s] Failed to read the digger config snapshot in test: %v", testName, err)
	}

	assert.Equal(t, string(snapshot), buffer.String(), "[%s] Generated digger config doesn't match the snapshot (-got +want)", testName)
}

func TestSnapshots(t *testing.T) {
	dirs, err := listDirectories("tests")
	if err != nil {
		t.Fatalf("Failed to list directories in the tests directory: %v", err)
	}

	if len(dirs) < 0 {
		t.Fatalf("No tests found")
	}

	for _, dir := range dirs {
		assertSnapshotEquals(t, dir)
	}
}
