package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestReadConfig(t *testing.T) {
	// Create a temporary config file for testing
	tmpfile, err := os.CreateTemp("", "config.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	config := &Config{
		Directory:    "test_directory",
		SearchPrefix: "test_prefix",
		ReplaceWith:  "test_replace",
	}
	data, err := json.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := tmpfile.Write(data); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test readConfig function
	result, err := readConfig(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	if *result != *config {
		t.Errorf("Expected config %+v, but got %+v", *config, *result)
	}
}

func TestCheckFlagErr(t *testing.T) {
	expectedErr := errors.New("test error")
	err := checkFlagErr(expectedErr)

	if err == nil {
		t.Error("Expected an error, but got nil")
	}

	if err.Error() != "unable to parse arguments: test error" {
		t.Errorf("Expected error message 'unable to parse arguments: test error', but got '%s'", err.Error())
	}

	err = checkFlagErr(nil)
	if err != nil {
		t.Errorf("Expected nil, but got an error: %s", err.Error())
	}
}

func TestRun(t *testing.T) {
	testDir := "test_directory"
	os.Mkdir(testDir, 0755)
	defer os.RemoveAll(testDir)

	// Create a temporary config file for testing
	tmpfile, err := os.CreateTemp("", "config.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	config := &Config{
		Directory:    testDir,
		SearchPrefix: "test_prefix",
		ReplaceWith:  "test_replace",
	}
	data, err := json.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := tmpfile.Write(data); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	err = run([]string{"", "-config", tmpfile.Name()})
	if err != nil {
		t.Errorf("run() returned error: %s", err)
	}
}

func TestRenameFilesWithPrefix(t *testing.T) {
	testDir := "test_directory"
	os.Mkdir(testDir, 0755)
	defer os.RemoveAll(testDir)

	// Create test files
	files := []string{"test_prefix_1.txt", "test_prefix_2.txt", "unrelated_file.txt"}
	for _, file := range files {
		f, err := os.Create(filepath.Join(testDir, file))
		if err != nil {
			t.Fatal(err)
		}
		f.Close()
	}

	config := &Config{
		Directory:    testDir,
		SearchPrefix: "test_prefix",
		ReplaceWith:  "renamed",
	}

	err := renameFilesWithPrefix(config)
	if err != nil {
		t.Fatal(err)
	}

	// Check if the files are renamed correctly
	renamedFiles := []string{"renamed_1.txt", "renamed_2.txt", "unrelated_file.txt"}
	for _, file := range renamedFiles {
		if _, err := os.Stat(filepath.Join(testDir, file)); os.IsNotExist(err) {
			t.Errorf("Expected file '%s' not found", file)
		}
	}
}
