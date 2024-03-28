package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/antonio-petrillo/dixieflatline/message"
)

func SetupDB(dbfile string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		return nil, err
	}

	setupHistoryTableStmt := `
CREATE TABLE IF NOT EXISTS
history(channel TEXT, user TEXT, message TEXT, timestamp DATETIME DEFAULT CURRENT_TIMESTAMP);`

	if _, err = db.Exec(setupHistoryTableStmt); err != nil {
		return nil, err
	}

	return db, nil
}

func StoreHistoryEntry(db *sql.DB, channel string, entry message.HistoryEntry) error {
	transaction, err := db.Begin()

	insertStmt, err := transaction.Prepare(`
INSERT INTO history(channel, user, message)
VALUES(?, ?, ?)`)

	if err != nil {
		return err
	}
	defer insertStmt.Close()

	_, err = insertStmt.Exec(channel, entry.From, entry.Msg)
	if err != nil {
		if errRoll := transaction.Rollback(); errRoll != nil {
			return errRoll
		}
		return err
	}

	err = transaction.Commit()
	return err
}

func RetrieveHistoryEntries(db *sql.DB, channel string, limit int) ([]message.HistoryEntry, error) {
	// I don't care if the reading is not perfectly in sync
	// In other words no need for transaction
	// transaction, err := db.Begin()

	rows, err := db.Query(`
SELECT * FROM history
WHERE channel = ?
ORDER BY timestamp DESC LIMIT ?`, channel, limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []message.HistoryEntry
	for rows.Next() {
		var ch, user, msg, timestamp string
		if err := rows.Scan(&ch, &user, &msg, &timestamp); err != nil {
			// return partial result
			return entries, err
		}
		entries = append(entries, message.HistoryEntry{From: user, Msg: msg})
	}
	// if err = rows.Err(); err != nil {
	// 	return entries, err
	// }
	// // err = transaction.Commit()

	err = rows.Err()
	// return partial result
	return entries, err
}
