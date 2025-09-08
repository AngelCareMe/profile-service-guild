package dto

type MythScoreDto struct {
	Current struct {
		Rating float64 `json:"rating"`
	} `json:"current_mythic_rating"`
}

type CharacterSummary struct {
	Name  string `json:"name"`
	Level int    `json:"level"`
	Realm struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	} `json:"realm"`
	PlayableClass struct {
		Name string `json:"name"`
	} `json:"playable_class"`
	PlayableRace struct {
		Name string `json:"name"`
	} `json:"playable_race"`
	Faction struct {
		Name string `json:"name"`
	} `json:"faction"`
}

type UserDTO struct {
	ID        string `json:"id"`
	Battletag string `json:"battletag"`
}

type BlizzardProfileResponse struct {
	WowAccounts []struct {
		Characters []CharacterSummary `json:"characters"`
	} `json:"wow_accounts"`
}

type CharacterDetailsResponse struct {
	Spec struct {
		Name string `json:"name"`
	} `json:"active_spec"`
	Guild struct {
		Name string `json:"name"`
	} `json:"guild"`
	Ilvl int `json:"average_item_level"`
}
