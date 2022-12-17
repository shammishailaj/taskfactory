package utils

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

type RDBMSConfig struct {
	DriverName      string // Name of the SQL driver viz., mysql, postgres, sqlite etc
	HostName        string // Hostname of the SQL server
	Port            int64  // Port number of the SQL server
	Database        string // Name of the database to connect to on the SQL server
	UserName        string // Username part of the SQL server credentials
	Password        string // Password for the afore-defined UserName
	MaxOpenConns    int    // Maximum number of connections to be opened
	MaxIdleConns    int    // Maximum number of idle connections to be kept
	MaxLifeTimeConn int    // Maximum life of connection (in seconds)
}

func (c *RDBMSConfig) Values(driverName, hostName string, port int64, database, restrictedUsers, userName, password string, maxOpenconns, maxIdleConns, maxLifeTimeConn int) {
	c.DriverName = driverName
	c.HostName = hostName
	c.Port = port
	c.Database = database
	c.UserName = userName
	c.Password = password
	c.MaxOpenConns = maxOpenconns
	c.MaxIdleConns = maxIdleConns
	c.MaxLifeTimeConn = maxLifeTimeConn
}

func (c *RDBMSConfig) DSN() string {
	return c.UserName + ":" + c.Password + "@tcp(" + c.HostName + ":" + strconv.FormatInt(c.Port, 10) + ")/" + c.Database
}

func (c *RDBMSConfig) String() string {
	return fmt.Sprintf("DBConfig = %#v", c)
}

func (c *RDBMSConfig) Connect() (*sql.DB, error) {
	db, err := sql.Open(c.DriverName, c.DSN())
	if err == nil {
		db.SetMaxOpenConns(c.MaxOpenConns)
		db.SetMaxIdleConns(c.MaxIdleConns)
		db.SetConnMaxLifetime(time.Duration(c.MaxLifeTimeConn) * time.Second)
		log.Printf("Successfully connected to: %s", c.DSN())
	}
	return db, err
}
