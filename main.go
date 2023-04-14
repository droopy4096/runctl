package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"gopkg.in/yaml.v3"
)

type EnvVar struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
	Type  string `yaml:"type,omitempty"`
}

type EnvVarList []EnvVar

type EnvVarConfig map[string]EnvVarList

var (
	logFileName    string
	shell          string
	shellCommand   string
	configFileName string
	configName     string
)

func init() {
	defaultShell, shellSet := os.LookupEnv("SHELL")
	if !shellSet {
		defaultShell = "/bin/sh"
	}

	flag.StringVar(&shellCommand, "command", "ls /", "comand line to run (one string)")
	flag.StringVar(&shell, "shell", defaultShell, "shell to use for command interpretation")
	flag.StringVar(&logFileName, "log", "", "log file name")
	flag.StringVar(&configFileName, "config-file", ".envctl.yaml", "Environment list file")
	flag.StringVar(&configName, "config", ".envctl.yaml", "Environment list file")
}

func compileEnv(envVarList EnvVarList) []string {
	var eList []string
	for _, myVar := range envVarList {
		eList = append(eList, fmt.Sprintf("%s=%s", myVar.Name, myVar.Value))
	}
	return eList
}

func main() {
	flag.Parse()
	var envConfig EnvVarConfig
	var stdout, stderr bytes.Buffer
	var err error

	configFile, err := os.Open(configFileName)
	if err != nil {
		fmt.Println(err)
	}
	defer configFile.Close()
	configBytes, _ := ioutil.ReadAll(configFile)
	yaml.Unmarshal(configBytes, &envConfig)
	cmd := exec.Command(shell, "-c", shellCommand)
	cmd.Env = append(os.Environ(), compileEnv(envConfig[configName])...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err = cmd.Run(); err != nil {
		fmt.Println("Error occured", err)
		return
	}
	fmt.Println(stdout.String())
	fmt.Fprintf(os.Stderr, "%s", stderr.String())
}
