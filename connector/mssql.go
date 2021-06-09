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
	cs = fmt.Sprintf("server=%s;user id=%s;port=%d;database=%s", hostname, username, port, databaseName)
	if !windowsAuthentication {
		cs = fmt.Sprintf("%s;password=%s", cs, password)
	} else {
		cs = fmt.Sprintf("%s;Integrated Security=sspi", cs)
		currentUser, err := user.Current()
		if err != nil {
			currentUser = &user.User{
				Username: "Unknown",
			}
		}
		o.Log.Infof("Attempting to connect to %s:%d (%s) using Windows Authentication. We are running as %+v", hostname, port, databaseName, currentUser)
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
