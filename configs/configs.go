package configs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

//Config ...
type Config struct {
	DBDialect  string
	DBHost     string
	DBPort     string
	DBUser     string
	DBName     string
	DBPassword string
	SSLMode    string
	HTTP       string
}

// Configs ..
var Configs *Config

// ReadConfig ...
func ReadConfig(file string) {
	Configs = new(Config)
	rawdata, _ := ioutil.ReadFile("config.json")
	var config map[string]interface{}
	err := json.Unmarshal(rawdata, &config)
	if err != nil {
		panic("Cannot unmarshal the json ")
	}
	Configs.DBDialect = fmt.Sprintf("%v", config["DBDialect"])
	Configs.DBHost = fmt.Sprintf("%v", config["DBHost"])
	Configs.DBPort = fmt.Sprintf("%v", config["DBPort"])
	Configs.DBUser = fmt.Sprintf("%v", config["DBUser"])
	Configs.DBName = fmt.Sprintf("%v", config["DBName"])
	Configs.DBPassword = fmt.Sprintf("%v", config["DBPassword"])
	Configs.SSLMode = fmt.Sprintf("%v", config["SSLMode"])
	Configs.HTTP = fmt.Sprintf("%v", config["HTTP"])
}

// DBconfig ...
func DBconfig() string {
	dbConfig := Configs.DBDialect + "/" +
		" host=" + Configs.DBHost +
		" port=" + Configs.DBPort +
		" user=" + Configs.DBUser +
		" dbname=" + Configs.DBName +
		" password=" + Configs.DBPassword +
		" sslmode=" + Configs.SSLMode
	return dbConfig
}
