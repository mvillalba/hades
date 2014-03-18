package main

import "github.com/mvillalba/hades"

type actionFunc func(hades.LicenseDatastore, []string) error

// [CLASS][ACTION] -> fn()
var actionMap = map[string]map[string]actionFunc {
    "producer": {
        "create": producerCreateAction,
        "list": producerListAction,
        "dump": producerDumpAction,
        "export": producerExportAction,
    },
}
