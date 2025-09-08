package database

import (
	"context"
	"profile-service/internal/entity"
	"profile-service/pkg/errors"

	sq "github.com/Masterminds/squirrel"
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
	query := psql.
		Insert("characters").
		Columns(
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

	sql, args, err := query.Suffix("ON CONFLICT (user_id, name, realm) DO UPDATE SET " +
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

	pr.log.Infof("Saved/updated %d characters", len(characters))
	return nil
}

func (pr *postgresRepository) GetCharacters(ctx context.Context, blizzardID string) ([]entity.Character, error) {
	query := psql.Select(
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
		From("characters").
		Where(sq.Eq{"blizzard_id": blizzardID})

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
	var char entity.Character

	for rows.Next() {
		err := rows.Scan(
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
			pr.log.WithError(err).Error("failed to scan character row")
			return nil, errors.NewAppError("failed to scan character row", err)
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

