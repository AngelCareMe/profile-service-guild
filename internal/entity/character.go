package entity

type Character struct {
	UserID      string  `json:"user_id" db:"user_id"`
	Battletag   string  `json:"battletag" db:"battletag"`
	Name        string  `json:"name" db:"name"`
	Race        string  `json:"race" db:"user_id"`
	Faction     string  `json:"faction" db:"user_id"`
	Class       string  `json:"class" db:"user_id"`
	Spec        string  `json:"spec" db:"user_id"`
	Lvl         int     `json:"lvl" db:"user_id"`
	Ilvl        int     `json:"ilvl" db:"ilvl"`
	Guild       string  `json:"guild" db:"guild"`
	MythicScore float64 `json:"mythic_score" db:"mythic_score"`
	IsMain      bool    `json:"is_main" db:"is_main"`
}
