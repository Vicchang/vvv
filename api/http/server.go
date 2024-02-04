package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	ln     net.Listener
	server *http.Server
	router *gin.Engine

	// Bind address  for the server's listener.
	Addr string

	// TODO: move to service
	lbURL       string
	connWatcher *ConnectionWatcher
}

// NewServer returns a new instance of Server.
func NewServer(lburl string) *Server {
	connWatcher := &ConnectionWatcher{
		m: make(map[net.Conn]struct{}),
	}

	// Create a new server that wraps the net/http server & add a gorilla router.
	s := &Server{
		server: &http.Server{
			ConnState: connWatcher.OnStateChange,
		},
		router:      gin.New(),
		lbURL:       lburl,
		connWatcher: connWatcher,
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

	go s.signalLB()
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

	authorized := r.Group("/")
	{
		authorized.GET("heartbeat", s.heartbeat)
		authorized.POST("echo", EchoEndpoint)
	}
}

func (s *Server) signalLB() {
	endpointURL, err := url.JoinPath(s.lbURL, "endpoint")
	if err != nil {
		panic(err)
	}

	surl := s.URL()

	post := func() {
		data, err := json.Marshal(map[string]string{
			"url": surl,
		})
		if err != nil {
			fmt.Printf("marshal json error, %s\n", err)
			return
		}

		resp, err := http.Post(endpointURL, "application/json", bytes.NewReader(data))
		if err != nil {
			fmt.Printf("post endpoint, %s\n", err)
			return
		}

		fmt.Printf("endpoint response code %d\n", resp.StatusCode)
	}

	post()

	go func() {
		ticker := time.NewTicker(time.Second * 60)
		defer ticker.Stop()
		for range ticker.C {
			post()
		}
	}()
}

type ConnectionWatcher struct {
	// mu protects remaining fields
	mu sync.Mutex

	// open connections are keys in the map
	m map[net.Conn]struct{}
}

func (cw *ConnectionWatcher) OnStateChange(conn net.Conn, state http.ConnState) {
	switch state {
	case http.StateNew:
		cw.mu.Lock()
		if cw.m == nil {
			cw.m = make(map[net.Conn]struct{})
		}
		cw.m[conn] = struct{}{}
		cw.mu.Unlock()
	case http.StateHijacked, http.StateClosed:
		cw.mu.Lock()
		delete(cw.m, conn)
		cw.mu.Unlock()
	}
}

func (cw *ConnectionWatcher) Connections() int {
	return len(cw.m)
}
