package hades

import "fmt"

const VERSION_MAJOR = 0
const VERSION_MINOR = 1
const VERSION_PATCH = 0

func VersionNumber() string {
    return fmt.Sprintf("%v.%v.%v", VERSION_MAJOR, VERSION_MINOR, VERSION_PATCH)
}
