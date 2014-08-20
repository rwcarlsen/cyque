package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"code.google.com/p/go.crypto/ssh"
)

var keyfile = flag.String("keyfile", filepath.Join(os.Getenv("HOME"), ".ssh/id_rsa"), "path to ssh private key file")

func main() {
	flag.Parse()

	config := buildConfig()

	cae, err := ssh.Dial("tcp", "best-tux.cae.wisc.edu:22", config)
	if err != nil {
		log.Fatal(err)
	}

	chtc, err := Hop(cae, "submit-3.chtc.wisc.edu:22", config)
	if err != nil {
		log.Fatal(err)
	}

	data, err := combined(chtc, "condor_status")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", data)

	data, err = combined(chtc, "pwd")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", data)
}

func combined(c *ssh.Client, cmd string) ([]byte, error) {
	s, err := c.NewSession()
	if err != nil {
		return nil, err
	}
	defer s.Close()

	return s.CombinedOutput(cmd)
}

func buildConfig() *ssh.ClientConfig {
	// config for ssh connection
	pemBytes, err := ioutil.ReadFile(*keyfile)
	if err != nil {
		log.Fatal(err)
	}

	signer, err := ssh.ParsePrivateKey(pemBytes)
	if err != nil {
		log.Fatal(err)
	}

	return &ssh.ClientConfig{
		User: "rcarlsen",
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
	}
}

func Hop(through *ssh.Client, toaddr string, c *ssh.ClientConfig) (*ssh.Client, error) {
	hopconn, err := through.Dial("tcp", toaddr)
	if err != nil {
		return nil, err
	}

	conn, chans, reqs, err := ssh.NewClientConn(hopconn, toaddr, c)
	if err != nil {
		return nil, err
	}

	return ssh.NewClient(conn, chans, reqs), nil
}
