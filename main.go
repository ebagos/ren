package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Directory    string `json:"directory"`
	SearchPrefix string `json:"search_prefix"`
	ReplaceWith  string `json:"replace_with"`
}

func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	configFilename := flag.String("config", "config.json", "path to the configuration file")

	err := checkFlagErr(flag.CommandLine.Parse(args[1:]))
	if err != nil {
		return fmt.Errorf("error parsing flags: %w", err)
	}

	config, err := readConfig(*configFilename)
	if err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	err = renameFilesWithPrefix(config)
	if err != nil {
		return fmt.Errorf("error renaming files: %w", err)
	}

	return nil
}

func readConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func renameFilesWithPrefix(config *Config) error {
	return filepath.Walk(config.Directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasPrefix(info.Name(), config.SearchPrefix) {
			oldPath := filepath.Join(config.Directory, info.Name())

			// filepath.Join() を使用してファイル名を作成する
			newPath := filepath.Join(config.Directory, strings.Replace(info.Name(),
				config.SearchPrefix, config.ReplaceWith, 1))

			err := os.Rename(oldPath, newPath)
			if err != nil {
				return err
			}

			// 標準エラー出力への出力に変更
			_, _ = io.WriteString(os.Stderr, filepath.Base(newPath)+"\n")
		}
		return nil
	})
}

func checkFlagErr(err error) error {
	if err != nil {
		return fmt.Errorf("unable to parse arguments: %w", err)
	}
	return nil
}
