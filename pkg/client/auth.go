package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"software.sslmate.com/src/go-pkcs12"
)

// loadP12Bundle loads a PKCS#12 bundle and returns TLS certificates
func loadP12Bundle(path string) (tls.Certificate, *x509.CertPool, error) {
	// Read the P12 file
	p12Data, err := os.ReadFile(path)
	if err != nil {
		return tls.Certificate{}, nil, fmt.Errorf("failed to read P12 bundle: %w", err)
	}

	// Get password from environment variable
	password := os.Getenv("VES_P12_PASSWORD")
	if password == "" {
		return tls.Certificate{}, nil, fmt.Errorf("VES_P12_PASSWORD environment variable is not set")
	}

	// Decode the P12 bundle
	privateKey, cert, caCerts, err := pkcs12.DecodeChain(p12Data, password)
	if err != nil {
		return tls.Certificate{}, nil, fmt.Errorf("failed to decode P12 bundle: %w", err)
	}

	// Create the TLS certificate
	tlsCert := tls.Certificate{
		Certificate: [][]byte{cert.Raw},
		PrivateKey:  privateKey,
		Leaf:        cert,
	}

	// Add any intermediate certificates to the chain
	for _, caCert := range caCerts {
		tlsCert.Certificate = append(tlsCert.Certificate, caCert.Raw)
	}

	// Create CA pool from the chain certificates
	caPool := x509.NewCertPool()
	for _, caCert := range caCerts {
		caPool.AddCert(caCert)
	}

	return tlsCert, caPool, nil
}

// loadCertKeyPair loads a certificate and key from separate files
func loadCertKeyPair(certPath, keyPath string) (tls.Certificate, error) {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to load cert/key pair: %w", err)
	}
	return cert, nil
}

// loadCACert loads a CA certificate file
func loadCACert(path string) (*x509.CertPool, error) {
	caCert, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA cert: %w", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to parse CA cert")
	}

	return caCertPool, nil
}

// createTLSConfig creates a TLS configuration from the provided credentials
func createTLSConfig(cfg *Config) (*tls.Config, error) {
	tlsConfig := &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: cfg.InsecureSkipVerify,
	}

	// Load client certificate
	if cfg.P12Bundle != "" {
		// Use P12 bundle
		cert, caPool, err := loadP12Bundle(cfg.P12Bundle)
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
		// Only use the CA pool from P12 if we're not skipping verification
		if caPool != nil && !cfg.InsecureSkipVerify {
			tlsConfig.RootCAs = caPool
		}
	} else if cfg.Cert != "" && cfg.Key != "" {
		// Use cert/key pair
		cert, err := loadCertKeyPair(cfg.Cert, cfg.Key)
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	// Load CA cert if provided separately
	if cfg.CACert != "" && !cfg.InsecureSkipVerify {
		caPool, err := loadCACert(cfg.CACert)
		if err != nil {
			return nil, err
		}
		tlsConfig.RootCAs = caPool
	}

	return tlsConfig, nil
}
