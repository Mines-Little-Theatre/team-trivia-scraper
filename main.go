package main

import (
	"context"
	"log"

	"github.com/Mines-Little-Theatre/team-trivia-scraper/bot"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	bot.Run(context.Background())
}
