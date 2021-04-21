/*
multiserver_converter is a utility designed to convert MT authentication
databases to the multiserver authentication database scheme

Usage:
	multiserver_converter <in> <out>
where in is the path to the MT auth database
and out is the desired path to the newly created database
*/
package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func convertDB(in, out string) error {
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

	if _, err := db2.Exec(`
	CREATE TABLE auth (
		name VARCHAR(32) NOT NULL,
		password VARCHAR(512) NOT NULL
	);
	CREATE TABLE privileges (
		name VARCHAR(32) NOT NULL,
		privileges VARCHAR(1024)
	);
	CREATE TABLE ban (
		addr VARCHAR(39) NOT NULL,
		name VARCHAR(32) NOT NULL
	);`); err != nil {
		return err
	}

	result := db.QueryRow("SELECT name, password FROM auth;")

	for {
		var name, password string

		err := result.Scan(&name, &password)
		if err != nil {
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

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: multiserver_converter <in> <out>")
		os.Exit(1)
	}

	if err := convertDB(os.Args[1], os.Args[2]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Database converted successfully")
}
