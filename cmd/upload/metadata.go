package main

import (
	"database/sql"
	"log"
  "fmt"

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

func update_meta_status(sqlDB *sql.DB, column string, tag string, date string) {
  query := fmt.Sprintf(`UPDATE metadata SET %s = TRUE WHERE tag = ? AND date = ?`, column)
  _, err := sqlDB.Exec(query, tag, date)
	if err != nil {
		log.Fatal(err)
	}
}

func get_meta_status(sqlDB *sql.DB, column string, tag string, date string) (bool) {
  var status bool
  query := fmt.Sprintf(`SELECT %s FROM metadata WHERE tag = ? AND date = ?`, column)
  row := sqlDB.QueryRow(query, tag, date)
  err:= row.Scan(&status)
	if err != nil {
		log.Fatal(err)
	}
  if !status {
    return true
  }
  return false
}

func ready_for_upload(sqlDB *sql.DB, tag string, date string) (bool) {
  var total int
	query := `
		SELECT (SUM(mixcloud) + SUM(radiocult) + SUM(drive)) AS total_sum
		FROM metadata
		WHERE tag = ? AND date = ?`
  row := sqlDB.QueryRow(query, tag, date)
  err:= row.Scan(&total)
	if err != nil {
		log.Fatal(err)
    return false 
	}
  if total >= 3 {
    return true
  } else {
    return false
  }
}

func check_row(sqlDB *sql.DB, tag string, date string) (bool, error) {
	query := "SELECT 1 FROM metadata WHERE tag = ? AND date = ? LIMIT 1"
	var exists int

	err := sqlDB.QueryRow(query, tag, date).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
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
  
  exists, err := check_row(sqlDB, metadata.tag, date)
  if !exists {
    stmt, err := sqlDB.Prepare(`
      INSERT INTO metadata
      VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
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
      0,
      1, // because we havent implemented radiocult yet
      0,
    )
    if err != nil {
        log.Fatal(err)
    }
  } else {
    log.Println("Row already exists")
  }
}
