package main

import (
	"crypto/x509"
	"testing"
)

func TestGenerateSelfSignedCert(t *testing.T) {
	cert, err := GenerateSelfSignedCert()
	if err != nil {
		t.Fatalf("Failed to generate cert: %v", err)
	}

	if len(cert.Certificate) == 0 {
		t.Fatal("Generated certificate is empty")
	}

	leaf, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		t.Fatalf("Failed to parse generated certificate: %v", err)
	}

	expectedCN := "proxy.whatsapp.net"
	if leaf.Subject.CommonName != expectedCN {
		t.Errorf("Expected CN %s, got %s", expectedCN, leaf.Subject.CommonName)
	}

	foundWhatsApp := false
	for _, name := range leaf.DNSNames {
		if name == "whatsapp.net" {
			foundWhatsApp = true
			break
		}
	}
	if !foundWhatsApp {
		t.Error("whatsapp.net missing from DNS SANs")
	}
}
