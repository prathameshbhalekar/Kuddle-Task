package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/KuddleTask/api/schemas"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB = nil

func GetDB() *gorm.DB {
	if db != nil {
		return db
	}

	db = Connect()
	return db
}

func Connect() *gorm.DB {
	username := viper.GetString("DB_USER")
	password := viper.GetString("DB_PASS")
	dbName := viper.GetString("DB_NAME")
	dbHost := viper.GetString("DB_HOST")
	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s port=5432 password=%s sslmode=require", dbHost, username, dbName, password)
	sqlDB, err := sql.Open("postgres", dbUri)
	if err != nil {
		log.Fatal(err)
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(50)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB}), &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Info),
		PrepareStmt: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	db.Migrator().CreateTable(&schemas.Class{})
	db.Migrator().CreateTable(&schemas.Registration{})

	return db
}
