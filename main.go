// nba-schedule returns all the NBA games played on an input date //TODO proper
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/games", GamesHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

type Games struct {
	Games []Game `json:"games"`
}

type Game struct {
	Id           string    `json:"gameId"`
	StartTime    time.Time `json:"startTimeUTC"`
	VisitingTeam Team      `json:"vTeam"`
	HomeTeam     Team      `json:"hTeam"`
	Period       Period    `json:"period"`
}

type Team struct {
	Id      string `json:"teamId"`
	TriCode string `json:"triCode"`
}

type Period struct {
	Current int `json:"current"`
}

func GamesHandler(w http.ResponseWriter, r *http.Request) {
	games, err := GetGames()
	if err != nil {
		log.Printf("error retrieving games: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, `{"error":"server error occurred"}`)
		return
	}

	enc := json.NewEncoder(w)
	err = enc.Encode(&games)
	if err != nil {
		log.Printf("error encoding games: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, `{"error":"server error occurred"}`)
		return
	}
}

func GetGames() (Games, error) {
	scoreboardURL, err := TodayScoreboardURL()
	if err != nil {
		log.Printf("error retrieving scoreboard url: %s\n", err)
		return Games{}, err
	}

	resp, err := http.Get(scoreboardURL)
	if err != nil {
		log.Printf("error retrieving scoreboard: %s\n", err)
		return Games{}, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	var games Games
	for dec.More() {
		err := dec.Decode(&games)
		if err != nil {
			log.Printf("error decoding response: %s\n", err)
			return Games{}, err
		}
	}

	return games, nil
}

type NBATodayResponse struct {
	Links map[string]string `json:"links"`
}

func TodayScoreboardURL() (string, error) {
	const NBABaseURL = "http://data.nba.net"
	const NBATodayRoute = "/10s/prod/v1/today.json"

	NBATodayURL := fmt.Sprintf("%s%s", NBABaseURL, NBATodayRoute)

	resp, err := http.Get(NBATodayURL)
	if err != nil {
		log.Printf("error retrieving today url: %s\n", err)
		return "", err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	var todayResp NBATodayResponse

	for dec.More() {
		err := dec.Decode(&todayResp)
		if err != nil {
			log.Printf("error decoding response: %s\n", err)
			return "", err
		}
	}

	NBATodayScoreboardURL := fmt.Sprintf("%s%s", NBABaseURL, todayResp.Links["todayScoreboard"])
	return NBATodayScoreboardURL, nil
}
