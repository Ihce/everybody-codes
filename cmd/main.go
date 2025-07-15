package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"everybody-codes/events/kingdom_algorithmia"
)

type Quest interface {
	Solve(input io.Reader) (string, error)
}

var quests = make(map[string]map[int]Quest)

func init() {
	// Register Kingdom of Algorithmia
	registerEvent("kingdom_algorithmia", kingdom_algorithmia.GetQuests())

	// Add more events here:
	// registerEvent("song_ducks_dragons", song_ducks_dragons.GetQuests())
}

func registerEvent(name string, eventQuests map[int]Quest) {
	quests[name] = eventQuests
}

func main() {
	var (
		event     = flag.String("event", "", "Event name (e.g., kingdom_algorithmia)")
		questID   = flag.String("quest", "", "Quest ID")
		inputFile = flag.String("input", "", "Input file path (optional, uses stdin if not provided)")
		list      = flag.Bool("list", false, "List available events and quests")
	)
	flag.Parse()

	if *list {
		printAvailable()
		return
	}

	if *event == "" || *questID == "" {
		fmt.Println("Usage:")
		fmt.Println("  -event string    Event name")
		fmt.Println("  -quest string    Quest ID")
		fmt.Println("  -input string    Input file (optional)")
		fmt.Println("  -list           List available content")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  go run main.go -event kingdom_algorithmia -quest 1 -input testdata/kingdom_algorithmia/quest01_input.txt")
		fmt.Println("  echo 'test' | go run main.go -event kingdom_algorithmia -quest 1")
		os.Exit(1)
	}

	qID, err := strconv.Atoi(*questID)
	if err != nil {
		log.Fatalf("Invalid quest ID: %v", err)
	}

	// Get quest
	eventQuests, exists := quests[*event]
	if !exists {
		log.Fatalf("Event '%s' not found", *event)
	}

	quest, exists := eventQuests[qID]
	if !exists {
		log.Fatalf("Quest %d not found in event '%s'", qID, *event)
	}

	// Get input
	var input *os.File
	if *inputFile != "" {
		input, err = os.Open(*inputFile)
		if err != nil {
			log.Fatalf("Failed to open input file: %v", err)
		}
		defer input.Close()
	} else {
		// Try to find matching .txt file
		defaultPath := fmt.Sprintf("events/%s/quest%02d.txt", *event, qID)
		if _, err := os.Stat(defaultPath); err == nil {
			input, err = os.Open(defaultPath)
			if err != nil {
				log.Fatalf("Failed to open default input file %s: %v", defaultPath, err)
			}
			defer input.Close()
			fmt.Printf("Using default input file: %s\n", defaultPath)
		} else {
			input = os.Stdin
		}
	}

	// Run quest
	result, err := quest.Solve(input)
	if err != nil {
		log.Fatalf("Quest failed: %v", err)
	}

	fmt.Println(result)
}

func printAvailable() {
	fmt.Println("Available Events and Quests:")
	for event, eventQuests := range quests {
		fmt.Printf("\n%s:\n", event)
		for questID := range eventQuests {
			fmt.Printf("  Quest %d\n", questID)
		}
	}
}
