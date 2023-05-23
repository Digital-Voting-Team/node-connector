package httpserver

import (
	"github.com/Digital-Voting-Team/node-connector/pkg/node"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	broadcastTime = 5 * time.Second
	nodeTimeout   = 10 * time.Second
	pingInterval  = 15 * time.Second
)

// Server represents the HTTP server.
type Server struct {
	Nodes *node.Nodes
	Echo  *echo.Echo
}

func NewServer() *Server {
	s := &Server{
		Nodes: &node.Nodes{
			NodesMap: make(map[string]*node.Node),
		},
		Echo: echo.New(),
	}
	s.Echo.HideBanner = true

	s.InitRouters()

	//broadcast each 5 seconds
	go func() {
		for {
			s.broadcast()
			time.Sleep(broadcastTime)
		}
	}()

	// ping each 5 seconds, remove inactive nodes after 1 hour
	go func() {
		ticker := time.NewTicker(pingInterval)

		for {
			select {
			case <-ticker.C:
				log.Println("Removing inactive nodes")
				s.Nodes.RemoveInactiveNodes(nodeTimeout)
			default:
				time.Sleep(broadcastTime)
				s.Ping()
			}
		}
	}()

	return s
}

// AddNodeHandler handles the route for adding a new node.
func (s *Server) AddNodeHandler(c echo.Context) error {
	req := &node.Node{}
	err := c.Bind(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = s.Nodes.AddNode(req.Hostname, req.ValidatorKey)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// log the new currentNode
	log.Printf("New currentNode added: %s", req.Hostname)

	return c.JSON(http.StatusCreated, s.Nodes.NodesMap[req.Hostname])
}

// ListNodesHandler handles the route for listing all nodes.
func (s *Server) ListNodesHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{"node_list": s.Nodes.GetNodeList()})
}

// InitRouters initializes the route handlers.
func (s *Server) InitRouters() {
	s.Echo.POST("/nodes", s.AddNodeHandler)
	s.Echo.GET("/nodes", s.ListNodesHandler)
}

// broadcast sends list of nodes to all nodes.
func (s *Server) broadcast() {
	for _, currentNode := range s.Nodes.NodesMap {
		go func(node *node.Node) {
			u := url.URL{Scheme: "ws", Host: node.Hostname, Path: "/ws"}
			log.Printf("connecting to %s", u.String())

			// Establish a WebSocket connection
			conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
			if err != nil {
				log.Println("dial:", err)
				return
			}

			defer func(conn *websocket.Conn) {
				log.Println("Closing connection")
				err := conn.Close()
				if err != nil {
					log.Println("close:", err)
				}
			}(conn)

			// Send nodes list to the currentNode in JSON format
			err = conn.WriteJSON(s.Nodes.GetNodeList())
			if err != nil {
				log.Println("write:", err)
				return
			}

			//close connection
			err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
		}(currentNode)
	}
}

// Ping sends a ping to all nodes, wait for pong for 5 seconds and update time when answer received.
func (s *Server) Ping() {
	for _, currentNode := range s.Nodes.NodesMap {
		go func(node *node.Node) {
			u := url.URL{Scheme: "ws", Host: node.Hostname, Path: "/ping"}
			log.Printf("connecting to %s", u.String())

			// Establish a WebSocket connection
			conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
			if err != nil {
				log.Println("dial:", err)
				return
			}

			conn.SetPongHandler(func(string) error {
				println("Received pong")
				node.LastResponse = time.Now()
				return nil
			})

			defer func(conn *websocket.Conn) {
				err := conn.Close()
				if err != nil {
					log.Println("close:", err)
				}
			}(conn)

			// Send ping to the currentNode in JSON format
			println("Sending ping")

			err = conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(15*time.Second))
			if err != nil {
				log.Println("write:", err)
				return
			}

			// close connection
			err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}

			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				return
			}
			log.Println("message received:", string(message))
		}(currentNode)
	}
}
