/*
multiserver_converter is a utility designed to convert MT authentication
databases to the multiserver authentication database scheme

Usage:
	multiserver_converter <sqlite3 <in> <out> | psql <in_db> <in_user> <in_password> <in_host> <in_port> <out_db> <out_user> <out_password> <out_host> <out_port>>
*/
package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func convertSQLite3(in, out string) error {
	db, err := sql.Open("sqlite3", in)
	if err != nil {
		return err
	}
	defer db.Close()

	db2, err := sql.Open("sqlite3", out)
	if err != nil {
		return err
	}
	defer db2.Close()

	if _, err := db2.Exec(`CREATE TABLE auth (
	name VARCHAR(32) PRIMARY KEY NOT NULL,
	password VARCHAR(512) NOT NULL
);
CREATE TABLE privileges (
	name VARCHAR(32) PRIMARY KEY NOT NULL,
	privileges VARCHAR(1024)
);
CREATE TABLE ban (
	addr VARCHAR(39) PRIMARY KEY NOT NULL,
	name VARCHAR(32) NOT NULL
);`); err != nil {
		return err
	}

	result := db.QueryRow("SELECT name, password FROM auth;")

	for {
		var name, password string

		if err := result.Scan(&name, &password); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				break
			}

			return err
		}

		_, err = db2.Exec(`INSERT INTO auth (
	name, password
) VALUES (
	?,
	?
);`, name, password)
		if err != nil {
			return err
		}
	}

	return nil
}

func convertPSQL(inDB, inUser, inPassword, inHost string, inPort int, outDB, outUser, outPassword, outHost string, outPort int) error {
	inConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", inHost, inPort, inUser, inPassword, inDB)
	outConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", outHost, outPort, outUser, outPassword, outDB)

	db, err := sql.Open("postgres", inConn)
	if err != nil {
		return err
	}
	defer db.Close()

	db2, err := sql.Open("postgres", outConn)
	if err != nil {
		return err
	}
	defer db2.Close()

	if _, err := db2.Exec(`CREATE TABLE auth (
	name VARCHAR(32) PRIMARY KEY NOT NULL,
	password VARCHAR(512) NOT NULL
);
CREATE TABLE privileges (
	name VARCHAR(32) PRIMARY KEY NOT NULL,
	privileges VARCHAR(1024)
);
CREATE TABLE ban (
	addr VARCHAR(39) PRIMARY KEY NOT NULL,
	name VARCHAR(32) NOT NULL
);`); err != nil {
		return err
	}

	result := db.QueryRow("SELECT name, password FROM auth;")

	for {
		var name, password string

		if err := result.Scan(&name, &password); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				break
			}

			return err
		}

		_, err = db2.Exec(`INSERT INTO auth (
	name, password
) VALUES (
	$1,
	$2
);`, name, password)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: multiserver_converter <sqlite3 <in> <out> | psql <in_db> <in_user> <in_password> <in_host> <in_port> <out_db> <out_user> <out_password> <out_host> <out_port>>")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "sqlite3":
		if len(os.Args) != 4 {
			fmt.Println("Usage: multiserver_converter <sqlite3 <in> <out> | psql <in_db> <in_user> <in_password> <in_host> <in_port> <out_db> <out_user> <out_password> <out_host> <out_port>>")
			os.Exit(1)
		}

		if err := convertSQLite3(os.Args[2], os.Args[3]); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "psql":
		if len(os.Args) != 12 {
			fmt.Println("Usage: multiserver_converter <sqlite3 <in> <out> | psql <in_db> <in_user> <in_password> <in_host> <in_port> <out_db> <out_user> <out_password> <out_host> <out_port>>")
			os.Exit(1)
		}

		inPort, err := strconv.Atoi(os.Args[6])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		outPort, err := strconv.Atoi(os.Args[11])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := convertPSQL(os.Args[2], os.Args[3], os.Args[4], os.Args[5], inPort, os.Args[7], os.Args[8], os.Args[9], os.Args[10], outPort); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	fmt.Println("Database converted successfully")
}
