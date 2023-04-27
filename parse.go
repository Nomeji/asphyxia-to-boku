package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type asphyxiaData map[string]interface{}
type Difficulity uint8
type Grade uint8
type Medal uint8

const (
	Easy Difficulity = iota
	Normal
	Hyper
	EX

	E   Grade = 1
	D   Grade = 2
	C   Grade = 3
	B   Grade = 4
	A   Grade = 5
	AA  Grade = 6
	AAA Grade = 7
	S   Grade = 8

	FailedCircle     Medal = 1
	FailedDiamond    Medal = 2
	FailedStar       Medal = 3
	EasyClear        Medal = 4
	ClearCircle      Medal = 5
	ClearDiamond     Medal = 6
	ClearStar        Medal = 7
	FullComboCircle  Medal = 8
	FullComboDiamond Medal = 9
	FullComboStar    Medal = 10
	Perfect          Medal = 11
)

func (d Difficulity) String() string {
	switch d {
	case Easy:
		return "Easy"
	case Normal:
		return "Normal"
	case Hyper:
		return "Hyper"
	case EX:
		return "EX"
	}
	return "unknown"
}

func (c Grade) String() string {
	switch c {
	case S:
		return "S"
	case AAA:
		return "AAA"
	case AA:
		return "AA"
	case A:
		return "A"
	case B:
		return "B"
	case C:
		return "C"
	case D:
		return "D"
	case E:
		return "E"
	default:
		return "unknown"
	}
}

func (m Medal) String() string {
	switch m {
	case FailedCircle:
		return "failedCircle"
	case FailedDiamond:
		return "failedDiamond"
	case FailedStar:
		return "failedStar"
	case EasyClear:
		return "easyClear"
	case ClearCircle:
		return "clearCircle"
	case ClearDiamond:
		return "clearDiamond"
	case ClearStar:
		return "clearStar"
	case FullComboCircle:
		return "fullComboCircle"
	case FullComboDiamond:
		return "fullComboDiamond"
	case FullComboStar:
		return "fullComboStar"
	case Perfect:
		return "perfect"
	default:
		return "unkown"
	}
}

func (m Medal) Lamp() string {
	switch {
	case m == Perfect:
		return "Perfect"
	case m >= 8:
		return "FULL COMBO"
	case m >= 5:
		return "CLEAR"
	case m == EasyClear:
		return "EASY CLEAR"
	default:
		return "FAILED"
	}
}

func parseScores() []score {
	file, err := os.Open(*asphyxiaDbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var asphyxiaDatas []asphyxiaData
	i := 0
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var asphyxiaData asphyxiaData
		err = json.Unmarshal(scanner.Bytes(), &asphyxiaData)
		if err != nil {
			log.Fatal(err)
		}
		asphyxiaDatas = append(asphyxiaDatas, asphyxiaData)
		i++
	}

	var scores []score

	for i := 0; i < len(asphyxiaDatas); i++ {
		if asphyxiaDatas[i]["collection"] == "scores" {
			currentRefID := asphyxiaDatas[i]["__refid"].(string)

			if !confirmRefID(currentRefID) {
				continue
			}

			date := asphyxiaDatas[i]["updatedAt"].(map[string]interface{})["$$date"].(float64)

			for key, value := range asphyxiaDatas[i]["scores"].(map[string]interface{}) {
				scoreMap := value.(map[string]interface{})
				scoreNb := scoreMap["score"].(float64)

				chartId := strings.Split(key, ":")[0]

				difficulityIndex, err := strconv.Atoi(strings.Split(key, ":")[1])
				if err != nil {
					log.Fatalln(err)
				}
				difficulty := Difficulity(difficulityIndex).String()
				// Difficulty unknown probably means battle mode (needs confirmation)
				if difficulty == "unknown" {
					continue
				}

				clearType := scoreMap["clear_type"]
				// Player exited before end of song: no medal
				if clearType == nil {
					continue
				}
				clearTypeFloat := clearType.(float64)
				// clear_types in unilab are *100?
				if clearTypeFloat >= 100 {
					clearTypeFloat = clearTypeFloat / 100
				}
				medal := Medal(int(clearTypeFloat))
				clearMedal := medal.String()

				var score = score{
					Score:        uint(scoreNb),
					ClearMedal:   clearMedal,
					Difficulity:  difficulty,
					MatchType:    "inGameID",
					Identifier:   chartId,
					TimeAchieved: uint(date),
				}

				scores = append(scores, score)
			}
		}
	}

	return scores
}

func confirmRefID(currentRefID string) bool {
	if *cardID != "" {
		return *cardID == currentRefID
	}

	fmt.Printf("Import scores for %v? [Y/n] ", currentRefID)

	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalln(err)
	}
	text = strings.TrimSpace(text)

	return text == "y" || text == "Y" || text == ""
}
