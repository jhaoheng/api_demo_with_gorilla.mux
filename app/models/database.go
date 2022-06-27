package models

import (
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB *gorm.DB
)

// // NewDBPostgres -
// func NewDBPostgres(host, user, password, dbname string) {
// 	if DB != nil {
// 		return
// 	}

// 	var err error
// 	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=5432 sslmode=disable TimeZone=Asia/Taipei", host, user, password, dbname)
// 	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
// 		Logger: logger.Default.LogMode(logger.Info),
// 	})
// 	if err != nil {
// 		logrus.Fatal(err)
// 	}
// 	DB.Debug()
// }

type DBSet struct {
	Host    string
	User    string
	Pass    string
	DBName  string
	IsDebug bool
}

func NewDBMySQL(dbset DBSet) {
	if DB != nil {
		return
	}

	dsn := fmt.Sprintf("%v:%v@tcp(%v:3306)/%v?charset=utf8mb4&parseTime=true", dbset.User, dbset.Pass, dbset.Host, dbset.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logrus.Fatalf(err.Error())
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(30)

	//
	var loggerMode logger.Interface = logger.Default.LogMode(logger.Silent)
	if dbset.IsDebug {
		loggerMode = logger.Default.LogMode(logger.Info)
	}

	DB, err = gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{
		Logger: loggerMode,
	})
	if err != nil {
		logrus.Fatalf(err.Error())
	}
}
