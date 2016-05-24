package server

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	dbProtocol string
	dbAddress  string
	dbName     string
	dbUsername string
	dbPassword string
}

func NewConfig(configPath string) *Config {
	viper.SetConfigName("config")
	viper.AddConfigPath(configPath)

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	config := &Config{
		dbProtocol: viper.GetString("database.protocol"),
		dbAddress:  viper.GetString("database.address"),
		dbName:     viper.GetString("database.name"),
		dbUsername: viper.GetString("database.username"),
		dbPassword: viper.GetString("database.password"),
	}

	if config.dbName == "" {
		panic("Fatal error: dbName not specified")
	}

	return config
}

func (config *Config) GetDataSourceName() string {

	// https://github.com/go-sql-driver/mysql#dsn-data-source-name
	// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]

	dsn := "/" + config.dbName + "?parseTime=true"

	if config.dbAddress != "" {
		dsn = config.dbProtocol + "(" + config.dbAddress + ")" + dsn
	}

	var dsnUser string

	if config.dbUsername != "" {
		dsnUser = config.dbUsername
		if config.dbPassword != "" {
			dsnUser += ":" + config.dbPassword
		}
	}

	if dsnUser != "" {
		dsn = dsnUser + "@" + dsn
	}

	return dsn
}
