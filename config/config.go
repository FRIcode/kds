package config

import (
	"flag"
	"os"

	"github.com/goccy/go-yaml"
)

type ConfigImpl struct {
	Server      ServerConfig       `yaml:"server"`
	Logging     LoggingConfig      `yaml:"logging"`
	Deployments []DeploymentConfig `yaml:"deployments"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
}

type LoggingConfig struct {
	Level string `yaml:"level"`
}

type DeploymentConfig struct {
	Name    string                `yaml:"name"`
	WorkDir string                `yaml:"workdir"`
	Auth    AuthConfig            `yaml:"auth"`
	Env     []EnvDeploymentConfig `yaml:"env"`
	Script  []string              `yaml:"script"`
}

type AuthConfig struct {
	Type string                 `yaml:"type"`
	Opts map[string]interface{} `yaml:"opts"`
}

type EnvDeploymentConfig struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

var Config ConfigImpl

func InitConfig() {
	configFlag := flag.String("config", "config.yml", "Path to the configuration file")
	flag.Parse()

	configPath := getenvordefault("CONFIG", *configFlag)

	f, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(&Config); err != nil {
		panic(err)
	}
}

func getenvordefault(key, def string) string {
	value := os.Getenv(key)
	if value == "" {
		return def
	}
	return value
}
