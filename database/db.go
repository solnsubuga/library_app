package database

import (
	"database/sql"

	"github.com/go-gorp/gorp"
	"github.com/snsubuga/library/models"
)

//DB connection to the Database
var DB *sql.DB

//DbMap gorp connection to the database
var DbMap *gorp.DbMap

//GetDbConnection gets a connection to the database
func GetDbConnection() *sql.DB {
	return DB
}

//GetDbMap returns a dbmap connection object
func GetDbMap() *gorp.DbMap {
	return DbMap
}

//InitDb initializes dbmap
func InitDb() {

	//TODO: make the Dialect engine InnoDB and Encoding UTF-8
	DbMap = &gorp.DbMap{Db: DB, Dialect: gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8"}}

	//Add tables to the database
	DbMap.AddTableWithName(models.Book{}, "books").SetKeys(true, "pk")
	DbMap.AddTableWithName(models.User{}, "users").SetKeys(false, "username")

	DbMap.CreateTablesIfNotExists()
}
