package main

import "flag"

var asphyxiaDbPath *string
var cardID *string

func main() {
	configPath := flag.String("config", "silent-config.json", "path to config file containing { \"key\": \"apikey\", \"url\": \"bokutachiurl/ir/direct-manual/import\"")
	asphyxiaDbPath = flag.String("asphyxia", "popn:db", "path to asphyxia database file")
	cardID = flag.String("cardid", "", "import a specific card ID without prompt")

	flag.Parse()

	config := readConfig(*configPath)
	sendScores(config)
}
