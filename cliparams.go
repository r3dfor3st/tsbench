package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

type cliParamsType struct {
	sql      string
	workers  uint
	params   string
	host     string
	port     uint
	database string
	user     string
	password string
}

const (
	defaultSQL = `
	SELECT time_bucket('1 minutes', ts) bucket, COUNT(*) cnt, MAX(usage) max, MIN(usage) min
	FROM cpu_usage
	WHERE host = $1 AND ts >= $2 AND ts <= $3
	GROUP BY bucket
	ORDER BY bucket, max DESC;`
	pingSQL = "SELECT 'DBD::Pg ping test'"
)

// initCliParams returns a structure containing all relevant commandline arguments.
func initCliParams() (clip cliParamsType) {
	flag.StringVar(&clip.sql, "sql", "", "SQL query to execute (uses provided params)")
	flag.UintVar(&clip.workers, "workers", 1, "number of concurrent workers")
	flag.StringVar(&clip.params, "params", "", "file containing parameters for the query (default uses input stream)")
	flag.StringVar(&clip.host, "host", "localhost", "database server host or socket directory")
	flag.UintVar(&clip.port, "port", 5432, "database server port number")
	flag.StringVar(&clip.database, "database", "homework", "database name")
	flag.StringVar(&clip.user, "user", "postgres", "connect as specified database user")
	flag.StringVar(&clip.password, "password", "", "connect using a specified password (default none)")
	flag.Parse()

	if len(clip.sql) == 0 {
		clip.sql = defaultSQL
	}

	return
}

// connectionString returns a database connecton string constructed from commandline arguments.
func (clip *cliParamsType) connectionString() string {
	var conn strings.Builder
	fmt.Fprintf(&conn, "user=%s dbname=%s host=%s port=%d", clip.user, clip.database, clip.host, clip.port)
	if len(clip.password) > 0 {
		fmt.Fprintf(&conn, " password=%s", clip.password)
	}
	return conn.String()
}

// assert does some common sense verification of provided params e.g. it checks if the file with
// input parameters exists, restricts the total number of workers to some reasonable number and
// even tries to establish a quick connection to the database before proceeding.
func (clip *cliParamsType) assert() (err error) {
	if len(clip.params) > 0 {
		if _, err = os.Stat(clip.params); os.IsNotExist(err) {
			return
		}
	}

	if clip.workers > 200 {
		err = errors.New("number of workers is limited to 200")
		return
	}

	db, err := sql.Open("postgres", clip.connectionString())
	if err != nil {
		return
	}
	defer db.Close()

	var rows *sql.Rows
	rows, err = db.Query(pingSQL)
	if err != nil {
		return
	}
	rows.Close()
	return
}
