package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const apiBaseURL = "http://localhost:8081/api"

type Game struct {
	ID     int      `json:"id"`
	Name   string   `json:"name"`
	Year   int      `json:"year"`
	Studio string   `json:"studio"`
	Genres []string `json:"genres"`
}

type Tournament struct {
	ID     int     `json:"id"`
	GameID int     `json:"gameId"`
	Name   string  `json:"name"`
	City   string  `json:"city"`
	Date   string  `json:"date"`
	Lat    float64 `json:"lat"`
	Lng    float64 `json:"lng"`
}

var apiClient = &http.Client{Timeout: 3 * time.Second}

func fetchJSON(url string, target any) error {
	resp, err := apiClient.Get(url)
	if err != nil {
		return fmt.Errorf("impossible de joindre l'API (%s) : %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("l'API a répondu %d %s pour %s", resp.StatusCode, http.StatusText(resp.StatusCode), url)
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("réponse JSON illisible : %w", err)
	}
	return nil
}

func FetchGames() ([]Game, error) {
	var games []Game
	err := fetchJSON(apiBaseURL+"/games", &games)
	return games, err
}

func FetchTournaments() ([]Tournament, error) {
	var tournaments []Tournament
	err := fetchJSON(apiBaseURL+"/tournaments", &tournaments)
	return tournaments, err
}