package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/CodeSingerGnC/MicroBank/util"
	_ "github.com/go-sql-driver/mysql"
)

var testQueries *Queries
var testdb *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	testdb, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	
	testQueries = New(testdb)

	os.Exit(m.Run())
}