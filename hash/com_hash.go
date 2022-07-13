package hash

import (
	"flag"
	"fmt"
)

func hash(key int, nodes int) int {
	return key % nodes
}

func ComHash() {
	flag.Parse()
	var keys = *keysPtr
	var oldNodes = *nodePtr
	var newNodes = *newNodePtr

	migrate := 0
	for i := 0; i < keys; i++ {
		if hash(i, oldNodes) != hash(i, newNodes) {
			migrate++
		}
	}

	migrateRatio := float64(migrate) / float64(keys)
	fmt.Printf("%f%%\n", migrateRatio*100)
}
