package db

import (
	"biblia-be/internal/model"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDB(host, user, password, db_name, db_addr string, maxOpenConns, maxIdleConns, maxIdleTime int) (*gorm.DB, error) {
	addr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, db_addr, db_name)
	log.Println("connecting to the database: " + addr)
	sqlDB, err := sql.Open("mysql", addr)

	if err != nil {
		log.Printf("error connecting db with %s", addr)
		return nil, err
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&model.User{}, &model.Record{})
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxIdleTime(time.Duration(maxIdleTime))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, err
	}

	log.Println("connected to the database successfully!")
	return db, nil
}
