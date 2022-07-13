package hash

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	"stathat.com/c/consistent"
)

func ConsistentHash() {
	flag.Parse()
	var keys = *keysPtr
	var nodes = *nodePtr
	var newNodes = *newNodePtr

	c := consistent.New()
	for i := 0; i < nodes; i++ {
		c.Add(strconv.Itoa(i))
	}

	newC := consistent.New()
	for i := 0; i < newNodes; i++ {
		newC.Add(strconv.Itoa(i))
	}

	migrate := 0
	for i := 0; i < keys; i++ {
		server, err := c.Get(strconv.Itoa(i))
		if err != nil {
			log.Fatal(err)
		}

		newServer, err := newC.Get(strconv.Itoa(i))
		if err != nil {
			log.Fatal(err)
		}

		if server != newServer {
			migrate++
		}
	}

	migrateRatio := float64(migrate) / float64(keys)
	fmt.Printf("%f%%\n", migrateRatio*100)
}
