package hades

import (
    "os"
    "encoding/json"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

// SQLite3Datastore implements LicenseDatastore for SQLite3-based storage.
type SQLite3Datastore struct {
    db      *sql.DB
}

func (s *SQLite3Datastore) SetupDatabase(dbpath string) error {
    // Does the DB exist or do we need to create it?
    init := false
    if _, err := os.Stat(dbpath); os.IsNotExist(err) {
        init = true
    }

    // Open DB
    db, err := sql.Open("sqlite3", dbpath)
    if err != nil { return err }

    // Let's make sure everything is ok…
    err = db.Ping()
    if err != nil { return err }

    // Initialize database
    if init {
        err := s.initDatabase()
        if err != nil { return err }
    }

    return nil
}

func (s *SQLite3Datastore) initDatabase() error {
    _, err := s.db.Exec(`
    CREATE TABLE licenses (
        id              CHAR(36)    NOT NULL    PRIMARY KEY,
        parent_id       CHAR(36)    NOT NULL,
        license         BLOB        NOT NULL,
    );
    CREATE INDEX parent_index ON licenses (parent_id);
    `)
    if err != nil { return err }

    return nil
}

func (s *SQLite3Datastore) StoreLicense(license *License) error {
    // Encode
    jsonLicense, err := json.Marshal(license)
    if err != nil { return err }

    // Store
    _, err = s.db.Exec("INSERT INTO licenses(id, parent_id, license) VALUES (?, ?, ?)",
                       license.ID, license.ParentID, jsonLicense)
    if err != nil { return err }

    return nil
}

func (s *SQLite3Datastore) DeleteLicense(id LicenseKey) error {
    // Delete
    _, err := s.db.Exec("DELETE FROM licenses WHERE ID=?", id)
    if err != nil { return err }

    return nil
}

func (s *SQLite3Datastore) GetLicense(id LicenseKey) (*License, error) {
    // Query DB
    var jsonLicense []byte
    err := s.db.QueryRow("SELECT license FROM licenses WHERE id=?",
                         id).Scan(&jsonLicense)
    if err != nil { return nil, err }

    // Decode
    var license License
    err = json.Unmarshal(jsonLicense, &license)
    if err != nil { return nil, err }

    return &license, nil
}

func (s *SQLite3Datastore) GetLicenseList() (keys []LicenseKey, err error) {
    // Query DB
    rows, err := s.db.Query("SELECT id FROM licenses")
    if err != nil { return keys, err }
    defer rows.Close()

    // Build key list
    var key string
    for rows.Next() {
        err = rows.Scan(&key)
        if err != nil { return keys, err }
        lk, err := LicenseKeyFromString(key)
        if err != nil { return keys, err }
        keys = append(keys, lk)
    }
    err = rows.Err()
    if err != nil { return keys, err }

    return keys, nil
}

func (s *SQLite3Datastore) GetLicenseListByParent(parent LicenseKey) (keys []LicenseKey, err error) {
    // Query DB
    rows, err := s.db.Query("SELECT id FROM licenses WHERE parent_id=?",
                            parent)
    if err != nil { return keys, err }
    defer rows.Close()

    // Build key list
    var key string
    for rows.Next() {
        err = rows.Scan(&key)
        if err != nil { return keys, err }
        lk, err := LicenseKeyFromString(key)
        if err != nil { return keys, err }
        keys = append(keys, lk)
    }
    err = rows.Err()
    if err != nil { return keys, err }

    return keys, nil
}
