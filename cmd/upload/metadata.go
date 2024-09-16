package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Metadata struct {
	tag         string
	show_name   string
	show_nr     int
	dj_name     string
	picture     string
	description string
	tags0       string
	tags1       string
	tags2       string
	tags3       string
	tags4       string
}

func get_metadata(db_path string, tag string) Metadata {
	sqlDB, _ := sql.Open("sqlite3", db_path)
	defer sqlDB.Close()
	row := sqlDB.QueryRow(`SELECT * FROM base_data WHERE tag =  ?`, tag)

	var metadata Metadata
	err := row.Scan(&metadata.tag, &metadata.show_name, &metadata.show_nr, &metadata.dj_name, &metadata.picture, &metadata.description, &metadata.tags0, &metadata.tags1, &metadata.tags2, &metadata.tags3, &metadata.tags4)
	if err != nil {
		log.Fatal(err)
	}

	return metadata
}

func update_show_nr(db_path string, tag string) {
	sqlDB, _ := sql.Open("sqlite3", db_path)
	defer sqlDB.Close()
	result, err := sqlDB.Exec("UPDATE base_data SET show_nr = show_nr + ? WHERE tag = ?", 1, tag)
	if err != nil {
		log.Fatal(err)
	}

	_, err = result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
}
