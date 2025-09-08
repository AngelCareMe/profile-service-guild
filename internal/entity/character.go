package entity

type Character struct {
	BlizzardID  string  `json:"blizzard_id" db:"blizzard_id"`
	Battletag   string  `json:"battletag" db:"battletag"`
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
