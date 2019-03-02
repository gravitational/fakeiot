/*
Copyright 2019 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

package utils

import (
	"crypto/x509"
	"encoding/pem"
	"os"

	"github.com/gravitational/trace"
	log "github.com/sirupsen/logrus"
)

// InitLogger configures the global logger for a given purpose / verbosity level
func InitLogger(level log.Level) {
	log.StandardLogger().SetHooks(make(log.LevelHooks))
	log.SetLevel(level)
	log.SetFormatter(&trace.TextFormatter{
		EnableColors: true,
	})
	log.SetOutput(os.Stderr)
}

// ParseCertificatePEM parses PEM-encoded certificate
func ParseCertificatePEM(bytes []byte) (*x509.Certificate, error) {
	if len(bytes) == 0 {
		return nil, trace.BadParameter("missing PEM encoded block")
	}
	block, _ := pem.Decode(bytes)
	if block == nil {
		return nil, trace.BadParameter("expected PEM-encoded block")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, trace.BadParameter(err.Error())
	}
	return cert, nil
}
