package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ksKao/dota/lib"
)

func main() {
	// Create a new parser instance from a file. Alternatively see NewParser([]byte)
	f, err := os.Open("replay.dem")
	if err != nil {
		log.Fatalf("unable to open file: %s", err)
	}
	defer f.Close()

	startTime, hooks, err := lib.ParseDemFile(f)
	if err != nil {
		log.Fatalf("failed to parse dem file: %s", err)
	}

	for _, hook := range hooks {
		time := hook.Timestamp - startTime
		content := ""
		if hook.IsMissed {
			content = "Missed"
		} else if *hook.IsEnemy {
			if *hook.IsCreep {
				content = "Hooked enemy creep"
			} else {
				content = "Hooked enemy"
			}
		} else {
			if *hook.IsCreep {
				content = "Hooked ally creep"
			} else {
				content = "Hooked ally"
			}
		}
		fmt.Printf("[%s] %s\n", lib.FormatMinutesSeconds(time), content)
	}

	log.Printf("Parse Complete!\n")
}
