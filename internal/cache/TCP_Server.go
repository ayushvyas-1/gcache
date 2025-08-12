package cache

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type Server struct {
	cache    *LRUCache
	listener net.Listener
	address  string
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewServer(address string, cacheCapacity int) *Server {
	ctx, cancel := context.WithCancel(context.Background())

	return &Server{
		cache:   NewLRUCache(cacheCapacity),
		address: address,
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (s *Server) Start() error {
	var err error
	s.listener, err = net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}

	log.Printf("GCache server started on %s", s.address)
	log.Printf("Cache capacity: %d", s.cache.capacity)

	// Handle graceful shutdown
	go s.handleShutdown()

	// Accept connections
	for {
		select {
		case <-s.ctx.Done():
			return nil
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				if s.ctx.Err() != nil {
					return nil // Server is shutting down
				}
				log.Printf("Failed to accept connection: %v", err)
				continue
			}

			// Handle connection in goroutine
			go s.handleConnection(conn)
		}
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	log.Printf("Client connected: %s", clientAddr)
	defer log.Printf("Client disconnected: %s", clientAddr)

	scanner := bufio.NewScanner(conn)
	writer := bufio.NewWriter(conn)

	for scanner.Scan() {
		select {
		case <-s.ctx.Done():
			return
		default:
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}

			response := s.processCommand(line)

			if _, err := writer.WriteString(response + "\r\n"); err != nil {
				log.Printf("Error writing to client %s: %v", clientAddr, err)
				return
			}

			if err := writer.Flush(); err != nil {
				log.Printf("Error flushing to client %s: %v", clientAddr, err)
				return
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Connection error with %s: %v", clientAddr, err)
	}
}

func (s *Server) processCommand(command string) string {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return "-ERR empty command"
	}

	cmd := strings.ToUpper(parts[0])

	switch cmd {
	case "GET":
		return s.handleGet(parts)
	case "SET":
		return s.handleSet(parts)
	case "DEL":
		return s.handleDel(parts)
	case "SIZE":
		return s.handleSize(parts)
	case "CLEAR":
		return s.handleClear(parts)
	case "PING":
		return s.handlePing(parts)
	case "INFO":
		return s.handleInfo(parts)
	case "STATS":
		return s.handleStats(parts)
	case "QUIT":
		return s.handleQuit(parts)
	default:
		return fmt.Sprintf("-ERR unknown command '%s'", cmd)
	}
}

func (s *Server) handleGet(parts []string) string {
	if len(parts) != 2 {
		return "-ERR wrong number of arguments for 'GET' command"
	}

	key := parts[1]
	if value, exists := s.cache.Get(key); exists {
		return fmt.Sprintf("+%s", value)
	}
	return "-ERR key not found"
}

func (s *Server) handleSet(parts []string) string {
	if len(parts) < 3 {
		return "-ERR wrong number of arguments for 'SET' command"
	}

	key := parts[1]
	// Join remaining parts as value (allows spaces in values)
	value := strings.Join(parts[2:], " ")

	s.cache.Put(key, value)
	return "+OK"
}

func (s *Server) handleDel(parts []string) string {
	if len(parts) != 2 {
		return "-ERR wrong number of arguments for 'DEL' command"
	}

	key := parts[1]
	if s.cache.Delete(key) {
		return "+OK"
	}
	return "-ERR key not found"
}

func (s *Server) handleSize(parts []string) string {
	if len(parts) != 1 {
		return "-ERR wrong number of arguments for 'SIZE' command"
	}

	return fmt.Sprintf(":%d", s.cache.Size())
}

func (s *Server) handleClear(parts []string) string {
	if len(parts) != 1 {
		return "-ERR wrong number of arguments for 'CLEAR' command"
	}

	s.cache.Clear()
	return "+OK"
}

func (s *Server) handlePing(parts []string) string {
	if len(parts) == 1 {
		return "+PONG"
	} else if len(parts) == 2 {
		return fmt.Sprintf("+%s", parts[1])
	}
	return "-ERR wrong number of arguments for 'PING' command"
}

func (s *Server) handleInfo(parts []string) string {
	if len(parts) != 1 {
		return "-ERR wrong number of arguments for 'INFO' command"
	}

	// Single line info to avoid parsing issues
	info := fmt.Sprintf("gcache_version:1.0 cache_capacity:%d cache_size:%d uptime_seconds:%.0f",
		s.cache.capacity,
		s.cache.Size(),
		time.Since(startTime).Seconds())

	return fmt.Sprintf("+%s", info)
}

func (s *Server) handleStats(parts []string) string {
	if len(parts) != 1 {
		return "-ERR wrong number of arguments for 'STATS' command"
	}

	// Single line stats to avoid multi-line parsing issues
	stats := fmt.Sprintf("size:%d capacity:%d",
		s.cache.Size(), s.cache.capacity)

	return fmt.Sprintf("+%s", stats)
}

func (s *Server) handleQuit(parts []string) string {
	return "+BYE"
}

func (s *Server) handleShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Printf("Received signal %v, shutting down gracefully...", sig)

	s.Stop()
}

func (s *Server) Stop() {
	s.cancel()
	if s.listener != nil {
		s.listener.Close()
	}
	log.Println("Server stopped")
}

var startTime time.Time

func init() {
	startTime = time.Now()
}
