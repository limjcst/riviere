package models

import (
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

// testDBDriver is the database driver used for test
const testDBDriver = "sqlite3"

// testDBSourceName is the database data source name used for test
const testDBSourceName = "file:test.db?cache=shared&mode=memory"

func TestNewDB(t *testing.T) {
	// Unknown driver
	_, err := NewDB("unknown", testDBSourceName)
	if err == nil {
		t.Error("Create db based on unknon driver")
	}

	// data source name out of reach
	_, err = NewDB(testDBDriver, "file:test.db?cache=shared&mode=ro")
	if err == nil {
		t.Error("Create db based on unreachable data source name")
	}

	_, err = NewDB(testDBDriver, testDBSourceName)
	if err != nil {
		t.Error(err)
	}
}

func TestNewTunnel(t *testing.T) {
	db, err := NewDB(testDBDriver, testDBSourceName)
	if err != nil {
		t.Error(err)
	}
	tunnel := &Tunnel{Port: 10000, ForwardAddress: "127.0.0.1", ForwardPort: 80}
	err = db.NewTunnel(tunnel)
	if err != nil {
		t.Error("Failed to insert tunnel into database")
	}
	err = db.NewTunnel(tunnel)
	if err == nil {
		t.Error("Insert duplicated tunnel into database")
	}
	n, err := db.DeleteTunnel(tunnel)
	if n != 1 || err != nil {
		t.Error("Failed to delete tunnel from database")
	}
}

func TestDeleteTunnel(t *testing.T) {
	db, err := NewDB(testDBDriver, testDBSourceName)
	if err != nil {
		t.Error(err)
	}
	tunnel := &Tunnel{Port: 10000, ForwardAddress: "127.0.0.1", ForwardPort: 80}
	err = db.NewTunnel(tunnel)
	if err != nil {
		t.Log(err)
		t.Error("Failed to insert tunnel into database")
	}
	tunnels, err := db.ListTunnel()
	if err != nil || len(tunnels) != 1 {
		t.Error("Failed to fetch tunnels")
	}
	n, err := db.DeleteTunnel(tunnel)
	if n != 1 || err != nil {
		t.Error("Failed to delete tunnel from database")
	}
	tunnels, err = db.ListTunnel()
	if err != nil || len(tunnels) != 0 {
		t.Error("Failed to fetch tunnels")
	}
	n, err = db.DeleteTunnel(tunnel)
	if n != 0 || err != nil {
		t.Error("Failed to delete unexisted tunnel")
	}
}

func TestSqlInject(t *testing.T) {
	db, err := NewDB(testDBDriver, testDBSourceName)
	if err != nil {
		t.Error(err)
	}
	fakeAddress := `127.0.0.1", 80);
    DELETE FROM tunnel WHERE port=10000;
    INSERT INTO tunnel (port, forward_port, forward_address) VALUES(10001, "127.0.0.1`
	tunnel := &Tunnel{Port: 10000, ForwardAddress: fakeAddress, ForwardPort: 80}
	err = db.NewTunnel(tunnel)
	if err != nil {
		t.Error("Failed to insert tunnel into database")
	}
	n, err := db.DeleteTunnel(tunnel)
	if n != 1 || err != nil {
		t.Error("Injected sql is applied")
	}
}
