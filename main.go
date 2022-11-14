package main

import (
	"log"

	"github.com/Mines-Little-Theatre/team-trivia-scraper/bot"
	"github.com/Mines-Little-Theatre/team-trivia-scraper/utils"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	bot.Run(utils.ReadConfig())
}
