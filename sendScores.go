package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type config struct {
	Key string
	Url string
}

type score struct {
	Score        uint   `json:"score"`
	ClearMedal   string `json:"clearMedal"`
	MatchType    string `json:"matchType"`
	Identifier   string `json:"identifier"`
	Difficulity  string `json:"difficulty"`
	TimeAchieved uint   `json:"timeAchieved"`
}

type meta struct {
	Game     string `json:"game"`
	Playtype string `json:"playtype"`
	Service  string `json:"service"`
}

type batchManual struct {
	Meta   meta    `json:"meta"`
	Scores []score `json:"scores"`
}

func buildBatchManual(asphyxiaDbPath string) batchManual {
	scores := parseScores(asphyxiaDbPath)

	return batchManual{
		Meta: meta{
			Game:     "popn",
			Playtype: "9B",
			Service:  "Asphyxia exporter",
		},
		Scores: scores,
	}
}

func sendScores(c *config, asphyxiaDbPath string) {
	batchManual := buildBatchManual(asphyxiaDbPath)

	if len(batchManual.Scores) == 0 {
		fmt.Println("Nothing to import exiting")
		return
	}

	jsonb, err := json.Marshal(batchManual)
	if err != nil {
		log.Fatalln(err)
	}

	body := bytes.NewBuffer(jsonb)
	req, err := http.NewRequest(http.MethodPost, c.Url, body)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(respBody))
}

func readConfig(configFile string) *config {
	file, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalln(err)
	}

	c := &config{}
	err = json.Unmarshal(file, c)
	if err != nil {
		log.Fatalln(err)
	}

	return c
}
