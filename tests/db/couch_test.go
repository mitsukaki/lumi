package db_test

import (
	"testing"

	"github.com/mitsukaki/lumi/internal/db"
)

func TestCreateDatabase(t *testing.T) {
	db := db.CreateDatabase("http://localhost:5984")
	if db == nil {
		t.Error("Expected a database, got nil")
	}
}

func TestCouchDatabase_CreateTable(t *testing.T) {
	db := db.CreateDatabase("http://localhost:5984")
	table := db.CreateTable("test")
	if table == nil {
		t.Error("Expected a table, got nil")
	}
}

