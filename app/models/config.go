package models

import (
	"github.com/spf13/viper"
	"fmt"
	"log"
)

// config management

type dbConfig struct {
	Host    string
	User    string
	Password	string
	Database string
}

func GetDbConfigs(env string) (db dbConfig){
	viper.SetConfigName("db_config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./configs")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	conf := viper.Sub(env)
	var C dbConfig

	uErr := conf.Unmarshal(&C)
	if uErr != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	fmt.Println(C.Host)
	return C
}
