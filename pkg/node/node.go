package node

import (
	"fmt"
	"log"
	"time"
)

// Node represents a single node with its Hostname address and last response time.
type Node struct {
	Hostname     string    `json:"hostname"`
	ValidatorKey [33]byte  `json:"validator_key"`
	LastResponse time.Time `json:"-"`
}

// TODO : add mutex for thread safety
type Nodes struct {
	NodesMap map[string]*Node
}

// AddNode add node to nodes with all fields if not exists
func (n *Nodes) AddNode(hostname string, validatorKey [33]byte) error {
	if _, ok := n.NodesMap[hostname]; ok {
		return fmt.Errorf("node %s already exists", hostname)
	}

	n.NodesMap[hostname] = &Node{
		Hostname:     hostname,
		ValidatorKey: validatorKey,
		LastResponse: time.Now(),
	}

	return nil
}

// RemoveNode remove node from nodes if exists
func (n *Nodes) RemoveNode(hostname string) error {
	if _, ok := n.NodesMap[hostname]; !ok {
		return fmt.Errorf("node %s does not exist", hostname)
	}

	delete(n.NodesMap, hostname)

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

// Update node last response time by hostname
func (n *Nodes) Update(hostname string) error {
	if _, ok := n.NodesMap[hostname]; !ok {
		return fmt.Errorf("node %s does not exist", hostname)
	}

	n.NodesMap[hostname].LastResponse = time.Now()

	return nil
}

// RemoveInactiveNodes remove nodes that have not responded for 1 hour
func (n *Nodes) RemoveInactiveNodes(duration time.Duration) {
	for hostname, node := range n.NodesMap {
		log.Println("Node", hostname, "last response", time.Since(node.LastResponse), "ago")
		if time.Since(node.LastResponse) > duration {
			log.Println("Deleting node", hostname)
			delete(n.NodesMap, hostname)
		}
	}
}
