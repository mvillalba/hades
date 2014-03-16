package main

import "database/sql"

type actionFunc func(*sql.DB, []string) error

// [CLASS][ACTION] -> fn()
var actionMap = map[string]map[string]actionFunc {
    "producer": {
        "create": producerCreateAction,
        "list": producerListAction,
        "dump": producerDumpAction,
        "export": producerExportAction,
    },
}
