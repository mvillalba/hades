package hades

import (
    "encoding/hex"
    "crypto/sha1"
    "crypto/rand"
    "errors"
    "strings"
    "fmt"
)

// Key format is as follows, in NAME(BYTES) format:
// CHECKSUM(4): truncated SHA-1 hash of parent keys (producer, product,
//              version) and the current key with the CHECKSUM field zeroed
//              out. Inclusion of the parent keys is mandatory if they exist.
// EXTRA-1(2): available for custom use.
// EXTRA-2(2): available for custom use.
// EXTRA-3(2): available for custom use.
// RANDOM(6): this is the key proper, so to speak. It makes a given key unique.
// Please note that unused EXTRA fields can be used to increase the key space.
type LicenseKey []byte

func (k LicenseKey) IsValid() bool {
    if len([]byte(k)) != 16 { return false }
    return true
}

func (k LicenseKey) String() string {
    b := k.Bytes()
    return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

func (k LicenseKey) Bytes() []byte {
    return []byte(k)
}

type LicenseKeyExtra []byte

func (k LicenseKeyExtra) IsValid() bool {
    if len([]byte(k)) != 2 { return false }
    return true
}

func genLicenseKeyExtra() (LicenseKeyExtra, error) {
    extra := make([]byte, 2)
    _, err := rand.Read(extra)
    if err != nil { return nil, err }
    return extra, nil
}

func NewLicenseKey(extra1, extra2, extra3 LicenseKeyExtra, parents ...LicenseKey) (LicenseKey, error) {
    key := LicenseKey(make([]byte, 16))

    // Generate key proper
    random := make([]byte, 6)
    _, err := rand.Read(random)
    if err != nil { return nil, err }

    // Generate missing key extras
    if extra1 == nil {
        extra1, err = genLicenseKeyExtra()
        if err != nil { return nil, err }
    }
    if extra2 == nil {
        extra2, err = genLicenseKeyExtra()
        if err != nil { return nil, err }
    }
    if extra3 == nil {
        extra3, err = genLicenseKeyExtra()
        if err != nil { return nil, err }
    }

    // Validate key extras
    if !extra1.IsValid() || !extra2.IsValid() || !extra3.IsValid() {
        return nil, errors.New("Received invalid key extras!")
    }

    // Set EXTRA-{1,2,3}
    key[4] = extra1[0]
    key[5] = extra1[1]
    key[6] = extra2[0]
    key[7] = extra2[1]
    key[8] = extra3[0]
    key[9] = extra3[1]

    // Set RANDOM
    key[10] = random[0]
    key[11] = random[1]
    key[12] = random[2]
    key[13] = random[3]
    key[14] = random[4]
    key[15] = random[5]

    // Calculate CHECKSUM
    h := sha1.New()
    for _, k := range parents {
        if !k.IsValid() {
            return nil, errors.New("Received invalid key parents!")
        }
        h.Write(k)
    }
    h.Write(key)
    checksum := h.Sum(nil)

    // Set CHECKSUM
    key[0] = checksum[0]
    key[1] = checksum[1]
    key[2] = checksum[2]
    key[3] = checksum[3]

    return key, nil
}

func LicenseKeyFromString(k string) (key LicenseKey, err error) {
    // Validate
    err = errors.New("Not a key.")
    parts := strings.Split(k, "-")
    if len(parts[0]) != 8  { return key, err }
    if len(parts[1]) != 4  { return key, err }
    if len(parts[2]) != 4  { return key, err }
    if len(parts[3]) != 4  { return key, err }
    if len(parts[4]) != 12 { return key, err }

    // Build
    k = strings.Replace(k, "-", "", -1)
    bytes, err := hex.DecodeString(k)
    if err != nil { return key, err }
    key = LicenseKey(bytes)

    return key, nil
}
