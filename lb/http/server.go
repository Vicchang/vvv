package http

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vvv/lb"
)

type Server struct {
	ln     net.Listener
	server *http.Server
	router *gin.Engine

	// Bind address  for the server's listener.
	Addr string

	// Servics used by the various HTTP routes.
	SelectService lb.SelectService
	PodService    lb.PodService
}

// NewServer returns a new instance of Server.
func NewServer() *Server {
	// Create a new server that wraps the net/http server & add a gorilla router.
	s := &Server{
		server: &http.Server{},
		router: gin.New(),
	}

	// Our router is wrapped by another function handler to perform some
	// middleware-like tasks that cannot be performed by actual middleware.
	// This includes changing route paths for JSON endpoints & overridding methods.
	s.server.Handler = http.HandlerFunc(s.serveHTTP)

	// Setup handling routes.
	s.setupRouter()

	return s
}

// Port returns the TCP port for the running server.
// This is useful in tests where we allocate a random port by using ":0".
func (s *Server) Port() int {
	if s.ln == nil {
		return 0
	}
	return s.ln.Addr().(*net.TCPAddr).Port
}

// URL returns the local base URL of the running server.
func (s *Server) URL() string {
	return fmt.Sprintf("http://localhost:%d", s.Port())
}

func (s *Server) serveHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// Open validates the server options and begins listening on the bind address.
func (s *Server) Open() (err error) {
	if s.ln, err = net.Listen("tcp", s.Addr); err != nil {
		return err
	}

	// Begin serving requests on the listener. We use Serve() instead of
	// ListenAndServe() because it allows us to check for listen errors (such
	// as trying to use an already open port) synchronously.
	go s.server.Serve(s.ln)

	return nil
}

func (s *Server) setupRouter() {
	r := s.router
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.NoRoute(s.handlPostRequest)

	authorized := r.Group("/")
	{
		authorized.POST("endpoint", s.handleEndpoint)
	}
}
