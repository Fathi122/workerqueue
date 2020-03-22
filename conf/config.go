package conf

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type (
	// Config parameters
	Config struct {
		Parameters struct {
			Grpc struct {
				Port string
				Host string
			}
			Etcd struct {
				Host string
			}
		}
	}
)

// ServerConfig global config
var serverConfig Config

// GetConfig get configuration
func GetConfig() *Config {
	return &serverConfig
}

func init() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name

	viper.AddConfigPath("../conf") // path to look for the config file in
	viper.AddConfigPath("./conf")  // call multiple times to add many search paths

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Errorln("Fatal error config file:  ", err.Error())
	} else {
		err := viper.Unmarshal(&serverConfig)
		if err != nil {
			log.Debugln("Unable to decode into struct ", err.Error())
		}
	}
}
