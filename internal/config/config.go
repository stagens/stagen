package config

import (
	"os"
	"path/filepath"
	"time"

	coreConfig "github.com/pixality-inc/golang-core/config"
	"github.com/pixality-inc/golang-core/http"
	"github.com/pixality-inc/golang-core/logger"

	"stagen/pkg/stagen"
)

type Config struct {
	Logger logger.YamlConfig     `env-prefix:"STAGEN_LOG_"  yaml:"logger"`
	Stagen stagen.ConfigYaml     `env-prefix:"STAGEN_"      yaml:"stagen"`
	Site   stagen.SiteConfigYaml `env-prefix:"STAGEN_SITE_" yaml:"site"`
}

func NewConfig(workDir string) *Config {
	return &Config{
		Logger: logger.YamlConfig{},
		Stagen: stagen.ConfigYaml{
			EnvValue: "dev",
			HttpValue: http.ConfigYaml{
				HostValue:            "127.0.0.1",
				PortValue:            8001,
				ShutdownTimeoutValue: 10 * time.Second,
			},
			SettingsValue: stagen.ConfigSettingsYaml{},
			DirsValue: stagen.ConfigDirsYaml{
				WorkValue: workDir,
			},
		},
		Site: stagen.SiteConfigYaml{},
	}
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

func NewConfigFromEnv() (*Config, error) {
	return coreConfig.NewConfigFromEnv[Config]()
}

func NewConfigFromFile(filename string) (*Config, error) {
	if filename == "" {
		filename = configFile()
	}

	return coreConfig.NewConfig[Config](filename)
}

func LoadConfig(filename string) *Config {
	if filename == "" {
		filename = configFile()
	}

	return coreConfig.LoadConfig[Config](filename)
}
