package node

import (
	"fmt"
	"log"
	"sync"
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
	nodes map[string]*Node
	mutex sync.Mutex
}

// NewNodes create new nodes
func NewNodes() *Nodes {
	return &Nodes{
		nodes: make(map[string]*Node),
		mutex: sync.Mutex{},
	}
}

// AddNode add node to nodes with all fields if not exists
func (n *Nodes) AddNode(hostname string, validatorKey [33]byte) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if _, ok := n.nodes[hostname]; ok {
		return fmt.Errorf("node %s already exists", hostname)
	}

	n.nodes[hostname] = &Node{
		Hostname:     hostname,
		ValidatorKey: validatorKey,
		LastResponse: time.Now(),
	}

	return nil
}

// RemoveNode remove node from nodes if exists
func (n *Nodes) RemoveNode(hostname string) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if _, ok := n.nodes[hostname]; !ok {
		return fmt.Errorf("node %s does not exist", hostname)
	}

	delete(n.nodes, hostname)

	return nil
}

// GetNodeList get node list from nodes
func (n *Nodes) GetNodeList() []*Node {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	var nodeList = []*Node{}

	for ip := range n.nodes {
		nodeList = append(nodeList, n.nodes[ip])
	}

	return nodeList
}

// Update node last response time by hostname
func (n *Nodes) Update(hostname string) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if _, ok := n.nodes[hostname]; !ok {
		return fmt.Errorf("node %s does not exist", hostname)
	}

	n.nodes[hostname].LastResponse = time.Now()

	return nil
}

// RemoveInactiveNodes remove nodes that have not responded for 1 hour
func (n *Nodes) RemoveInactiveNodes(duration time.Duration) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	for hostname, node := range n.nodes {
		log.Println("Node", hostname, "last response", time.Since(node.LastResponse), "ago")
		if time.Since(node.LastResponse) > duration {
			log.Println("Deleting node", hostname)
			delete(n.nodes, hostname)
		}
	}
}
