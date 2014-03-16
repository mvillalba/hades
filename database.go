package hades

// LicenseDatastore is an interface designed to make the license database
// easily replaceable so the license server and the CLI application can use
// different database engines based on their different use cases. Plus, this
// way I get to change my mind about the datastore later on without having to
// make a lot of changes.
type LicenseDatastore interface {
    SetupDatabase(url string) error
    GetLicense(LicenseKey) (*License, error)
    GetLicenseList() (*[]LicenseKey, error)
    GetLicenseListByParent(LicenseKey) (*License, error)
    StoreLicense(*License) error
    DeleteLicense(LicenseKey) error
}

