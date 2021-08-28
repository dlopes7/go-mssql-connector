package connector

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/sirupsen/logrus"
	"os/user"
)

type MSSQLConnector struct {
	Log *logrus.Logger
}

func (o *MSSQLConnector) GetDB(hostname string, port int, username string, password string, databaseName string, windowsAuthentication bool) (*sql.DB, error) {

	var db *sql.DB
	var err error

	var cs string
	cs = fmt.Sprintf("server=%s;port=%d;database=%s", hostname, port, databaseName)
	if !windowsAuthentication {
		cs = fmt.Sprintf("%s;user id=%s;password=%s", cs, username, password)
	} else {
		cs = fmt.Sprintf("%s;trusted_connection=yes", cs)
		currentUser, err := user.Current()
		if err != nil {
			currentUser = &user.User{
				Username: username,
			}
		}
		o.Log.Infof("Attempting to connect to %s:%d (%s) using Windows Authentication. We are running as %s", hostname, port, databaseName, currentUser.Username)
	}

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
