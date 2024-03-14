package controllers

import (
	"database/sql"
	"fmt"

	// "fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func connect() *sql.DB {
	// dbHost := os.Getenv("DB_HOST")
	// fmt.Println(dbHost)
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/db_latihan_pbp?parseTime=true&loc=Asia%2FJakarta")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func connectGorm() *gorm.DB {
	dbHost := os.Getenv("DB_HOST")
	fmt.Println(dbHost)
	db, err := gorm.Open(mysql.Open("root:@tcp(localhost:3306)/db_latihan_pbp"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return db
}
