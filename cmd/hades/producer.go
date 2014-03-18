package main

import (
    "fmt"
    "time"
    "bytes"
    "io/ioutil"
    "database/sql"
    "github.com/mvillalba/hades"
)

func producerCreateAction(db hades.LicenseDatastore, args []string) error {
    // Args
    if len(args) < 2 { usage() }
    slug := args[0]
    name := args[1]
    comment := ""
    if len(args) == 3 { comment = args[2] }

    // Generate unique key and timestamp
    timestamp := time.Now().Unix()
    key, err := NewLicenseKey(nil, nil, nil)
    if err != nil { return err }

    // Store

    return nil
}

func producerListAction(db hades.LicenseDatastore, args []string) error {
    // Args
    if len(args) != 0 { usage() }

    // Query DB
    rows, err := db.Query("SELECT id FROM producers")
    if err != nil { return err }
    defer rows.Close()

    // Print results
    var slug string
    for rows.Next() {
        err = rows.Scan(&slug)
        if err != nil { return err }
        fmt.Println(slug)
    }
    err = rows.Err()
    if err != nil { return err }

    return nil
}

func producerDumpActionImpl(db hades.LicenseDatastore, args[]string) ([]byte, error) {
    // Args
    if len(args) < 1 { usage() }

    // Prepare statement
    stmt, err := db.Prepare("SELECT id, name, key, created, comment FROM producers WHERE ID=?")
    if err != nil { return nil, err }
    defer stmt.Close()

    // Query, encode, dump
    buf := new(bytes.Buffer)
    var id string
    var name string
    var key string
    var created int64
    var comment string
    for _, slug := range args {
        err := stmt.QueryRow(slug).Scan(&id, &name, &key, &created, &comment)
        if err != nil { return nil, err }
        output, err := exportProducer(id, name, key, created, comment)
        if err != nil { return nil, err }
        _, err = buf.Write(output)
        if err != nil { return nil, err }
    }

    return buf.Bytes(), nil
}

func producerDumpAction(db hades.LicenseDatastore, args []string) error {
    output, err := producerDumpActionImpl(db, args)
    if err != nil { return err }
    fmt.Println(string(output))
    return nil
}

func producerExportAction(db hades.LicenseDatastore, args []string) error {
    if len(args) < 2 { usage() }
    output, err := producerDumpActionImpl(db, args[1:])
    if err != nil { return err }
    err = ioutil.WriteFile(args[0], output, 0664)
    if err != nil { return err }
    return nil
}
