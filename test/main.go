package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"profile-service/internal/entity"
	"profile-service/pkg/dto"
	"time"
)

type CharacterResponseDTO struct {
	Name  string `json:"name"`
	Realm struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	} `json:"realm"`
	Class struct {
		Name string `json:"name"`
	} `json:"playable_class"`
	Race struct {
		Name string `json:"name"`
	} `json:"playable_race"`
	Faction struct {
		Type string `json:"type"`
	} `json:"faction"`
	Lvl int `json:"level"`
}

func main() {
	var token string
	ctx := context.Background()
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			IdleConnTimeout: 30 * time.Second,
		},
	}
	fmt.Println("Entry blizz token")
	fmt.Scanln(&token)

	req, err := http.NewRequestWithContext(ctx, "GET", "https://eu.api.blizzard.com/profile/user/wow?namespace=profile-eu&locale=ru_RU", nil)
	if err != nil {
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	fmt.Printf("Response body: %s\n", string(body))

	var profile dto.BlizzardProfileResponse
	err = json.Unmarshal(body, &profile) // Исправлено: добавлена проверка ошибки
	if err != nil {
		fmt.Printf("Error unmarshaling JSON: %v\n", err)
		return
	}
	var characters []entity.Character
	for _, account := range profile.WowAccounts {
		for _, char := range account.Characters {
			character := entity.Character{
				Name:    char.Name,
				Race:    char.PlayableRace.Name,
				Faction: char.Faction.Name,
				Class:   char.PlayableClass.Name,
				Lvl:     char.Level,
				IsMain:  false,
			}
			characters = append(characters, character)
		}
	}

	fmt.Printf("Found %d characters:\n", len(characters))
	for _, char := range characters {
		fmt.Printf("Name: %s, Race: %s, Class: %s, Faction: %s, Level: %d\n",
			char.Name, char.Race, char.Class, char.Faction, char.Lvl)
	}
}
