package cache

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

type Client struct {
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

func NewClient(address string) (*Client, error) {
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to connect : %v", err)
	}

	return &Client{
		conn:   conn,
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
	}, nil
}

func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Client) SendCommand(command string) (string, error) {

	//send command to server with writer
	if _, err := c.writer.WriteString(command + "\r\n"); err != nil {
		return "", fmt.Errorf("failed to send command: %v", err)
	}
	if err := c.writer.Flush(); err != nil {
		return "", fmt.Errorf("failed to flush: %v", err)
	}

	//read response
	response, err := c.reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("Failed to read response: %v", err)
	}

	response = strings.TrimSpace(response)

	if strings.Contains(response, "\r\n") {
		// For now, just replace \r\n with spaces for single-line display
		response = strings.ReplaceAll(response, "\r\n", " ")
	}

	return response, nil

}

func (c *Client) Get(key string) (string, error) {
	response, err := c.SendCommand(fmt.Sprintf("GET %s", key))
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(response, "+") {
		return response[1:], nil
	} else if strings.HasPrefix(response, "-ERR") {
		return "", fmt.Errorf("%s", response[5:])
	}

	return "", fmt.Errorf("unexpected response : %s", response)
}

func (c *Client) Set(key, value string) error {
	response, err := c.SendCommand(fmt.Sprintf("SET %s %s", key, value))
	if err != nil {
		return err
	}

	if response == "+OK" {
		return nil
	}

	return fmt.Errorf("unexpected response: %s", response)
}

func (c *Client) Delete(key string) error {
	response, err := c.SendCommand(fmt.Sprintf("DEL %s", key))
	if err != nil {
		return err
	}

	if response == "+OK" {
		return nil
	} else if strings.HasPrefix(response, "-ERR") {
		return fmt.Errorf("%s", response[5:])
	}

	return fmt.Errorf("unexpected response: %s", response)
}

func (c *Client) Size() (int, error) {
	response, err := c.SendCommand("SIZE")
	if err != nil {
		return 0, err
	}

	if strings.HasPrefix(response, ":") {
		var size int
		if _, err := fmt.Sscanf(response[1:], "%d", &size); err != nil {
			return 0, fmt.Errorf("invalid size response: %s", response)
		}
		return size, nil
	}
	return 0, fmt.Errorf("unexpected response: %s", response)

}

func (c *Client) Ping() error {
	response, err := c.SendCommand("PING")
	if err != nil {
		return err
	}

	if response == "+PONG" {
		return nil
	}

	return fmt.Errorf("unexpected response: %s", response)
}

func (c *Client) Clear() error {
	response, err := c.SendCommand("CLEAR")
	if err != nil {
		return err
	}

	if response == "+OK" {
		return nil
	}

	return fmt.Errorf("unexpected response: %s", response)
}

func (c *Client) Info() (string, error) {
	response, err := c.SendCommand("INFO")
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(response, "+") {
		return response[1:], nil
	}

	return "", fmt.Errorf("unexpected response: %s", response)
}

func printHelp() {
	help := `
Available Commands:
  GET key          - Get value for key
  SET key value    - Set key to value
  DEL key          - Delete key
  SIZE             - Get cache size
  CLEAR            - Clear all items
  PING [message]   - Ping server
  INFO             - Server information
  STATS            - Cache statistics
  QUIT/EXIT        - Exit client

Examples:
  SET mykey "hello world"
  GET mykey
  DEL mykey
`
	fmt.Println(help)
}

func (c *Client) InteractiveMode() {
	fmt.Println("GCache Client - Interactive Mode")
	fmt.Println("Commands: GET, SET, DEL, SIZE, CLEAR, PING, INFO, STATS, QUIT")
	fmt.Println("Type 'help' for more information")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("gcache> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		if input == "help" {
			printHelp()
			continue
		}
		if strings.ToUpper(input) == "QUIT" || strings.ToUpper(input) == "EXIT" {
			break
		}
		response, err := c.SendCommand(input)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		// Pretty print response
		if strings.HasPrefix(response, "+") {
			fmt.Printf("OK: %s\n", response[1:])
		} else if strings.HasPrefix(response, "-ERR") {
			fmt.Printf("ERROR: %s\n", response[5:])
		} else if strings.HasPrefix(response, ":") {
			fmt.Printf("VALUE: %s\n", response[1:])
		} else {
			fmt.Printf("RESPONSE: %s\n", response)
		}
	}
}
