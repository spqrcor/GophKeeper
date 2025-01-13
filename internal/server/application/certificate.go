package application

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"time"
)

// certOptions тип конфига сертификата
type certOptions struct {
	certPath     string
	keyPath      string
	validMonths  int
	organization string
	country      string
}

// certCfg конфиг сертификата
var certCfg = certOptions{
	certPath:     "./server.crt",
	keyPath:      "./server.key",
	validMonths:  12,
	organization: "Рога и копыта",
	country:      "RU",
}

// initCertificate инициализация сертификата
func initCertificate() error {
	fileInfo, err := os.Stat(certCfg.certPath)
	if err == nil && time.Now().Before(fileInfo.ModTime().AddDate(0, certCfg.validMonths, 0)) {
		return nil
	}
	return createCertificate()
}

// createCertificate создание сертификата
func createCertificate() error {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization: []string{certCfg.organization},
			Country:      []string{certCfg.country},
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(0, certCfg.validMonths, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		return err
	}

	var certPEM bytes.Buffer
	if err := pem.Encode(&certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}); err != nil {
		return err
	}
	var privateKeyPEM bytes.Buffer
	if err := pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}); err != nil {
		return err
	}

	if err := os.WriteFile(certCfg.certPath, certPEM.Bytes(), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(certCfg.keyPath, privateKeyPEM.Bytes(), 0600); err != nil {
		return err
	}
	return nil
}
