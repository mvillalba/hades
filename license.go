package hades

import (
    "bytes"
    "strings"
    "encoding/json"
    "encoding/base64"
    "time"
)

var ArmorHeader = strings.Repeat("-", 25) + " BEGIN LICENSE " + strings.Repeat("-", 24)
var ArmorFooter = strings.Repeat("-", 25) + " END LICENSE " + strings.Repeat("-", 24)

const ARMOR_MAGIC = "HADES-ARMOR"
const ARMOR_VERSION = 1

type ArmorEnvelope struct {
    Magic       string
    Version     int
    Payload     []byte
    Signature   []byte
}

// A license's type specifies what the licence is intended to represent (a
// producer, a product, etc.) All but LICENSE_CLIENT are master licenses.
type LicenseType string
var LICENSE_PRODUCER = LicenseType("PRODUCER")
var LICENSE_PRODUCT = LicenseType("PRODUCT")
var LICENSE_VERSION = LicenseType("VERSION")
var LICENSE_ENDUSER = LicenseType("ENDUSER")

// A license's class (master, client) essentially specifies what the license is
// to be used for, in broad terms. A master license is that of a software
// producer or program, and is used to sign other master licenses and client
// licenses and there can only be one master license for a given thing
// (software producer, program, program version, etc.). A client license, in
// turn has no key pair, can't be used for signing other licenses, and is the
// only license class a Hades-protected program can load.
type LicenseClass string
var LICENSE_MASTER = LicenseClass("MASTER")
var LICENSE_CLIENT = LicenseClass("CLIENT")

// Signing key pair type and bit-length.
type LicenseKeyPairType string
var KEYPAIR_ECDSA = LicenseKeyPairType("ECDSA-2048")

type License struct {
    ID          LicenseKey
    ParentID    LicenseKey
    Type        LicenseType
    Class       LicenseClass
    Created     time.Time
    Metadata    map[string]string
    KeyPairType LicenseKeyPairType
    PrivateKey  []byte
    PublicKey   []byte
}

func (l *License) Armor() ([]byte, error) {
    // Encode payload
    jsonPayload, err := json.Marshal(l)
    if err != nil { return nil, err }

    // Sign payload
    var signature []byte
    signature = []byte("SIGNATURE-GOES-HERE")
    // TODO

    // Create envelope
    var envelope ArmorEnvelope
    envelope.Magic = ARMOR_MAGIC
    envelope.Version = ARMOR_VERSION
    envelope.Payload = jsonPayload
    envelope.Signature = signature

    // Encode envelope
    jsonEnvelope, err := json.Marshal(envelope)
    if err != nil { return nil, err }

    base64Envelope := base64.StdEncoding.EncodeToString(jsonEnvelope)

    // Write envelope to buffer
    buf := new(bytes.Buffer)
    buf.WriteString(ArmorHeader + "\n")

    for i := 0; true; i += 64 {
        if len(base64Envelope) >= i {
            buf.WriteString(base64Envelope[i:i+64])
            buf.WriteString("\n")
            continue
        }
        buf.WriteString(base64Envelope[i:])
        buf.WriteString("\n")
        break
    }

    buf.WriteString(ArmorFooter + "\n")

    return buf.Bytes(), nil
}
/*
// TODO
func NewProducerLicense() (*License, error) {
    license := new(License)
    license.Class = LICENSE_MASTER
    license.Type = LICENSE_PRODUCER
    license.Created = time.Now()
    
}*/
