package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"code.google.com/p/go.crypto/ssh"
)

var (
	keyfile = flag.String("keyfile", filepath.Join(os.Getenv("HOME"), ".ssh/id_rsa"), "path to ssh private key file")
	user    = flag.String("user", "rcarlsen", "condor (and via node) ssh username")
	via     = flag.String("via", "best-tux.cae.wisc.edu:22", "intermediate server (if needed) prior to submit node ssh")
	dst     = flag.String("dst", "submit-3.chtc.wisc.edu:22", "condor submit node URI")
	submit  = flag.String("submit", "", "submit file to send")
)

func main() {
	flag.Usage = func() {
		fmt.Println("Usage: condor [FILE...]")
		fmt.Println("Copy listed files to condor submit node and possibly submit a job.\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	config := buildConfig()

	cae, err := ssh.Dial("tcp", *via, config)
	if err != nil {
		log.Fatal(err)
	}

	chtc, err := Hop(cae, *dst, config)
	if err != nil {
		log.Fatal(err)
	}

	fnames := flag.Args()
	if *submit != "" {
		fnames = append(fnames, *submit)
	}

	for _, fname := range flag.Args() {
		f, err := os.Open(fname)
		if err != nil {
			log.Fatal(err)
		}
		err = copyFile(chtc, f, fname)
		if err != nil {
			log.Fatal(err)
		}
		f.Close()
	}

	if *submit != "" {
		out, err := combined(chtc, "condor_submit "+*submit)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", out)
	}
}

func copyFile(c *ssh.Client, r io.Reader, path string) error {
	s, err := c.NewSession()
	if err != nil {
		return err
	}
	defer s.Close()

	w, err := s.StdinPipe()
	if err != nil {
		return err
	}

	s.Start("tee " + path)

	_, err = io.Copy(w, r)
	if err != nil {
		return err
	}
	w.Close()

	return s.Wait()
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
		User: *user,
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
