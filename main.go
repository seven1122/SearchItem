package main

import (
	_ "SearchItem/machers"
	"SearchItem/search"
	"log"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)

}
func main() {
	search.Run("范冰冰")
}
