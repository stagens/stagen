package config

import (
	"os"
	"path/filepath"

	coreConfig "github.com/pixality-inc/golang-core/config"
	"github.com/pixality-inc/golang-core/logger"

	"stagen/pkg/stagen"
)

type Config struct {
	Logger logger.YamlConfig     `env-prefix:"STAGEN_LOG_"  yaml:"logger"`
	Stagen stagen.ConfigYaml     `env-prefix:"STAGEN_"      yaml:"stagen"`
	Site   stagen.SiteConfigYaml `env-prefix:"STAGEN_SITE_" yaml:"site"`
}

func RootDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return cwd
}

func configFile() string {
	configFilename := os.Getenv("STAGEN_CONFIG_FILE")
	if configFilename == "" {
		configFilename = filepath.Join(RootDir(), "config.yaml")
	}

	return configFilename
}

func LoadConfig() *Config {
	return coreConfig.LoadConfig[Config](configFile())
}
