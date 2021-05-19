package connector

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/sirupsen/logrus"
)

type MSSQLConnector struct {
	Log *logrus.Logger
}

func (o *MSSQLConnector) GetDB(hostname string, port int, username string, password string, databaseName string) (*sql.DB, error) {

	var db *sql.DB
	var err error

	var cs string
	cs = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s", hostname, username, password, port, databaseName)

	db, err = sql.Open("mssql", cs)
	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}
	return db, nil

}
