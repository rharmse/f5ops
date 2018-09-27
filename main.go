package main

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/scottdware/go-bigip"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	uri := getConsoleURI()
	user := getUser()
	pass := getPasswd()
	passStr := string(pass)

	fmt.Printf("\n")
	fmt.Printf("Connecting to %s, using %s/%s\n", uri, user, hashAndSalt(pass))

	f5 := establishSession(uri, user, passStr)

	fmt.Println("at prompt, type help if needed")

	var read = ""
	for {
		fmt.Printf("f5>")
		_, err := fmt.Scan(&read)
		if err != nil {
			log.Println(err)
		}

		switch read {
		case "quit":
			fmt.Println("Switch you later aligator.")
			os.Exit(0)
		case "nodes":
			nodes, err := f5.Nodes()
			if err != nil {
				fmt.Println(err)
			}

			cnt := 0
			for _, node := range nodes.Nodes {
				fmt.Printf("%d{part:%s, name:%s, address:%s, state:%s}\n", cnt, node.Partition, node.Name, node.Address, node.State)
			}
			break
		}
	}
}

func getPasswd() []byte {
	// Prompt the user to enter a password
	fmt.Print("Password:")

	// Read the users input
	pass, _ := terminal.ReadPassword(int(syscall.Stdin))

	// Return the users input as a byte slice which will save us
	// from having to do this conversion later on
	return pass
}

func getUser() string {

	// Prompt the user to enter a username
	fmt.Print("username:")

	// We will use this to store the users input
	var username string

	// Read the users input
	_, err := fmt.Scan(&username)
	if err != nil {
		log.Println(err)
	}
	return username
}

func getConsoleURI() string {

	// Prompt the user to enter a username
	fmt.Print("BigIP F5 Management URI:")

	// We will use this to store the users input
	var uri string

	// Read the users input
	_, err := fmt.Scan(&uri)
	if err != nil {
		log.Println(err)
	}
	return uri
}

func hashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func establishSession(uri string, user string, pass string) (f5 *bigip.BigIP) {

	//Establish our session to the BIG-IP
	f5 = bigip.NewSession(uri, user, pass, &bigip.ConfigOptions{APICallTimeout: 30 * time.Second})
	return f5
}
