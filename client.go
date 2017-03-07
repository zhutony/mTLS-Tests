package main

import (
	"crypto/x509"
	"crypto/tls"
	"net/http"
	"log"
	"io/ioutil"
	"fmt"
	"path/filepath"
)

func handleError(err error) {
	if err != nil {
		log.Fatal("Fatal", err)
	}
}


func main() {
	absPathClientCrt, err := filepath.Abs("certs/client.crt")
	handleError(err)
	absPathClientKey, err := filepath.Abs("certs/client.key")
	handleError(err)
	absPathServerCrt, err := filepath.Abs("certs/server.crt")
	handleError(err)

	cert, err := tls.LoadX509KeyPair(absPathClientCrt, absPathClientKey)
	if err != nil {
		log.Fatalln("Unable to load cert", err)
	}

	roots := x509.NewCertPool()

	// We're going to load the server cert and add all the intermediates and CA from that.
	// Alternatively if we have the CA directly we could call AppendCertificate method
	fakeCA, err := ioutil.ReadFile(absPathServerCrt)
	ok := roots.AppendCertsFromPEM([]byte(fakeCA))
	if !ok {
		panic("failed to parse root certificate")
	}

	tlsConf := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs: roots,
	}
	tr := &http.Transport{TLSClientConfig: tlsConf}
	client := &http.Client{Transport: tr}

	resp, err := client.Get("https://mysite.local")
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(resp.Status)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(string(body))
}
