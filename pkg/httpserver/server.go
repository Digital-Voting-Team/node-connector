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

// Server represents the HTTP server.
type Server struct {
	Nodes *node.Nodes
	Echo  *echo.Echo
}

func NewServer() *Server {
	s := &Server{
		Nodes: &node.Nodes{
			Nodes:    []*node.Node{},
			NodesMap: make(map[string]struct{}),
		},
		Echo: echo.New(),
	}
	s.Echo.HideBanner = true

	s.InitRouters()

	// broadcast each 5 seconds
	go func() {
		for {
			s.broadcast()
			time.Sleep(5 * time.Second)
		}
	}()

	return s
}

// AddNodeHandler handles the route for adding a new node.
func (s *Server) AddNodeHandler(c echo.Context) error {
	currentNode := new(node.Node)
	if err := c.Bind(currentNode); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	currentNode.LastResponse = time.Now()

	err := s.Nodes.AddNode(currentNode)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// log the new currentNode
	log.Printf("New currentNode added: %s", currentNode.IP)

	return c.JSON(http.StatusCreated, currentNode)
}

// ListNodesHandler handles the route for listing all nodes.
func (s *Server) ListNodesHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{"ip_list": s.Nodes.Nodes})
}

// InitRouters initializes the route handlers.
func (s *Server) InitRouters() {
	s.Echo.POST("/nodes", s.AddNodeHandler)
	s.Echo.GET("/nodes", s.ListNodesHandler)
}

// broadcast sends list of nodes to all nodes.
func (s *Server) broadcast() {
	for _, currentNode := range s.Nodes.Nodes {
		go func(node *node.Node) {
			u := url.URL{Scheme: "wss", Host: node.IP, Path: "/ws"}
			log.Printf("connecting to %s", u.String())

			// Establish a WebSocket connection
			conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
			if err != nil {
				log.Println("dial:", err)
				return
			}

			defer func(conn *websocket.Conn) {
				err := conn.Close()
				if err != nil {
					log.Println("close:", err)
				}
			}(conn)

			// Send nodes list to the currentNode in JSON format
			err = conn.WriteJSON(s.Nodes.Nodes)
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
		}(currentNode)
	}
}
