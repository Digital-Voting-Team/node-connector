package node

import (
	"fmt"
	"time"
)

// Node represents a single node with its IP address and last response time.
type Node struct {
	IP           string    `json:"ip"`
	LastResponse time.Time `json:"-"`
}

type Nodes struct {
	Nodes    []*Node
	NodesMap map[string]struct{}
}

// GetNodesIPs returns list of nodes ips
func (n *Nodes) GetNodesIPs() []string {
	var ips []string
	for _, node := range n.Nodes {
		ips = append(ips, node.IP)
	}
	return ips
}

// AddNode add node to nodes if not exists
func (n *Nodes) AddNode(node *Node) error {
	if _, ok := n.NodesMap[node.IP]; ok {
		return fmt.Errorf("node %s already exists", node.IP)
	}

	n.Nodes = append(n.Nodes, node)
	n.NodesMap[node.IP] = struct{}{}

	return nil
}
