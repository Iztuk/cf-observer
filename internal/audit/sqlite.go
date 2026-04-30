package audit

import (
	"context"
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func InitDatabase(dbPath string, overwrite bool) error {
	if overwrite {
		if err := os.Remove(dbPath); err != nil && !os.IsExist(err) {
			return err
		}
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	ctx := context.Background()

	_, err = db.ExecContext(ctx, `PRAGMA foreign_keys = ON;`)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS audit_jobs (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL,

    request_id TEXT NOT NULL,
    host TEXT NOT NULL,
    method TEXT NOT NULL,
    path TEXT NOT NULL,
    query TEXT,
    upstream TEXT NOT NULL,
    status INTEGER,
    timestamp TEXT NOT NULL,
    duration_ms INTEGER NOT NULL,

    headers TEXT,
    body TEXT,
    error TEXT
);
`)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS findings (
    id TEXT PRIMARY KEY,
    job_id TEXT NOT NULL,

    rule_id TEXT NOT NULL,
    title TEXT NOT NULL,
    message TEXT NOT NULL,

    request_id TEXT NOT NULL,
    host TEXT NOT NULL,
    method TEXT NOT NULL,
    path TEXT NOT NULL,
    query TEXT,
    status_code INTEGER,

    created_at TEXT NOT NULL,

    FOREIGN KEY (job_id) REFERENCES audit_jobs(id)
);
`)
	if err != nil {
		return err
	}

	return nil
}
