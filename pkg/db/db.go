package db

import (
	"database/sql"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"io"
	"time"
)

type Connection struct {
	io.Closer
	File string
	Cfg  gorm.Config
	db   *gorm.DB
}

func (c *Connection) Open() error {
	db, openErr := gorm.Open(sqlite.Open(c.File+"?_pragma=foreign_keys(1)"), &c.Cfg)
	if openErr != nil {
		return openErr
	}
	c.db = db
	sqlDB, sqlErr := c.getRealConnection()
	if sqlErr != nil {
		return sqlErr
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(30)
	sqlDB.SetConnMaxIdleTime(10 * time.Second)
	sqlDB.SetConnMaxLifetime(60 * time.Second)
	return nil
}

func (c *Connection) IsConnected() bool {
	return c.db != nil
}

func (c Connection) getRealConnection() (*sql.DB, error) {
	return c.db.DB()
}

func (c *Connection) API() *gorm.DB {
	return c.db
}

func (c *Connection) Close() error {
	conn, err := c.getRealConnection()
	if err != nil {
		return err
	}
	return conn.Close()
}

func NewConnection(file string) *Connection {
	return &Connection{
		File: file,
		Cfg:  gorm.Config{},
	}
}
