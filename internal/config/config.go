package config

import (
	"os"
	"path/filepath"
	"runtime"

	coreConfig "github.com/pixality-inc/golang-core/config"
	"github.com/pixality-inc/golang-core/logger"
)

type Config struct {
	Logger logger.YamlConfig `env-prefix:"STAGEN_LOG_" yaml:"logger"`
}

func RootDir() string {
	var (
		_, b, _, _ = runtime.Caller(0)
		basepath   = filepath.Join(filepath.Dir(b), "../..")
	)

	return basepath
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
