package node

import "time"

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
