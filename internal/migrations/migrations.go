package migrations

import (
    "database/sql"
)

func Up(tx *sql.Tx) error {
    _, err := tx.Exec(`INSERT INTO users (name, email) VALUES ('admin', 'admin@example.com')`)
    return err
}

func Down(tx *sql.Tx) error {
    _, err := tx.Exec(`DELETE FROM users WHERE email = 'admin@example.com'`)
    return err
}