package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	analysisdomain "valorant-tactical-trainer/internal/domain/analysis"
	matchdomain "valorant-tactical-trainer/internal/domain/match"
	"valorant-tactical-trainer/internal/domain/player"

	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

func Open(ctx context.Context) (*Store, error) {
	dataDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("resolve config dir: %w", err)
	}

	appDir := filepath.Join(dataDir, "ValorantTacticalTrainer")
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		return nil, fmt.Errorf("create app data dir: %w", err)
	}

	return OpenPath(ctx, filepath.Join(appDir, "trainer.db"))
}

func OpenPath(ctx context.Context, dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	store := &Store{db: db}
	if err := store.migrate(ctx); err != nil {
		db.Close()
		return nil, err
	}

	return store, nil
}

func (store *Store) Close() error {
	if store == nil || store.db == nil {
		return nil
	}

	return store.db.Close()
}

func (store *Store) SavePlayerWithConsent(ctx context.Context, account player.Account, consent player.Consent) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin save player: %w", err)
	}
	defer tx.Rollback()

	now := time.Now().UTC()
	_, err = tx.ExecContext(ctx, `
		insert into players (puuid, name, tag, region, account_level, card_small, card_large, last_update, created_at, updated_at)
		values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		on conflict(puuid) do update set
			name = excluded.name,
			tag = excluded.tag,
			region = excluded.region,
			account_level = excluded.account_level,
			card_small = excluded.card_small,
			card_large = excluded.card_large,
			last_update = excluded.last_update,
			updated_at = excluded.updated_at
	`, account.PUUID, account.Name, account.Tag, account.Region, account.AccountLevel, account.CardSmall, account.CardLarge, account.LastUpdate, now, now)
	if err != nil {
		return fmt.Errorf("save player: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
		insert into consents (player_puuid, name, tag, region, provider, consent_version, consented_at)
		values (?, ?, ?, ?, ?, ?, ?)
	`, consent.PlayerPUUID, consent.Name, consent.Tag, consent.Region, consent.Provider, consent.ConsentVersion, consent.ConsentedAt)
	if err != nil {
		return fmt.Errorf("save consent: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `insert into settings (key, value, updated_at) values ('current_player_puuid', ?, ?) on conflict(key) do update set value = excluded.value, updated_at = excluded.updated_at`, account.PUUID, now); err != nil {
		return fmt.Errorf("save current player setting: %w", err)
	}

	return tx.Commit()
}

func (store *Store) CurrentPlayer(ctx context.Context) (player.Account, bool, error) {
	var puuid string
	err := store.db.QueryRowContext(ctx, `select value from settings where key = 'current_player_puuid'`).Scan(&puuid)
	if errors.Is(err, sql.ErrNoRows) {
		return player.Account{}, false, nil
	}
	if err != nil {
		return player.Account{}, false, fmt.Errorf("read current player setting: %w", err)
	}

	var account player.Account
	err = store.db.QueryRowContext(ctx, `
		select puuid, name, tag, region, account_level, card_small, card_large, last_update
		from players
		where puuid = ?
	`, puuid).Scan(&account.PUUID, &account.Name, &account.Tag, &account.Region, &account.AccountLevel, &account.CardSmall, &account.CardLarge, &account.LastUpdate)
	if errors.Is(err, sql.ErrNoRows) {
		return player.Account{}, false, nil
	}
	if err != nil {
		return player.Account{}, false, fmt.Errorf("read current player: %w", err)
	}

	return account, true, nil
}

func (store *Store) SaveSetting(ctx context.Context, key string, value string) error {
	_, err := store.db.ExecContext(ctx, `
		insert into settings (key, value, updated_at)
		values (?, ?, ?)
		on conflict(key) do update set value = excluded.value, updated_at = excluded.updated_at
	`, key, value, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("save setting: %w", err)
	}

	return nil
}

func (store *Store) SaveAPICache(ctx context.Context, key string, endpoint string, payload string, ttl time.Duration) error {
	now := time.Now().UTC()
	_, err := store.db.ExecContext(ctx, `
		insert into api_cache (cache_key, endpoint, response_json, expires_at, created_at)
		values (?, ?, ?, ?, ?)
		on conflict(cache_key) do update set
			endpoint = excluded.endpoint,
			response_json = excluded.response_json,
			expires_at = excluded.expires_at,
			created_at = excluded.created_at
	`, key, endpoint, payload, now.Add(ttl), now)
	if err != nil {
		return fmt.Errorf("save api cache: %w", err)
	}

	return nil
}

func (store *Store) APICache(ctx context.Context, key string) (string, bool, error) {
	var payload string
	var expiresAt time.Time
	err := store.db.QueryRowContext(ctx, `select response_json, expires_at from api_cache where cache_key = ?`, key).Scan(&payload, &expiresAt)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false, nil
	}
	if err != nil {
		return "", false, fmt.Errorf("read api cache: %w", err)
	}
	if time.Now().UTC().After(expiresAt) {
		return "", false, nil
	}

	return payload, true, nil
}

func (store *Store) SaveMatches(ctx context.Context, summaries []matchdomain.Summary) (int, error) {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin save matches: %w", err)
	}
	defer tx.Rollback()

	inserted := 0
	now := time.Now().UTC()
	for _, summary := range summaries {
		result, err := tx.ExecContext(ctx, `
			insert into matches (
				match_id, player_puuid, map_name, mode, queue, season_id, region, cluster,
				game_start, game_length, rounds_played, agent, team, kills, deaths, assists,
				headshots, bodyshots, legshots, damage_made, raw_json, created_at, updated_at
			)
			values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			on conflict(match_id, player_puuid) do update set
				map_name = excluded.map_name,
				mode = excluded.mode,
				queue = excluded.queue,
				season_id = excluded.season_id,
				region = excluded.region,
				cluster = excluded.cluster,
				game_start = excluded.game_start,
				game_length = excluded.game_length,
				rounds_played = excluded.rounds_played,
				agent = excluded.agent,
				team = excluded.team,
				kills = excluded.kills,
				deaths = excluded.deaths,
				assists = excluded.assists,
				headshots = excluded.headshots,
				bodyshots = excluded.bodyshots,
				legshots = excluded.legshots,
				damage_made = excluded.damage_made,
				raw_json = excluded.raw_json,
				updated_at = excluded.updated_at
		`, summary.MatchID, summary.PlayerPUUID, summary.MapName, summary.Mode, summary.Queue, summary.SeasonID, summary.Region, summary.Cluster, summary.GameStart, summary.GameLength, summary.RoundsPlayed, summary.Agent, summary.Team, summary.Kills, summary.Deaths, summary.Assists, summary.Headshots, summary.Bodyshots, summary.Legshots, summary.DamageMade, summary.RawJSON, now, now)
		if err != nil {
			return 0, fmt.Errorf("save match: %w", err)
		}

		rows, err := result.RowsAffected()
		if err == nil && rows > 0 {
			inserted++
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit save matches: %w", err)
	}

	return inserted, nil
}

func (store *Store) MatchesForPlayer(ctx context.Context, puuid string) ([]matchdomain.Summary, error) {
	rows, err := store.db.QueryContext(ctx, `
		select match_id, player_puuid, map_name, mode, queue, season_id, region, cluster,
			game_start, game_length, rounds_played, agent, team, kills, deaths, assists,
			headshots, bodyshots, legshots, damage_made, raw_json
		from matches
		where player_puuid = ?
		order by game_start desc
	`, puuid)
	if err != nil {
		return nil, fmt.Errorf("list matches: %w", err)
	}
	defer rows.Close()

	summaries := []matchdomain.Summary{}
	for rows.Next() {
		var summary matchdomain.Summary
		if err := rows.Scan(&summary.MatchID, &summary.PlayerPUUID, &summary.MapName, &summary.Mode, &summary.Queue, &summary.SeasonID, &summary.Region, &summary.Cluster, &summary.GameStart, &summary.GameLength, &summary.RoundsPlayed, &summary.Agent, &summary.Team, &summary.Kills, &summary.Deaths, &summary.Assists, &summary.Headshots, &summary.Bodyshots, &summary.Legshots, &summary.DamageMade, &summary.RawJSON); err != nil {
			return nil, fmt.Errorf("scan match: %w", err)
		}
		summaries = append(summaries, summary)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate matches: %w", err)
	}

	return summaries, nil
}

func (store *Store) SaveReport(ctx context.Context, report analysisdomain.Report) (analysisdomain.Report, error) {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return analysisdomain.Report{}, fmt.Errorf("begin save report: %w", err)
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, `
		insert into analysis_reports (player_puuid, generated_at, match_count, average_kda, headshot_percent, average_damage, top_agent, top_map, summary)
		values (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, report.PlayerPUUID, report.GeneratedAt, report.MatchCount, report.AverageKDA, report.HeadshotPercent, report.AverageDamage, report.TopAgent, report.TopMap, report.Summary)
	if err != nil {
		return analysisdomain.Report{}, fmt.Errorf("save report: %w", err)
	}

	reportID, err := result.LastInsertId()
	if err != nil {
		return analysisdomain.Report{}, fmt.Errorf("read report id: %w", err)
	}
	report.ID = reportID

	for _, finding := range report.Findings {
		if _, err := tx.ExecContext(ctx, `
			insert into analysis_findings (report_id, type, severity, confidence, title, description, evidence)
			values (?, ?, ?, ?, ?, ?, ?)
		`, reportID, finding.Type, finding.Severity, finding.Confidence, finding.Title, finding.Description, joinEvidence(finding.Evidence)); err != nil {
			return analysisdomain.Report{}, fmt.Errorf("save finding: %w", err)
		}
	}

	for _, recommendation := range report.Recommendations {
		if _, err := tx.ExecContext(ctx, `
			insert into training_recommendations (report_id, title, drill, priority, reason, evidence, status, created_at, updated_at)
			values (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, reportID, recommendation.Title, recommendation.Drill, recommendation.Priority, recommendation.Reason, joinEvidence(recommendation.Evidence), recommendation.Status, time.Now().UTC(), time.Now().UTC()); err != nil {
			return analysisdomain.Report{}, fmt.Errorf("save recommendation: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return analysisdomain.Report{}, fmt.Errorf("commit report: %w", err)
	}

	return report, nil
}

func (store *Store) ResetAll(ctx context.Context) error {
	tables := []string{"training_recommendations", "analysis_findings", "analysis_reports", "matches", "api_cache", "consents", "players", "settings"}
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin reset: %w", err)
	}
	defer tx.Rollback()

	for _, table := range tables {
		if _, err := tx.ExecContext(ctx, fmt.Sprintf("delete from %s", table)); err != nil {
			return fmt.Errorf("reset %s: %w", table, err)
		}
	}

	return tx.Commit()
}

func (store *Store) Setting(ctx context.Context, key string) (string, bool, error) {
	var value string
	err := store.db.QueryRowContext(ctx, `select value from settings where key = ?`, key).Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false, nil
	}
	if err != nil {
		return "", false, fmt.Errorf("read setting: %w", err)
	}

	return value, true, nil
}

func (store *Store) migrate(ctx context.Context) error {
	statements := []string{
		`create table if not exists players (
			puuid text primary key,
			name text not null,
			tag text not null,
			region text not null,
			account_level integer not null default 0,
			card_small text not null default '',
			card_large text not null default '',
			last_update text not null default '',
			created_at datetime not null,
			updated_at datetime not null
		)`,
		`create table if not exists consents (
			id integer primary key autoincrement,
			player_puuid text not null,
			name text not null,
			tag text not null,
			region text not null,
			provider text not null,
			consent_version text not null,
			consented_at datetime not null,
			revoked_at datetime,
			foreign key(player_puuid) references players(puuid)
		)`,
		`create table if not exists settings (
			key text primary key,
			value text not null,
			updated_at datetime not null
		)`,
		`create table if not exists api_cache (
			cache_key text primary key,
			endpoint text not null,
			response_json text not null,
			expires_at datetime not null,
			created_at datetime not null
		)`,
		`create table if not exists matches (
			match_id text not null,
			player_puuid text not null,
			map_name text not null default '',
			mode text not null default '',
			queue text not null default '',
			season_id text not null default '',
			region text not null default '',
			cluster text not null default '',
			game_start integer not null default 0,
			game_length integer not null default 0,
			rounds_played integer not null default 0,
			agent text not null default '',
			team text not null default '',
			kills integer not null default 0,
			deaths integer not null default 0,
			assists integer not null default 0,
			headshots integer not null default 0,
			bodyshots integer not null default 0,
			legshots integer not null default 0,
			damage_made integer not null default 0,
			raw_json text not null default '',
			created_at datetime not null,
			updated_at datetime not null,
			primary key(match_id, player_puuid),
			foreign key(player_puuid) references players(puuid)
		)`,
		`create table if not exists analysis_reports (
			id integer primary key autoincrement,
			player_puuid text not null,
			generated_at datetime not null,
			match_count integer not null,
			average_kda real not null,
			headshot_percent real not null,
			average_damage real not null,
			top_agent text not null default '',
			top_map text not null default '',
			summary text not null,
			foreign key(player_puuid) references players(puuid)
		)`,
		`create table if not exists analysis_findings (
			id integer primary key autoincrement,
			report_id integer not null,
			type text not null,
			severity text not null,
			confidence real not null,
			title text not null,
			description text not null,
			evidence text not null,
			foreign key(report_id) references analysis_reports(id)
		)`,
		`create table if not exists training_recommendations (
			id integer primary key autoincrement,
			report_id integer not null,
			title text not null,
			drill text not null,
			priority text not null,
			reason text not null,
			evidence text not null,
			status text not null,
			created_at datetime not null,
			updated_at datetime not null,
			foreign key(report_id) references analysis_reports(id)
		)`,
	}

	for _, statement := range statements {
		if _, err := store.db.ExecContext(ctx, statement); err != nil {
			return fmt.Errorf("migrate sqlite: %w", err)
		}
	}

	return nil
}

func joinEvidence(values []string) string {
	result := ""
	for index, value := range values {
		if index > 0 {
			result += " | "
		}
		result += value
	}
	return result
}
