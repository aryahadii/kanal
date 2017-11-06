package configuration

import (
	"github.com/spf13/viper"
)

var (
	// KanalConfigPath is path of config file
	KanalConfigPath = "config.yaml"

	// KanalConfig is config of project
	KanalConfig *viper.Viper
)

// LoadConfig loads Kanal's config file from KanalConfigPath
func LoadConfig() error {
	KanalConfig = viper.New()
	KanalConfig.SetConfigFile(KanalConfigPath)
	if err := KanalConfig.ReadInConfig(); err != nil {
		return err
	}

	KanalConfig.SetDefault("debug", true)

	return nil
}
