package db

import "database/sql"

func ScanOne(row *sql.Row, dest ...interface{}) error {
	err := row.Scan(dest...)
	if err == sql.ErrNoRows {
		return nil // Not found, but not a hard error
	}
	return err
}
