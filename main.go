// nba-schedule returns all the NBA games played on an input date
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	// fmt.Printf(rawGames())
	// fmt.Printf(games())
	games()
}

type Games struct {
	Games []Game `json:"games"`
}

type Game struct {
	Id           string    `json:"gameId"`
	StartTime    time.Time `json:"startTimeUTC"`
	VisitingTeam Team      `json:"vTeam"`
	HomeTeam     Team      `json:"hTeam"`
}

type Team struct {
	Id      string `json:"teamId"`
	TriCode string `json:"triCode"`
}

func games() Games {
	resp, err := http.Get(scoreboardURL())
	if err != nil {
		log.Fatalf("error retrieving scoreboard: %s\n", err)
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	var games Games
	for dec.More() {
		err := dec.Decode(&games)
		if err != nil {
			log.Fatalf("error decoding response: %s\n", err)
		}
	}

	fmt.Printf("%+v\n", games)

	return games
}

func rawGames() string {
	resp, err := http.Get(scoreboardURL())
	if err != nil {
		log.Fatalf("error retrieving scoreboard: %s\n", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error reading scoreboard response: %s\n", err)
	}

	return string(body)
}

func scoreboardURL() string {
	return "http://data.nba.net/data/10s/prod/v1/20161112/scoreboard.json"
}
