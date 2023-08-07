package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type DbConfig struct {
	PostgresHost     string `json:"postgresHost"`
	PostgresPort     int    `json:"postgresPort"`
	PostgresUser     string `json:"postgresUser"`
	PostgresPassword string `json:"postgresPassword"`
	PostgresDb       string `json:"postgresDb"`
}

type Config struct {
	ClientOrigin    string   `json:"clientOrigin"`
	Host            string   `json:"host"`
	Port            int      `json:"port"`
	Database1       DbConfig `json:"database1"`
	Database2       DbConfig `json:"database2"`
	UseAuth         bool     `json:"useAuth"`
	AuthKey         string   `json:"authKey"`
	FlushDbAfterUse bool     `json:"flushDbAfterUse"`
}

func InitConfig() *Config {
	var confResult Config
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		log.Errorf("Error reading config file, %s", err)
		panic(err)
	}

	err = viper.Unmarshal(&confResult)
	//fmt.Println(config)
	if err != nil {
		log.Errorf("Error unmarshal config file, %s", err)
		panic(err)
	}

	return &confResult
}
