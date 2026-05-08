package config

import _ "embed"

//go:embed certs/ca.pem
var embeddedCA []byte

//go:embed certs/client.pem
var embeddedClientCert []byte

//go:embed certs/client-key.pem
var embeddedClientKey []byte

func EmbeddedCACert() []byte  { return embeddedCA }
func EmbeddedClientCert() []byte  { return embeddedClientCert }
func EmbeddedClientKey() []byte  { return embeddedClientKey }
