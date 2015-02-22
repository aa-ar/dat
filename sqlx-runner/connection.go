package runner

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/mgutz/dat"
)

// Connection is a connection to the database with an EventReceiver
type Connection struct {
	DB *sqlx.DB
	*Queryable
}

var standardConformingString string

// pgMustNotAllowEscapeSequence checks if Postgres prohibits using escaped
// sequences in strings when dat.EnableInterpolation == true. If escape
// sequences are allowed, then it is not safe to use interpoaltion, and
// this function panics.
func pgMustNotAllowEscapeSequence(conn *Connection) {
	dat.EnableInterpolation = true
	if !dat.EnableInterpolation {
		return
	}

	if standardConformingString == "" {
		err := conn.
			SQL("select setting from pg_settings where name='standard_conforming_strings'").
			QueryScalar(&standardConformingString)
		if err != nil {
			panic(err)
		}
	}

	if standardConformingString != "on" {
		log.Fatalf("Database allows escape sequences. Cannot be used with interpolation. "+
			"standard_conforming_strings=%q\n"+
			"See http://www.postgresql.org/docs/9.3/interactive/sql-syntax-lexical.html#SQL-SYNTAX-STRINGS-ESCAPE",
			standardConformingString)
	}
}

// NewConnection instantiates a Connection for a given database/sql connection
func NewConnection(db *sql.DB, driverName string) *Connection {
	DB := sqlx.NewDb(db, driverName)
	conn := &Connection{DB, &Queryable{DB}}
	if driverName == "postgres" {
		pgMustNotAllowEscapeSequence(conn)
	} else {
		panic("Unsupported driver: " + driverName)
	}
	return conn
}
