package node

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

// Node represents a single node with its Hostname address and last response time.
type Node struct {
	Hostname     string    `json:"hostname"`
	ValidatorKey [33]byte  `json:"validator_key"`
	LastResponse time.Time `json:"-"`
}

type Nodes struct {
	nodes map[string]*Node
	mutex sync.Mutex
}

// NewNodes create new nodes
func NewNodes() *Nodes {
	nodes := &Nodes{
		nodes: make(map[string]*Node),
	}

	err := nodes.LoadNodes()
	if err != nil {
		log.Println("Error loading nodes:", err)
		return nil
	}

	// save nodes every 5 seconds
	go func() {
		for {
			time.Sleep(5 * time.Second)
			log.Println("Saving nodes")
			err = nodes.SaveNodes()
			if err != nil {
				log.Println("Error saving nodes:", err)
			}
		}
	}()

	return nodes
}

const file = "nodes.dat"

func (n *Nodes) SaveNodes() error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			log.Println("Error closing file:", err)
		}
	}(f)

	enc := gob.NewEncoder(f)
	for _, node := range n.GetNodeList() {
		if err = enc.Encode(node); err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func (n *Nodes) LoadNodes() error {
	// open or create file
	f, err := os.OpenFile(file, os.O_CREATE|os.O_RDONLY, 0666)
	if err != nil {
		return err
	}

	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			log.Println("Error closing file:", err)
		}
	}(f)

	// exit if file is empty
	fi, err := f.Stat()
	if err != nil {
		return err
	}

	if fi.Size() == 0 {
		return nil
	}

	dec := gob.NewDecoder(f)
	for {
		node := &Node{}

		if err = dec.Decode(&node); err != nil {
			if err.Error() == "EOF" {
				break
			}
			return err
		}

		err = n.AddNode(node.Hostname, node.ValidatorKey)
		if err != nil {
			return err
		}
	}

	return nil
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

	return n.SaveNodes()
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
