package server

import (
	"bytes"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

// HandleFunc ...
type HandleFunc func(conn net.Conn)

// Server ...
type Server struct {
	addr     string
	mu       sync.RWMutex
	handlers map[string]HandleFunc
}

// NewServer ....
func NewServer(addr string) *Server {
	return &Server{addr: addr, handlers: make(map[string]HandleFunc)}
}

// Register .....
func (s *Server) Register(path string, handler HandleFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[path] = handler
}

// Start ....
func (s *Server) Start() error {

	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Print(err)
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		err = s.handle(conn)
		if err != nil {
			log.Print(err)
			// Идём обслуживать следующего
			continue
		}

	}

}

func (s *Server) handle(conn net.Conn) (err error) {
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			if err == nil {
				err = cerr
				return
			}
			log.Print(err)
		}
	}()
	// .....

	//conn.Write([]byte("Hello!\r\n"))

	buf := make([]byte, 4096)

	n, err := conn.Read(buf)
	if err == io.EOF {
		log.Printf("%s", buf[:n])
		return nil
	}

	if err != nil {
		log.Print(err)
		return err
	}
	//log.Printf("%s", buf[:n])

	data := buf[:n]
	requestLineDelim := []byte{'\r', '\n'}
	requestLineEnd := bytes.Index(data, requestLineDelim)
	if requestLineEnd == -1 {
	}

	requestLine := string(data[:requestLineEnd])
	parts := strings.Split(requestLine, " ")
	if len(parts) != 3 {
	}

	method, path, version := parts[0], parts[1], parts[2]
	if method != "GET" {

	}

	if version != "HTTP/1.1" {

	}

	s.mu.RLock()
	for name, handler := range s.handlers {
		if name == path {
			s.mu.RUnlock()
			handler(conn)
			break
		}
	}
	return nil

}
