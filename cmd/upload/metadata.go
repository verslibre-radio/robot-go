package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/api/sheets/v4"
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
	live        bool
}

func check_custom_metadata(sheet *sheets.ValueRange, metadata Metadata, date string, tag string) (bool, Metadata) {
	for _, row := range sheet.Values {
		if row[0] == date && row[1] == tag {
			for len(row) < 12 {
				row = append(row, "")
			}
			if row[2] != "" {
				metadata.show_name = row[2].(string)
			}
			if row[3] != "" {
				metadata.dj_name = row[3].(string)
			}
			if row[4] != "" {
				metadata.picture = row[4].(string)
			}
			if row[5] != "" {
				metadata.description = row[5].(string)
			}
			if row[6] != "" {
				metadata.tags0 = row[6].(string)
			}
			if row[7] != "" {
				metadata.tags1 = row[7].(string)
			}
			if row[8] != "" {
				metadata.tags2 = row[8].(string)
			}
			if row[9] != "" {
				metadata.tags3 = row[9].(string)
			}
			if row[10] != "" {
				metadata.tags4 = row[10].(string)
			}
			if row[11] != "" {
				metadata.live = row[11].(bool)
			}
			return true, metadata
		}
	}
	return false, metadata
}

func get_metadata(sheet_meta *sheets.ValueRange, sqlDB *sql.DB, tag string, date string) Metadata {
	var metadata Metadata
	row := sqlDB.QueryRow(`SELECT * FROM base_data WHERE tag =  ?`, tag)

	err := row.Scan(&metadata.tag, &metadata.show_name, &metadata.show_nr, &metadata.dj_name, &metadata.picture, &metadata.description, &metadata.tags0, &metadata.tags1, &metadata.tags2, &metadata.tags3, &metadata.tags4, &metadata.live)
	if err != nil {
		log.Fatal(err)
	}

	add, custom := check_custom_metadata(sheet_meta, metadata, date, tag)
	if add {
		return custom
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

func get_meta_status(sqlDB *sql.DB, column string, tag string, date string) bool {
	var status bool
	query := fmt.Sprintf(`SELECT %s FROM metadata WHERE tag = ? AND date = ?`, column)
	row := sqlDB.QueryRow(query, tag, date)
	err := row.Scan(&status)
	if err != nil {
		log.Fatal(err)
	}
	if !status {
		return true
	}
	return false
}

func ready_for_upload(sqlDB *sql.DB, tag string, date string) bool {
	var total int
	query := `
		SELECT (SUM(mixcloud) + SUM(radiocult) + SUM(drive)) AS total_sum
		FROM metadata
		WHERE tag = ? AND date = ?`
	row := sqlDB.QueryRow(query, tag, date)
	err := row.Scan(&total)
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
			metadata.live,
			0,
			0,
      0,
		)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("Row already exists")
	}
}
