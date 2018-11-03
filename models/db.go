package models

import (
	"database/sql"
)

// Database wraps all the db operations
type Database interface {
	NewTunnel(tunnel *Tunnel) error
	DeleteTunnel(tunnel *Tunnel) (int64, error)
	ListTunnel() ([]*Tunnel, error)
}

// DB wraps sql.DB, with a set of implementations of Database
type DB struct {
	*sql.DB
}

// NewDB creates a DB
func NewDB(driverName string, dataSourceName string) (*DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	wdb := &DB{db}
	err = PrepareTunnel(wdb)
	return wdb, err
}
