/*
 * Generates a csr and private key with OpenSSL.
 * OpenSSL config format taken from http://apfelboymchen.net/gnu/notes/openssl%20multidomain%20with%20config%20files.html
 */
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	bits := flag.String("bits", "2048", "The private key size")
	md := flag.String("md", "sha256", "The hash algorithm to use")
	keyName := flag.String("key-name", "", "The filename of the private key")
	csrName := flag.String("csr-name", "", "The filename for the CSR")
	countryName := flag.String("country", "", "The country name")
	stateOrProvinceName := flag.String("state-or-province", "", "The state or province name")
	localityName := flag.String("locality", "", "The locality name")
	postalCode := flag.String("postal-code", "", "The postal code")
	streetAddress := flag.String("street-address", "", "The street address")
	organizationName := flag.String("org-name", "", "The organization name")
	organizationalUnitName := flag.String("org-unit-name", "", "The organizational unit name")
	commonName := flag.String("common-name", "", "A single common name")
	emailAddress := flag.String("email-address", "", "The contact email address")
	flag.Parse()
	if flag.NFlag() < 11 {
		fmt.Println("Invalid arguments")
		os.Exit(1)
	}

	// Write a temporary OpenSSL config to disk.
	f, err := os.Create("csr.cnf")
	if err != nil {
		fmt.Println("Failed to open file: " + err.Error())
		os.Exit(1)
	}
	defer f.Close()
	defer os.Remove("csr.cnf")
	f.WriteString("[ req ]\n")
	f.WriteString("default_bits = " + *bits + "\n")
	f.WriteString("default_md = " + *md + "\n")
	f.WriteString("prompt = no\n")
	f.WriteString("encrypt_key = no\n")
	f.WriteString("distinguished_name = req_distinguished_name\n")
	f.WriteString("[ req_distinguished_name ]\n")
	f.WriteString("countryName = " + *countryName + "\n")
	f.WriteString("stateOrProvinceName  = " + *stateOrProvinceName + "\n")
	f.WriteString("localityName  = " + *localityName + "\n")
	f.WriteString("postalCode  = " + *postalCode + "\n")
	f.WriteString("streetAddress  = " + *streetAddress + "\n")
	f.WriteString("organizationName  = " + *organizationName + "\n")
	f.WriteString("organizationalUnitName  = " + *organizationalUnitName + "\n")
	f.WriteString("commonName  = " + *commonName + "\n")
	f.WriteString("emailAddress  = " + *emailAddress + "\n")

	// Generate the CSR and private key with openssl.
	cmd := "openssl"
	args := []string{"req", "-config", "csr.cnf", "-new", "-out", *csrName, "-keyout", *keyName}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		fmt.Println("Failed to generate CSR and private key: " + err.Error())
		os.Exit(1)
	}
}
