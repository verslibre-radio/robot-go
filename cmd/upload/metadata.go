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

func get_metadata(sqlDB *sql.DB, tag string) Metadata {
	var metadata Metadata
	row := sqlDB.QueryRow(`SELECT * FROM base_data WHERE tag =  ?`, tag)

	err := row.Scan(&metadata.tag, &metadata.show_name, &metadata.show_nr, &metadata.dj_name, &metadata.picture, &metadata.description, &metadata.tags0, &metadata.tags1, &metadata.tags2, &metadata.tags3, &metadata.tags4)
	if err != nil {
	}

	return metadata
}

func update_show_nr(sqlDB *sql.DB, tag string) {
	result, err := sqlDB.Exec("UPDATE base_data SET show_nr = show_nr + ? WHERE tag = ?", 1, tag)
	if err != nil {
		log.Fatal(err)
	}

	_, err = result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
}

func new_meta_row(sqlDB *sql.DB, date string, metadata Metadata) {
  var max_nr sql.NullInt64
  var show_nr int

	row := sqlDB.QueryRow(`SELECT MAX(show_nr) FROM metadata WHERE tag =  ?`, metadata.tag)
  err := row.Scan(&max_nr)
	if err != nil {
		log.Fatal(err)
	}
  if max_nr.Valid {
    show_nr = int(max_nr.Int64 + 1)
  } else {
    show_nr = 1
  }

  stmt, err := sqlDB.Prepare(`
    INSERT INTO metadata
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
  `)
	if err != nil {
		log.Fatal(err)
	}

  _, err = stmt.Exec(
    date,
    metadata.tag,
    metadata.show_name,
    show_nr,
    metadata.dj_name,
    metadata.picture,
    metadata.description,
    metadata.tags0,
    metadata.tags1,
    metadata.tags2,
    metadata.tags3,
    metadata.tags4,
    "FALSE",
    "FALSE",
    "FALSE",
    "FALSE",
  )
  if err != nil {
      log.Fatal(err)
  }
}
