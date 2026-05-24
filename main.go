package main

import (
	"fmt"
	"log"
	"os"

	"github.com/engram/engram/internal/config"
	"github.com/engram/engram/internal/store"
)

// Version is the current version of engram.
const Version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]

	// Load configuration from the .engram directory
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Initialize the chunk store
	s, err := store.New(cfg)
	if err != nil {
		log.Fatalf("failed to initialize store: %v", err)
	}
	defer s.Close()

	switch cmd {
	case "search":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "usage: engram search <query>")
			os.Exit(1)
		}
		query := os.Args[2]
		results, err := s.Search(query)
		if err != nil {
			log.Fatalf("search failed: %v", err)
		}
		for _, r := range results {
			fmt.Println(r)
		}

	case "ingest":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "usage: engram ingest <file>")
			os.Exit(1)
		}
		filePath := os.Args[2]
		if err := s.Ingest(filePath); err != nil {
			log.Fatalf("ingest failed: %v", err)
		}
		fmt.Printf("successfully ingested %s\n", filePath)

	case "version":
		fmt.Printf("engram version %s\n", Version)

	case "help":
		printUsage()

	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

// printUsage prints the CLI usage information.
func printUsage() {
	fmt.Printf(`engram - a memory and context store for AI-assisted development

Version: %s

Usage:
  engram <command> [arguments]

Commands:
  search <query>   Search stored memory chunks by query
  ingest <file>    Ingest a file into the memory store
  version          Print the current version
  help             Show this help message
`, Version)
}
