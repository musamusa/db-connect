package db_connect

import (
	"os"
	"fmt"
	"encoding/json"
	"github.com/go-xorm/xorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

type Configuration struct {
	DBType string
	DBUser string
	DBPass string
	DBHost string
	DBPort string
	DBName string
	DBPath string
	DBConnectionString string
}

func GetConfig () *Configuration {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	configuration := &Configuration{}
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
		defaultToSqlite(configuration)
	}
	return configuration
}

func GetConnection () (*xorm.Engine, error) {
	config := updateDBSettings()
	return xorm.NewEngine(config.DBType, config.DBConnectionString)
}

func SyncDB (bean interface{}) (*xorm.Engine, error) {
	orm, err := GetConnection()

	if err != nil {
		return orm, err
	}
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

func setDBString (config *Configuration) {
	config.DBConnectionString =
		config.DBUser + ":" + config.DBPass +
		"@tcp(" + config.DBHost + ":" + config.DBPort + ")/" + config.DBName
}

func updateDBSettings () *Configuration {
	config := GetConfig()
	switch config.DBType {
	case
	"postgres",
	"mysql":
		setDBString(config)
		return config
	case
	"sqlite",
	"sqlite3":
		// overwrite input type
		config.DBType = "sqlite3"
		if len(config.DBPath) > 0 {
			config.DBConnectionString = config.DBPath
		}
		return config
	}
	return config
}

func defaultToSqlite (config *Configuration) {
	if len(config.DBType) == 0 {
		config.DBType = "sqlite3"
		config.DBPath = "/tmp/system.db"
	}
}
