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
	ID   int `json:"id"`
	Spec struct {
		Name string `json:"name"`
	} `json:"active_spec"`
	Guild struct {
		Name  string `json:"name"`
		ID    int    `json:"id"`
		Realm struct {
			Name string `json:"name"`
			Slug string `json:"slug"`
		} `json:"realm"`
		Faction struct {
			Name string `json:"name"`
		} `json:"faction"`
	} `json:"guild"`
	Ilvl int `json:"average_item_level"`
}

type ProfileResponse struct {
	BlizzardID string
	Battletag  string
	Characters []CharacterResponse
}

type CharacterResponse struct {
	Name        string  `json:"name" db:"name"`
	Realm       string  `json:"realm" db:"realm"`
	Race        string  `json:"race" db:"race"`
	Faction     string  `json:"faction" db:"faction"`
	Class       string  `json:"class" db:"class"`
	Spec        string  `json:"spec" db:"spec"`
	Lvl         int     `json:"lvl" db:"lvl"`
	Ilvl        int     `json:"ilvl" db:"ilvl"`
	Guild       string  `json:"guild" db:"guild"`
	MythicScore float64 `json:"mythic_score" db:"mythic_score"`
	IsMain      bool    `json:"is_main" db:"is_main"`
}
