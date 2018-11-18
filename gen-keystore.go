/*
 * Creates a Java keystore from a public and private key.
 */
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	privateKey := flag.String("private-key", "", "The file name of the private key")
	publicKey := flag.String("public-key", "", "The file name of the public key")
	alias := flag.String("alias", "", "The alias to store the key under")
	caCert := flag.String("ca-cert", "", "The CA root cert file name")
	intermediateCerts := flag.String("intermediate-certs", "", "A comma-separated list of intermediate cert file names")
	flag.Parse()
	if flag.NFlag() < 5 {
		fmt.Println("Invalid arguments")
		os.Exit(1)
	}

	// Store the chain in a pkcs12 archive so that it can be imported into a keystore.
	cmd := "openssl"
	args := []string{"pkcs12", "-export", "-in", *publicKey, "-inkey", *privateKey, "-out", "temp.p12", "-name", *alias, "-CAfile", *caCert, "-password", "pass:changeit"}
	intermediateCertList := strings.Split(*intermediateCerts, ",")
	for i, _ := range intermediateCertList {
		args = append(args, "-certfile")
		args = append(args, intermediateCertList[i])
	}
	// Generate pkcs12.
	if err := exec.Command(cmd, args...).Run(); err != nil {
		fmt.Println("Failed to generate pkcs12 archive: " + err.Error())
		os.Exit(1)
	}
	defer os.Remove("temp.p12")

	// Generate keystore and import pkcs12.
	cmd = "keytool"
	args = []string{"-importkeystore", "-deststorepass", "changeit", "-destkeypass", "changeit", "-destkeystore", *alias + ".jks", "-srckeystore", "temp.p12", "-srcstoretype", "PKCS12", "-srcstorepass", "changeit", "-alias", *alias}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		fmt.Println("Failed to generate keystore archive: " + err.Error())
		os.Exit(1)
	}
}
