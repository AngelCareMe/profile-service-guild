package entity

type Guild struct {
	CharacterID int    `json:"character_id" db:"character_id"`
	GuildID     int    `json:"guild_id" db:"guild_id"`
	Name        string `json:"name" db:"name"`
	NameSlug    string `json:"name_slug" db:"name_slug"`
	Realm       string `json:"realm" db:"realm"`
	RealmSlug   string `json:"realm_slug" db:"realm_slug"`
	Faction     string `json:"faction" db:"faction"`
}
