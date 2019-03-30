package main

import (
	"fmt"
	"log"
	"os"

	poker "github.com/bjrnt/tdd-app"
)

const dbFileName = "game.db.json"

func main() {
	store, err := poker.FileSystemPlayerStoreFromFile(dbFileName)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Let's play poker")
	fmt.Println("Type {Name} wins to record a win")
	game := poker.NewTexasHoldem(poker.BlindAlerterFunc(poker.StdOutAlerter), store)
	poker.NewCLI(os.Stdin, os.Stdout, game).PlayPoker()
}
