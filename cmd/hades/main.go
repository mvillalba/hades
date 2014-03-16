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
        fmt.Printf("Usage: %s ACTION-CLASS ACTION-NAME [ACTION-ARGS]\n", os.Args[0])
        fmt.Println(`
Valid actions and their args are:
    producer:
        create SLUG-NAME LONG-NAME [COMMENT]: create a new software producer.
        list: produce a list of all known producers as identified by their slug names.
        dump SLUG-1[...SLUG-N]: dump one or more producers to stdout.
        export FILENAME SLUG-1[...SLUG-N]: dump one or more producers to a file.
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
        db, err := setupDatabase()
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
