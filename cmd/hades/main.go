package main

import (
    "os"
    "fmt"
    "strings"
    "github.com/mvillalba/hades"
)

type Args struct {
    class       string
    action      string
    args        []string
}

func usage() {
        fmt.Println("Hades version", hades.VersionNumber())
        fmt.Println("Hades is a generic software licensing platform written in Go (http://github.com/mvillalba/hades).")
        fmt.Println("Copyright © 2014 Martín Raúl Villalba <martin@martinvillalba.com>")
        fmt.Println("")
        fmt.Printf("Usage: %s ACTION [ACTION-ARGS]\n", os.Args[0])
        fmt.Println(`
Valid actions and their arguments are:
    create-producer NAME CITY REGION COUNTRY [COMMENT]: create a new software
        producer.
    create-product PRODUCER-KEY NAME [COMMENT]: create a new product with the
        name NAME under the producer identified by key PRODUCER-KEY.
    create-version PRODUCT-KEY NAME [COMMENT]: create a new version with the
        name NAME under the product identified with key PRODUCT-KEY.
    create-enduser VERSION-KEY NAME CITY REGION COUNTRY [COMMENT]: create an
        end-user license for client NAME living in CITY, REGION, COUNTRY.
    create-alias KEY-ID ALIAS: create an alias ALIAS for KEY-ID. Whereever
        KEY-ID could be used, ALIAS will also be accepted. Note that aliases
        will not be exported with the license.
    list-producers: producer a list of known producers by key (with aliases).
    list-products PRODUCER-KEY: producer a list of known producers by key (with
        aliases).
    list-versions PRODUCT-KEY: producer a list of known producers by key (with
        aliases).
    list-enduser VERSION-KEY: producer a list of known producers by key (with
        aliases).
    list-all: producer a list of all known licenses by key (with aliases).
    export FILENAME KEY-ID-1[...KEY-ID-N]: dump one or more licenses to a file
        (no signing key).
    export-private FILENAME KEY-ID-1[...KEY-ID-N]: dump one or more licenses
        to a file (WITH signing key).
    export-alias FILENAME ALIAS-1[...ALIAS-N]: export one or more aliases to
        FILENAME.
    export-all FILENAME: dump all licenses to a file (no signing key).
    export-all-private FILENAME: dump all licenses to a file (WITH signing key).
    export-all-alias FILENAME: export all aliases to FILENAME.
    import FILENAME: import and verify all licenses contained in FILENAME.
    import-alias FILENAME: import all aliases contained in FILENAME.
    info KEY-ID: pretty-print all information for a given license identified by
        KEY-ID.
    verify KEY-ID: verify license identified by KEY-ID.
    verify-all: verify all known licenses.
    revoke KEY-ID: mark a given license as revoked (making it unavailable for
        future use as a parent license but available to verify pre-existing
        ones).
    destroy KEY-ID: permanently delete a license and all its descendants. Don't
        do this; it's a bad idea. Revoke the license instead.
`)
        os.Exit(1)
}

func processArgs() (a Args) {
    if len(os.Args) < 3 { usage() }

    a.class = strings.ToLower(os.Args[1])
    a.action = strings.ToLower(os.Args[2])

    if len(os.Args) > 3 {
        a.args = os.Args[3:]
    }

    return a
}

// TODO: make this cross-platform
func getDatabasePath() (string, error) {
    cu, err := user.Current()
    if err != nil { return "", err }
    hadespath := cu.HomeDir + "/.hades"
    os.Mkdir(hadespath, 0775)
    return hadespath + "/license.db", nil
}

func main() {
    args := processArgs()
    if fn, ok := actionMap[args.class][args.action]; ok {
        var db hades.SQLite3Datastore
        dbpath, err := getDatabasePath()
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
        err = db.SetupDatabase(dbpath)
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
        defer db.Close()

        err = fn(db, args.args)
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
    } else {
        fmt.Println("Action not found.")
        os.Exit(1)
    }
}
