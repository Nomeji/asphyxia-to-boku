package main

import "flag"

func main() {
	configPath := flag.String("config", "silent-config.json", "path to config file containing { \"key\": \"apikey\", \"url\": \"bokutachiurl/ir/direct-manual/import\"")
	asphyxiaDbPath := flag.String("asphyxia", "popn:db", "path to asphyxia database file")

	flag.Parse()

	config := readConfig(*configPath)
	sendScores(config, *asphyxiaDbPath)
}
