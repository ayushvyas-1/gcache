package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ayushvyas-1/gcache/internal/cache"
)

func main() {

	var (
		mode        = flag.String("mode", "server", "Mode: 'server' or 'client'")
		address     = flag.String("addr", "localhost:8080", "Server address")
		capacity    = flag.Int("capacity", 1000, "Cache capacity (server mode only)")
		interactive = flag.Bool("interactive", false, "Interactive client mode")
		command     = flag.String("cmd", "", "Single command to execute (client mode)")
	)
	flag.Parse()

	switch *mode {
	case "server":
		runServer(*address, *capacity)
	case "client":
		runClient(*address, *interactive, *command)
	default:
		fmt.Fprintf(os.Stderr, "Invalid mode: %s. Use 'server' or 'client'\n", *mode)
		os.Exit(1)
	}
}

func runClient(address string, Interactive bool, command string) {
	client, err := cache.NewClient(address)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer client.Close()

	if Interactive {
		client.InteractiveMode()
	} else if command != "" {
		response, err := client.SendCommand(command)
		if err != nil {
			log.Fatalf("Command failed: %v", err)
		}
		fmt.Print(response)
	} else {
		runClientDemo(client)
	}
}

func runServer(address string, capacity int) {
	fmt.Printf("Starting GCache Server...\n")
	fmt.Printf("Address: %s\n", address)
	fmt.Printf("Capacity: %d\n", capacity)

	server := cache.NewServer(address, capacity)
	if err := server.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func runClientDemo(client *cache.Client) {
	fmt.Println("=== GCache Client Demo ===")

	// Test ping
	fmt.Print("Testing connection... ")
	if err := client.Ping(); err != nil {
		fmt.Printf("FAILED: %v\n", err)
		return
	}
	fmt.Println("OK")

	// Set some values
	fmt.Println("\nSetting values:")
	testData := map[string]string{
		"name":     "GCache",
		"version":  "1.0",
		"author":   "You",
		"language": "Go",
	}

	for key, value := range testData {
		fmt.Printf("  SET %s %s... ", key, value)
		if err := client.Set(key, value); err != nil {
			fmt.Printf("FAILED: %v\n", err)
		} else {
			fmt.Println("OK")
		}
	}

	// Get values back
	fmt.Println("\nGetting values:")
	for key := range testData {
		fmt.Printf("  GET %s... ", key)
		if value, err := client.Get(key); err != nil {
			fmt.Printf("FAILED: %v\n", err)
		} else {
			fmt.Printf("OK: %s\n", value)
		}
	}

	// Check size
	fmt.Print("\nChecking cache size... ")
	if size, err := client.Size(); err != nil {
		fmt.Printf("FAILED: %v\n", err)
	} else {
		fmt.Printf("OK: %d items\n", size)
	}

	// Test delete
	fmt.Print("Deleting 'version'... ")
	if err := client.Delete("version"); err != nil {
		fmt.Printf("FAILED: %v\n", err)
	} else {
		fmt.Println("OK")
	}

	// Check size again
	fmt.Print("Checking size after delete... ")
	if size, err := client.Size(); err != nil {
		fmt.Printf("FAILED: %v\n", err)
	} else {
		fmt.Printf("OK: %d items\n", size)
	}

	// Get server info
	fmt.Print("Getting server info... ")
	if info, err := client.Info(); err != nil {
		fmt.Printf("FAILED: %v\n", err)
	} else {
		fmt.Printf("OK:\n%s\n", info)
	}

	fmt.Println("\n=== Demo Complete ===")
	fmt.Printf("Try interactive mode: %s -mode=client -interactive\n", os.Args[0])
}
