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

	controlLoop(f5)
}

func getPasswd() []byte {
	fmt.Print("Password:")

	// This supresses terminal echo when typing
	pass, _ := terminal.ReadPassword(int(syscall.Stdin))

	return pass
}

func getUser() string {
	fmt.Print("username:")
	var username string
	_, err := fmt.Scan(&username)
	if err != nil {
		log.Println(err)
	}
	return username
}

func getConsoleURI() string {
	fmt.Print("BigIP F5 Management URI:")
	var uri string
	_, err := fmt.Scan(&uri)
	if err != nil {
		log.Println(err)
	}
	return uri
}

func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func establishSession(uri string, user string, pass string) (f5 *bigip.BigIP) {
	f5 = bigip.NewSession(uri, user, pass, &bigip.ConfigOptions{APICallTimeout: 30 * time.Second})
	return f5
}

func listNodes(f5 *bigip.BigIP) {
	nodes, err := f5.Nodes()
	if err != nil {
		fmt.Println(err)
	}

	for cnt, node := range nodes.Nodes {
		fmt.Printf("%d\t{part:%s, name:%s, address:%s, state:%s}\n", cnt, node.Partition, node.Name, node.Address, node.State)
	}
}

func listPools(f5 *bigip.BigIP) {
	pools, err := f5.Pools()

	if err != nil {
		fmt.Println(err)
	}

	for _, pool := range pools.Pools {
		members, err := f5.PoolMembers("/" + pool.Partition + "/" + pool.Name)

		if err != nil {
			fmt.Println(err)
		}

		for memcnt, member := range members.PoolMembers {
			fmt.Printf("{%d, pool:%s, member_partition:%s, member_name:%s, member_ip:%s, member_state:%s}\n", memcnt, pool.Name, member.Partition, member.Name, member.Address, member.State)
		}
	}
}

func togglePoolMembers(f5 *bigip.BigIP) {
	members, err := f5.PoolMembers("/" + "{partition}" + "/" + "{poolname}")
	if err != nil {
		fmt.Println(err)
	}

	mems := make([]bigip.PoolMember, len(members.PoolMembers))
	for memcnt, member := range members.PoolMembers {
		fmt.Printf("{%d, member_partition:%s, member_name:%s, member_ip:%s, member_state:%s}\n", memcnt, member.Partition, member.Name, member.Address, member.State)
		member.Session = "user-disabled"
		mems[memcnt] = member
		fmt.Printf("%v", member)
	}

	error := f5.UpdatePoolMembers("/"+"{partition}"+"/"+"{poolname}", &mems)
	if error != nil {
		fmt.Println(error)
	}
}

func controlLoop(f5 *bigip.BigIP) {
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
			listNodes(f5)
			break
		case "listpools":
			listPools(f5)
			break
		case "toggle":
			togglePoolMembers(f5)
			break
		}
	}
}
