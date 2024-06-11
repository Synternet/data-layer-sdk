package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

const consoleTemplate = `User NKEY: %s
User JWT: %s
`

const credsTemplate = `-----BEGIN NATS USER JWT-----
%s
------END NATS USER JWT------

************************* IMPORTANT *************************
NKEY Seed printed below can be used to sign and prove identity.
NKEYs are sensitive and should be treated as secrets.

-----BEGIN USER NKEY SEED-----
%s
------END USER NKEY SEED------

*************************************************************
`

func main() {
	creds := flag.Bool("creds", false, "Generates Creds file contents instead.")
	h := flag.Bool("h", false, "Display this help.")
	jsManager := flag.Bool("js", false, "Enables JetStream Manager option.")
	flag.Parse()

	if *h {
		flag.Usage()
		return
	}

	inStat, err := os.Stdin.Stat()
	if err != nil || inStat.Size() == 0 {
		fmt.Print("Enter Seed(Account NKEY): ")
	}

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		fatal("Seed scanning failed")
	}

	seed := []byte(strings.TrimSpace(scanner.Text()))
	if len(seed) == 0 {
		fatal("Account seed is empty")
	}

	var opts []Opt
	if *jsManager {
		opts = append(opts, JetStreamManagerOpt)
	}
	userSeed, jwt, err := createUser(seed, opts)
	if err != nil {
		fatal("Failed to create user: %s", err)
	}

	template := consoleTemplate
	if *creds {
		template = credsTemplate
	}

	fmt.Printf(template, jwt, userSeed)
}

func fatal(template string, values ...interface{}) {
	fmt.Printf(template, values...)
	fmt.Println()
	os.Exit(1)
}
