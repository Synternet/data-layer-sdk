package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/syntropynet/data-layer-sdk/pkg/options"
)

const credsTemplate = `-----BEGIN NATS USER JWT-----
%s
------END NATS USER JWT------

************************* IMPORTANT *************************
NKEY Seed printed below can be used to sign and prove identity.
NKEYs are sensitive and should be treated as secrets.

-----BEGIN USER NKEY SEED-----
%s
------END USER NKEY SEED------

*************************************************************`

func main() {
	creds := flag.Bool("creds", false, "Generate Creds file contents instead.")
	h := flag.Bool("h", false, "Display this help")
	flag.Parse()

	if *h {
		flag.Usage()
		return
	}

	for {
		fmt.Print("Enter Seed(Account NKEY): ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		seed := scanner.Text()

		if seed == "" {
			return
		}
		fmt.Println()

		nkey, jwt, err := options.CreateUser(seed)
		if err != nil {
			panic(err)
		}

		if *creds {
			fmt.Printf(credsTemplate, *jwt, *nkey)
			fmt.Println()
		} else {
			fmt.Println("User NKEY: ", *nkey)
			fmt.Println("User JWT: ", *jwt)
		}
		fmt.Println()
	}
}
