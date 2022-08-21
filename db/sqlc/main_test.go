package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbDriver = "mysql"
	// dbSource = "local:!Elnino1903@tcp(localhost:3306)/simple_bank"

	dbSource = "root:secret@tcp(localhost:3356)/simple_bank"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to dbL:", err)
	}
	testQueries = New(conn)
	log.Output(1, "connect successfully")

	os.Exit(m.Run())
}
