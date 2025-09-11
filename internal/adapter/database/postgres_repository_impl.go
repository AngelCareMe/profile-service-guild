package database

import (
	"context"
	"profile-service/internal/entity"
	"profile-service/pkg/errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type postgresRepository struct {
	pool *pgxpool.Pool
	log  *logrus.Logger
}

func NewPostgresRepository(pool *pgxpool.Pool, log *logrus.Logger) *postgresRepository {
	return &postgresRepository{pool: pool, log: log}
}

func (pr *postgresRepository) SaveCharacters(ctx context.Context, characters []entity.Character) error {
	tx, err := pr.pool.Begin(ctx)
	if err != nil {
		pr.log.WithError(err).Error("failed to begin transaction")
		return errors.NewAppError("failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	query := psql.
		Insert("profile").
		Columns(
			"character_id",
			"blizzard_id",
			"battletag",
			"name",
			"realm",
			"race",
			"faction",
			"class",
			"spec",
			"lvl",
			"ilvl",
			"guild",
			"mythic_score",
			"is_main",
		)

	for _, char := range characters {
		query = query.Values(
			char.CharacterID,
			char.BlizzardID,
			char.Battletag,
			char.Name,
			char.Realm,
			char.Race,
			char.Faction,
			char.Class,
			char.Spec,
			char.Lvl,
			char.Ilvl,
			char.Guild,
			char.MythicScore,
			char.IsMain,
		)
	}

	sql, args, err := query.Suffix("ON CONFLICT (blizzard_id, name, realm) DO UPDATE SET " +
		"character_id = EXCLUDED.character_id, " +
		"battletag = EXCLUDED.battletag, " +
		"race = EXCLUDED.race, " +
		"faction = EXCLUDED.faction, " +
		"class = EXCLUDED.class, " +
		"spec = EXCLUDED.spec, " +
		"lvl = EXCLUDED.lvl, " +
		"ilvl = EXCLUDED.ilvl, " +
		"guild = EXCLUDED.guild, " +
		"mythic_score = EXCLUDED.mythic_score, " +
		"is_main = EXCLUDED.is_main").
		ToSql()

	if err != nil {
		pr.log.WithError(err).Error("failed build query for save characters")
		return errors.NewAppError("failed build query for save characters", err)
	}

	_, err = pr.pool.Exec(ctx, sql, args...)
	if err != nil {
		pr.log.WithError(err).Error("failed execute SQL save characters")
		return errors.NewAppError("failed execute SQL save characters", err)
	}

	if err := tx.Commit(ctx); err != nil {
		pr.log.WithError(err).Error("failed to commit transaction")
		return errors.NewAppError("failed to commit transaction", err)
	}

	pr.log.Infof("Saved/updated %d characters", len(characters))
	return nil
}

func (pr *postgresRepository) SaveGuilds(ctx context.Context, guilds []entity.Guild) error {
	if len(guilds) == 0 {
		return nil
	}

	query := psql.
		Insert("guild").
		Columns(
			"character_id",
			"guild_id",
			"name",
			"name_slug",
			"realm",
			"realm_slug",
			"faction",
		)

	for _, g := range guilds {
		query = query.Values(
			g.CharacterID,
			g.GuildID,
			g.Name,
			g.NameSlug,
			g.Realm,
			g.RealmSlug,
			g.Faction,
		)
	}

	sql, args, err := query.Suffix(`
		ON CONFLICT (character_id) DO UPDATE SET
			guild_id = EXCLUDED.guild_id,
			name = EXCLUDED.name,
			name_slug = EXCLUDED.name_slug,
			realm = EXCLUDED.realm,
			realm_slug = EXCLUDED.realm_slug,
			faction = EXCLUDED.faction
	`).ToSql()

	if err != nil {
		pr.log.WithError(err).Error("failed build query for save guilds")
		return errors.NewAppError("failed build query for save guilds", err)
	}

	_, err = pr.pool.Exec(ctx, sql, args...)
	if err != nil {
		pr.log.WithError(err).Error("failed execute SQL save guilds")
		return errors.NewAppError("failed execute SQL save guilds", err)
	}

	pr.log.Infof("Saved/updated %d guilds", len(guilds))
	return nil
}

func (pr *postgresRepository) GetGuildByName(ctx context.Context, nameSlug string) (*entity.Guild, error) {
	query, args, err := psql.Select(
		"character_id",
		"guild_id",
		"name",
		"name_slug",
		"realm",
		"realm_slug",
		"faction",
	).
		From("guild").
		Where(sq.Eq{"name_slug": nameSlug}).
		ToSql()
	if err != nil {
		pr.log.WithError(err).Error("failed build query for get guild")
		return nil, errors.NewAppError("failed build query for get guild", err)
	}

	var g entity.Guild
	err = pr.pool.QueryRow(ctx, query, args...).
		Scan(
			&g.CharacterID,
			&g.GuildID,
			&g.Name,
			&g.NameSlug,
			&g.Realm,
			&g.RealmSlug,
			&g.Faction,
		)

	if err != nil {
		if err == pgx.ErrNoRows {
			pr.log.WithField("name_slug", nameSlug).Info("guild not found")
			return nil, errors.NewAppError("guild not found", err)
		}
		pr.log.WithError(err).Error("failed to scan guild row")
		return nil, errors.NewAppError("failed to scan guild row", err)
	}

	pr.log.Infof("Guild %s get succeeded", nameSlug)
	return &g, nil
}

func (pr *postgresRepository) GetMainCharacterByBlizzardID(ctx context.Context, blizzardID string) (*entity.Character, error) {
	query, args, err := psql.Select(
		"character_id",
		"blizzard_id",
		"battletag",
		"name",
		"realm",
		"race",
		"faction",
		"class",
		"spec",
		"lvl",
		"ilvl",
		"guild",
		"mythic_score",
		"is_main",
	).
		From("profile").
		Where(sq.Eq{
			"blizzard_id": blizzardID,
			"is_main":     true,
		}).
		ToSql()
	if err != nil {
		pr.log.WithError(err).Error("failed build query for get character")
		return nil, errors.NewAppError("failed build query for get character", err)
	}

	var char entity.Character
	err = pr.pool.QueryRow(ctx, query, args...).
		Scan(
			&char.CharacterID,
			&char.BlizzardID,
			&char.Battletag,
			&char.Name,
			&char.Realm,
			&char.Race,
			&char.Faction,
			&char.Class,
			&char.Spec,
			&char.Lvl,
			&char.Ilvl,
			&char.Guild,
			&char.MythicScore,
			&char.IsMain,
		)
	if err != nil {
		if err == pgx.ErrNoRows {
			pr.log.WithField("name", char.Name).Info("character not found")
			return nil, errors.NewAppError("character not found", err)
		}
		pr.log.WithError(err).Error("failed to scan character row")
		return nil, errors.NewAppError("failed to scan character row", err)
	}

	pr.log.Infof("Character %s get succeeded", char.Name)
	return &char, nil
}

func (pr *postgresRepository) GetCharacters(ctx context.Context, blizzardID string) ([]entity.Character, error) {
	query := psql.Select(
		"character_id",
		"blizzard_id",
		"battletag",
		"name",
		"realm",
		"race",
		"faction",
		"class",
		"spec",
		"lvl",
		"ilvl",
		"guild",
		"mythic_score",
		"is_main",
	).
		From("profile").
		Where(sq.Eq{"blizzard_id": blizzardID}).
		OrderBy("mythic_score DESC")

	sql, args, err := query.ToSql()
	if err != nil {
		pr.log.WithError(err).Error("failed build query for get characters")
		return nil, errors.NewAppError("failed build query for get characters", err)
	}

	rows, err := pr.pool.Query(ctx, sql, args...)
	if err != nil {
		pr.log.WithError(err).Error("failed execute SQL get characters")
		return nil, errors.NewAppError("failed execute SQL get characters", err)
	}
	defer rows.Close()

	var characters []entity.Character

	for rows.Next() {
		char := entity.Character{}
		err := rows.Scan(
			&char.CharacterID,
			&char.BlizzardID,
			&char.Battletag,
			&char.Name,
			&char.Realm,
			&char.Race,
			&char.Faction,
			&char.Class,
			&char.Spec,
			&char.Lvl,
			&char.Ilvl,
			&char.Guild,
			&char.MythicScore,
			&char.IsMain,
		)
		if err != nil {
			pr.log.WithError(err).Error("failed to scan character rows")
			return nil, errors.NewAppError("failed to scan character rows", err)
		}
		characters = append(characters, char)
	}

	if err = rows.Err(); err != nil {
		pr.log.WithError(err).Error("rows error")
		return nil, errors.NewAppError("rows error", err)
	}

	pr.log.Infof("Got %d characters for %s", len(characters), blizzardID)
	return characters, nil
}

func (pr *postgresRepository) SetMainCharacter(ctx context.Context, blizzardID, charName string) error {
	tx, err := pr.pool.Begin(ctx)
	if err != nil {
		pr.log.WithError(err).Error("failed to begin transaction")
		return errors.NewAppError("failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	resetQuery, resetArgs, err := psql.
		Update("profile").
		Set("is_main", false).
		Where(sq.Eq{"blizzard_id": blizzardID}).
		ToSql()
	if err != nil {
		pr.log.WithError(err).Error("failed build query for reset main character")
		return errors.NewAppError("failed build query for reset main character", err)
	}

	_, err = tx.Exec(ctx, resetQuery, resetArgs...)
	if err != nil {
		pr.log.WithError(err).Errorf("failed to reset main charcter for user: %s", blizzardID)
		return errors.NewAppError("failed to reset main charcter", err)
	}

	setQuery, setArgs, err := psql.
		Update("profile").
		Set("is_main", true).
		Where(sq.ILike{"name": charName}).
		ToSql()
	if err != nil {
		pr.log.WithError(err).Error("failed build query for set main character")
		return errors.NewAppError("failed build query for set main character", err)
	}

	_, err = tx.Exec(ctx, setQuery, setArgs...)
	if err != nil {
		pr.log.WithError(err).Errorf("failed to set main charcter for user: %s", blizzardID)
		return errors.NewAppError("failed to set main charcter", err)
	}

	if err := tx.Commit(ctx); err != nil {
		pr.log.WithError(err).Error("failed to commit set main transaction")
		return errors.NewAppError("failed to commit set main transaction", err)
	}

	pr.log.WithFields(logrus.Fields{
		"blizzard_id": blizzardID,
		"character":   charName,
	}).Infof("Set main for %s succeeded", charName)

	return nil
}

func (pr *postgresRepository) GetCharacterByName(ctx context.Context, charName string) (*entity.Character, error) {
	query, args, err := psql.Select(
		"character_id",
		"blizzard_id",
		"battletag",
		"name",
		"realm",
		"race",
		"faction",
		"class",
		"spec",
		"lvl",
		"ilvl",
		"guild",
		"mythic_score",
		"is_main",
	).From("profile").
		Where(sq.ILike{"name": charName}).
		ToSql()

	if err != nil {
		pr.log.WithError(err).Error("failed build query for get character")
		return nil, errors.NewAppError("failed build query for get character", err)
	}

	var char entity.Character
	if err := pr.pool.QueryRow(ctx, query, args...).
		Scan(
			&char.CharacterID,
			&char.BlizzardID,
			&char.Battletag,
			&char.Name,
			&char.Realm,
			&char.Race,
			&char.Faction,
			&char.Class,
			&char.Spec,
			&char.Lvl,
			&char.Ilvl,
			&char.Guild,
			&char.MythicScore,
			&char.IsMain,
		); err != nil {
		pr.log.WithError(err).Error("failed to scan character row")
		return nil, errors.NewAppError("failed to scan character row", err)
	}

	pr.log.Infof("%s get succeeded", charName)
	return &char, nil
}
