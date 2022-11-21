package main

/*
I have used, two ways to connect to db in my work, will be discussing below 2 methods
1. connecting to maria client using Gorm ( in below code block)
2. Connecting to mongo client  (under mongo-clinet module)
*/

import (
	"fmt"
	"time"

	mysql "go.elastic.co/apm/module/apmgormv2/v2/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	DBCon         string
	DBTablePrefix string
)

// PostgresDB : Variable to hold the postgresql db connection
var MariaDB *gorm.DB

// InitializeDB : Initializes the Database connections
func InitializeDB() {

	db, err := gorm.Open(mysql.Open(DBCon+fmt.Sprintf("?%s", "&parseTime=True")), &gorm.Config{
		SkipDefaultTransaction: false,
		PrepareStmt:            true,
		NamingStrategy:         schema.NamingStrategy{SingularTable: true, TablePrefix: DBTablePrefix + "."},
	})
	if err != nil {
		panic(err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(time.Minute * 5)

	MariaDB = db
}

func main() {
	fmt.Println("Hello world")
	InitializeDB()
	// Since the Connectin string is not given, it won't connect for now
	fmt.Printf("DB connected succesfully")
}
