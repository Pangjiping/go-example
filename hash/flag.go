package hash

import "flag"

var (
	keysPtr    = flag.Int("keys", 10000000, "key number")
	nodePtr    = flag.Int("nodes", 3, "node number of old cluster")
	vnodesPtr  = flag.Int("vnode", 100, "node number of new cluster")
	newNodePtr = flag.Int("new-nodes", 4, "node number of new cluster")
)
