package main

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"net"
	"sync"
	"time"
)

type ServerOptions struct {
	Host string
}

type Server struct {
	listener    net.Listener
	accepted    chan net.Conn
	connMu      sync.Mutex
	connections map[int64]net.Conn
	conn2close  chan int64
	messages    chan string
}

func NewServer(_ ServerOptions) *Server {
	return &Server{
		accepted:    make(chan net.Conn),
		connections: make(map[int64]net.Conn),
		conn2close:  make(chan int64),
		messages:    make(chan string),
	}
}

func (s *Server) Run(host string) error {
	var err error
	s.listener, err = net.Listen("tcp", host)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	go s.accepter()
	go s.manager()
	go s.messanger()

	return nil
}

func (s *Server) accepter() {
	slog.Debug("running", slog.String("method", "accepter"))
	defer slog.Debug("ended", slog.String("method", "accepter"))

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			slog.Error("listener.accept", slog.Any("error", err))
			return
		}
		s.accepted <- conn
	}
}

func (s *Server) manager() {
	slog.Debug("running", slog.String("method", "manager"))
	defer slog.Debug("ended", slog.String("method", "manager"))

	for {
		select {
		case conn := <-s.accepted:
			id := rand.Int63()
			s.addConnection(id, conn)
			go s.scanConnection(id, conn)
		case id := <-s.conn2close:
			s.closeConnection(id)
		}
	}
}

func (s *Server) messanger() {
	slog.Debug("running", slog.String("method", "messanger"))
	defer slog.Debug("ended", slog.String("method", "messanger"))

	for message := range s.messages {
		if len(s.connections) == 0 {
			continue
		}

		s.connMu.Lock()
		for id, conn := range s.connections {
			if _, err := conn.Write(formatMessage(message)); err != nil {
				s.conn2close <- id
				continue
			}
		}
		s.connMu.Unlock()
	}
}

func formatMessage(msg string) []byte {
	return []byte(fmt.Sprintf("%s %s\n", time.Now().Format(time.TimeOnly), msg))
}

func (s *Server) scanConnection(id int64, conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		s.messages <- scanner.Text()
	}
	s.conn2close <- id
}

func (s *Server) addConnection(id int64, conn net.Conn) {
	s.connMu.Lock()
	defer s.connMu.Unlock()

	s.connections[id] = conn
	slog.Info("add_connection", slog.Int64("id", id), slog.String("remote_address", conn.RemoteAddr().String()))
}

func (s *Server) closeConnection(id int64) {
	s.connMu.Lock()
	defer s.connMu.Unlock()

	conn, ok := s.connections[id]
	if !ok {
		return
	}
	_ = conn.Close()
	delete(s.connections, id)
	slog.Info("close_connection", slog.Int64("id", id))
}
