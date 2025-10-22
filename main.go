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

/*
EnvVar Structure for represeting environment variable entry
*/
type EnvVar struct {
	Name      string `yaml:"name"`
	Value     string `yaml:"value"`
	Type      string `yaml:"type,omitempty"`      // *string*, array
	Separator string `yaml:"separator,omitempty"` // Default: ","
	Action    string `yaml:"action,omitempty"`    // *replace*, merge, new, unset
}

/*
EnvVarList - Each environment is in fact list of EnvVar's
*/
type EnvVarList []EnvVar

/*
EnvVarConfig - map of EnvVarList's represeting each individual "environment"
*/
type EnvVarConfig map[string]EnvVarList

/*
NoConfigError - custom error representing "No Config Found event"
*/
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
	debug          bool
)

const defaultConfigFile = ".runctl.yaml"

func init() {
	var selectedConfigFile string
	var selectedConfig string
	defaultShell, shellSet := os.LookupEnv("SHELL")
	if !shellSet {
		defaultShell = "/bin/sh"
	}
	envConfigFile, envConfigFileSet := os.LookupEnv("RUNCTL_CONFIG")
	if envConfigFileSet {
		selectedConfigFile = envConfigFile
	} else {
		selectedConfigFile = defaultConfigFile
	}

	envConfig, envConfigSet := os.LookupEnv("RUNCTL_ENV")
	if envConfigSet {
		selectedConfig = envConfig
	} else {
		selectedConfig = ""
	}

	flag.StringVar(&shell, "shell", defaultShell, "shell to use for command interpretation")
	flag.StringVar(&logFileName, "log", "", "log file name")
	flag.StringVar(&configFileName, "config-file", selectedConfigFile, "Environment list file (or $RUNCTL_CONFIG)")
	flag.StringVar(&configNames, "environment", selectedConfig, "Environment name (or $RUNCTL_ENV)")
	flag.BoolVar(&debug, "debug", false, "Print debug info")

	origUsage := flag.Usage
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, os.Args[0]+" [options] [--] command...\n\n")
		origUsage()
	}
}

func compileEnv(envVarList EnvVarList) []string {
	var eList []string
	for _, myVar := range envVarList {
		if myVar.Type == "" {
			if myVar.Action == "merge" {
				myVar.Type = "array"
			} else {
				myVar.Type = "string"
			}
		}
		if myVar.Action == "" {
			myVar.Action = "replace"
		}
		if myVar.Separator == "" {
			myVar.Separator = ","
		}
		envVar, defined := os.LookupEnv(myVar.Name)
		if debug && defined {
			fmt.Println("Found " + envVar + " defined as \"" + envVar + "\" action: " + myVar.Action)
		}
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
				if debug {
					fmt.Println("Redefined " + envVar + " as \"" + myVar.Value + "\"")
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
	var commandArgs []string

	configPaths := []string{configFileName, path.Join(homeDir, defaultConfigFile), path.Join("/etc/", defaultConfigFile)}
	configFile, err := openConfig(configPaths)
	if err != nil {
		fmt.Println(err)
	}
	defer configFile.Close()
	configBytes, _ := ioutil.ReadAll(configFile)
	yaml.Unmarshal(configBytes, &envConfig)

	commandArgs = flag.Args()

	if len(commandArgs) < 1 {
		fmt.Fprintln(os.Stderr, "No command specified")
		return
	}
	shellCommand = strings.Join(commandArgs, " ")
	cmd := exec.Command(shell, "-c", shellCommand)

	configNameList := strings.Split(configNames, ",")
	for _, configName := range configNameList {
		envList = append(envList, compileEnv(envConfig[configName])...)
	}
	if debug {
		fmt.Fprintln(os.Stderr, "Command executed: ", shellCommand)
		fmt.Fprintln(os.Stderr, envList)
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
