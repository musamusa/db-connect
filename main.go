package db_connect

import (
	"os"
	"fmt"
	"encoding/json"
	"github.com/go-xorm/xorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Configuration struct {
	DBType string
	DBUser string
	DBPass string
	DBHost string
	DBPort string
	DBName string
	DBPath string
}

func GetConfig () Configuration {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println(`
Provide JSON data in a file config.json containing the following:
 {
	"dbType": "sqlite",
	"dbUser": "admin",
	"dbPass": "admin",
	"dbHost": "localhost",
	"dbPort": "5432",
	"dbName": "cog-db",
	"dbPath": ""
 }
defaulting to sqlite
        `)
		configuration.DBType = "sqlite"
		configuration.DBPath = "/tmp/system.db"
	}
	return configuration
}

func GetConnection () (*xorm.Engine, error) {
	config := GetConfig()
	return xorm.NewEngine(config.DBType, config.DBPath)
}

func SyncDB (bean interface{}) (*xorm.Engine, error) {
	orm, err := GetConnection()
	err = orm.CreateTables(bean)
	if err != nil {
		return orm, err
	}
	err = orm.Sync2(bean)

	if err != nil {
		return orm, err
	}

	return orm, err
}