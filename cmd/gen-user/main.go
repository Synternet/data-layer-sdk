package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/synternet/data-layer-sdk/pkg/user"
)

const consoleTemplate = `User JWT: %s
User NKEY: %s
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
	credsFormat := flag.Bool("creds", false, "Generates Creds file contents instead.")
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

	var opts []user.Opt
	if *jsManager {
		opts = append(opts, user.JetStreamManagerOpt)
	}
	userSeed, jwt, err := user.CreateCreds(seed, opts...)
	if err != nil {
		fatal("Failed to create user: %s", err)
	}

	template := consoleTemplate
	if *credsFormat {
		template = credsTemplate
	}

	fmt.Printf(template, jwt, userSeed)
}

func fatal(template string, values ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(template, values...))
	os.Exit(1)
}
