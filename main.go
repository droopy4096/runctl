package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

type EnvVar struct {
	Name      string `yaml:"name"`
	Value     string `yaml:"value"`
	Type      string `yaml:"type,omitempty"`      // *string*, array
	Separator string `yaml:"separator,omitempty"` // Default: ","
	Action    string `yaml:"action,omitempty"`    // *replace*, merge, new
}

type EnvVarList []EnvVar

type EnvVarConfig map[string]EnvVarList

type NoConfigError struct {
	SearchPath []string
}

func (*NoConfigError) Error() string {
	return "No Config found in "
}

var (
	logFileName    string
	shell          string
	shellCommand   string
	configFileName string
	configNames    string
)

const defaultConfigFile = ".envctl.yaml"

func init() {
	defaultShell, shellSet := os.LookupEnv("SHELL")
	if !shellSet {
		defaultShell = "/bin/sh"
	}

	flag.StringVar(&shellCommand, "command", "ls /", "comand line to run (one string)")
	flag.StringVar(&shell, "shell", defaultShell, "shell to use for command interpretation")
	flag.StringVar(&logFileName, "log", "", "log file name")
	flag.StringVar(&configFileName, "config-file", defaultConfigFile, "Environment list file")
	flag.StringVar(&configNames, "config", "", "Environment name")
}

func compileEnv(envVarList EnvVarList) []string {
	var eList []string
	for _, myVar := range envVarList {
		if myVar.Type == "" {
			myVar.Type = "string"
		}
		if myVar.Action == "" {
			myVar.Action = "replace"
		}
		if myVar.Separator == "" {
			myVar.Separator = ","
		}
		envVar, defined := os.LookupEnv(myVar.Name)
		if myVar.Action == "new" {
			if defined {
				continue
			}
		}
		if myVar.Action == "unset" {
			os.Unsetenv(myVar.Name)
			continue
		}
		if myVar.Type == "array" {
			if myVar.Action == "merge" {
				if defined {
					elements := strings.Split(envVar, myVar.Separator)
					myVar.Value = strings.Join(append(elements, myVar.Value), myVar.Separator)
				}
				os.Setenv(myVar.Name, myVar.Value)
			}
		}
		eList = append(eList, fmt.Sprintf("%s=%s", myVar.Name, myVar.Value))
	}
	return eList
}

func openConfig(configList []string) (*os.File, error) {
	for _, fileName := range configList {
		configFile, err := os.Open(fileName)
		if err == nil {
			return configFile, err
		}
	}
	return nil, &NoConfigError{
		SearchPath: configList,
	}
}

func main() {
	homeDir, _ := os.LookupEnv("HOME")

	flag.Parse()
	var envConfig EnvVarConfig
	var stdout, stderr bytes.Buffer
	var err error
	var envList []string

	configPaths := []string{configFileName, path.Join(homeDir, defaultConfigFile), path.Join("/etc/", defaultConfigFile)}
	// configFile, err := os.Open(configFileName)
	configFile, err := openConfig(configPaths)
	if err != nil {
		fmt.Println(err)
	}
	defer configFile.Close()
	configBytes, _ := ioutil.ReadAll(configFile)
	yaml.Unmarshal(configBytes, &envConfig)
	cmd := exec.Command(shell, "-c", shellCommand)
	configNameList := strings.Split(configNames, ",")
	for _, configName := range configNameList {
		envList = append(envList, compileEnv(envConfig[configName])...)
	}
	cmd.Env = append(os.Environ(), envList...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err = cmd.Run(); err != nil {
		fmt.Println("Error occured", err)
		return
	}
	fmt.Println(stdout.String())
	fmt.Fprintf(os.Stderr, "%s", stderr.String())
}
