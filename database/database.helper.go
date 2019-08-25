package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"
	"skripsi-sso/config"
	"skripsi-sso/database/entities"
)

//query for parse time coloum created_at, updated_at, deleted_at not wuth custom table
const QueryTimeFormat = "to_char(created_at, 'YYYY/MM/DD HH24:MI:SS') as created_time," +
	"to_char(updated_at, 'YYYY/MM/DD HH24:MI:SS') as updated_time, " +
	"to_char(deleted_at, 'YYYY/MM/DD HH24:MI:SS') as deleted_time "

//query for parse time coloum created_at, updated_at, deleted_at with custom table / table alias
func QueryTimeFormatWithTable(tableName string) string {
	return "to_char(" + tableName + ".created_at, 'YYYY/MM/DD HH24:MI:SS') as created_time," +
		"to_char(" + tableName + ".updated_at, 'YYYY/MM/DD HH24:MI:SS') as updated_time, " +
		"to_char(" + tableName + ".deleted_at, 'YYYY/MM/DD HH24:MI:SS') as deleted_time "
}

var DB []*gorm.DB

func Init() {
	for i := range config.Databases {
		var database = &config.Databases[i]

		db, err := gorm.Open(database.DriverName, database.ConnectionString)
		if err != nil {
			logrus.Error(err.Error())
			panic("failed to connect database #" + database.Name)
		}

		db.DB().SetMaxOpenConns(database.MaxConnectionOpen)
		db.LogMode(true)
		db.SetLogger(logrus.StandardLogger())

		if config.App.Env != "prod" {
			logrus.WithFields(logrus.Fields{
				"config": database,
			}).Info("Connected to database")
		}

		gormMigration(database.Name, db)
		//append database to array
		DB = append(DB, db)
	}
}

//register entity for created table
func gormMigration(dbName string, db *gorm.DB) {
	db.SingularTable(true)
	db.AutoMigrate(&entities.User{}, &entities.Role{})
}
