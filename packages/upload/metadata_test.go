package main

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}

	schema := `
	CREATE TABLE metadata (
		date TEXT,
		tag TEXT,
		show_name TEXT,
		show_nr INTEGER,
		dj_name TEXT,
		picture TEXT,
		description TEXT,
		tags0 TEXT,
		tags1 TEXT,
		tags2 TEXT,
		tags3 TEXT,
		tags4 TEXT,
		live BOOLEAN,
		mixcloud BOOLEAN,
		radiocult BOOLEAN,
		drive BOOLEAN,
		PRIMARY KEY (date, tag)
	);

	CREATE TABLE base_data (
		tag TEXT PRIMARY KEY,
		show_name TEXT,
		show_nr INTEGER,
		dj_name TEXT,
		picture TEXT,
		description TEXT,
		tags0 TEXT,
		tags1 TEXT,
		tags2 TEXT,
		tags3 TEXT,
		tags4 TEXT,
		live BOOLEAN
	);`

	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("create schema: %v", err)
	}

	t.Cleanup(func() {
		_ = db.Close()
	})

	return db
}

func testMetadata() Metadata {
	return Metadata{
		tag:         "test-show",
		show_name:   "Test Show",
		dj_name:     "DJ Test",
		picture:     "cover.jpg",
		description: "Description",
		tags0:       "ambient",
		tags1:       "electronic",
		tags2:       "",
		tags3:       "",
		tags4:       "",
		live:        true,
	}
}

func TestMigrateMetadataSchemaAddsSoundcloudColumns(t *testing.T) {
	db := openTestDB(t)

	migrateMetadataSchema(db)
	migrateMetadataSchema(db)

	rows, err := db.Query(`PRAGMA table_info(metadata)`)
	if err != nil {
		t.Fatalf("query table info: %v", err)
	}
	defer rows.Close()

	columns := map[string]bool{}
	for rows.Next() {
		var (
			cid       int
			name      string
			columnType string
			notNull   int
			defaultVal sql.NullString
			pk        int
		)
		if err := rows.Scan(&cid, &name, &columnType, &notNull, &defaultVal, &pk); err != nil {
			t.Fatalf("scan table info: %v", err)
		}
		columns[name] = true
	}

	if !columns["soundcloud"] {
		t.Fatalf("soundcloud column not added")
	}
	if !columns["soundcloud_urn"] {
		t.Fatalf("soundcloud_urn column not added")
	}
}

func TestNewMetaRowInitializesSoundcloudFields(t *testing.T) {
	db := openTestDB(t)
	migrateMetadataSchema(db)

	new_meta_row(db, "20260318", testMetadata())

	var soundcloud bool
	var soundcloudURN sql.NullString
	var showNr int
	err := db.QueryRow(`
		SELECT soundcloud, soundcloud_urn, show_nr
		FROM metadata
		WHERE date = ? AND tag = ?
	`, "20260318", "test-show").Scan(&soundcloud, &soundcloudURN, &showNr)
	if err != nil {
		t.Fatalf("load inserted row: %v", err)
	}

	if soundcloud {
		t.Fatalf("soundcloud should default to false")
	}
	if soundcloudURN.Valid && soundcloudURN.String != "" {
		t.Fatalf("soundcloud_urn should default to empty, got %q", soundcloudURN.String)
	}
	if showNr != 1 {
		t.Fatalf("expected initial show_nr 1, got %d", showNr)
	}
}

func TestNewMetaRowMarksRadiocultCompleteForPrerecord(t *testing.T) {
	db := openTestDB(t)
	migrateMetadataSchema(db)

	metadata := testMetadata()
	metadata.live = false
	new_meta_row(db, "20260318", metadata)

	var radiocult bool
	err := db.QueryRow(`
		SELECT radiocult
		FROM metadata
		WHERE date = ? AND tag = ?
	`, "20260318", "test-show").Scan(&radiocult)
	if err != nil {
		t.Fatalf("load radiocult value: %v", err)
	}

	if !radiocult {
		t.Fatalf("prerecord should mark radiocult as already complete")
	}
}

func TestReadyForUploadRequiresAllFourDestinations(t *testing.T) {
	db := openTestDB(t)
	migrateMetadataSchema(db)
	new_meta_row(db, "20260318", testMetadata())

	tag := "test-show"
	date := "20260318"

	update_meta_status(db, "mixcloud", tag, date)
	update_meta_status(db, "soundcloud", tag, date)
	update_meta_status(db, "radiocult", tag, date)

	if ready_for_upload(db, tag, date) {
		t.Fatalf("row should not be ready before drive upload completes")
	}

	update_meta_status(db, "drive", tag, date)

	if !ready_for_upload(db, tag, date) {
		t.Fatalf("row should be ready after all four destinations complete")
	}
}

func TestUpdateMetaValueStoresSoundcloudURN(t *testing.T) {
	db := openTestDB(t)
	migrateMetadataSchema(db)
	new_meta_row(db, "20260318", testMetadata())

	update_meta_value(db, "soundcloud_urn", "soundcloud:tracks:123", "test-show", "20260318")

	var got string
	err := db.QueryRow(`
		SELECT soundcloud_urn
		FROM metadata
		WHERE date = ? AND tag = ?
	`, "20260318", "test-show").Scan(&got)
	if err != nil {
		t.Fatalf("load soundcloud urn: %v", err)
	}

	if got != "soundcloud:tracks:123" {
		t.Fatalf("expected soundcloud urn to persist, got %q", got)
	}
}
