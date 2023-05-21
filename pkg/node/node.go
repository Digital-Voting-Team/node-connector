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
	NodesMap map[string]*Node
}

// AddNode add node to nodes if not exists
func (n *Nodes) AddNode(ip string) error {
	if _, ok := n.NodesMap[ip]; ok {
		return fmt.Errorf("node %s already exists", ip)
	}

	n.NodesMap[ip] = &Node{
		IP:           ip,
		LastResponse: time.Now(),
	}

	return nil
}

// RemoveNode remove node from nodes if exists
func (n *Nodes) RemoveNode(ip string) error {
	if _, ok := n.NodesMap[ip]; !ok {
		return fmt.Errorf("node %s does not exist", ip)
	}

	delete(n.NodesMap, ip)

	return nil
}

// GetNodeList get node list from nodes
func (n *Nodes) GetNodeList() []*Node {
	var nodeList = []*Node{}

	for ip := range n.NodesMap {
		nodeList = append(nodeList, n.NodesMap[ip])
	}

	return nodeList
}
