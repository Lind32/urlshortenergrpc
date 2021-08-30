package config

import (
	"fmt"

	"github.com/spf13/viper"
)

//Config  параметры среды
type Config struct {
	Grpcadr string
	Httpadr string
	Data    struct {
		PgHost  string
		PgPort  string
		PgUser  string
		PgPass  string
		PgName  string
		PgClean bool
	}
}

//GetConfig чтение конфиг файла
func GetConfig() *Config {

	viper.SetConfigName("config")
	viper.AddConfigPath("$GOPATH/urlshortenergrpc")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w", err))
	}
	conf := &Config{}
	conf.Grpcadr = viper.GetString("grpcadress")
	conf.Httpadr = viper.GetString("httpadress")
	conf.Data.PgHost = viper.GetString("data.host")
	conf.Data.PgPort = viper.GetString("data.port")
	conf.Data.PgUser = viper.GetString("data.user")
	conf.Data.PgPass = viper.GetString("data.pass")
	conf.Data.PgName = viper.GetString("data.name")
	conf.Data.PgClean = viper.GetBool("data.clean")

	return conf

}
