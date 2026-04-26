package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Wild-sergunys/shrtic/internal/model"
)

type LinkRepository struct {
	db *pgxpool.Pool
}

func NewLinkRepository(db *pgxpool.Pool) *LinkRepository {
	return &LinkRepository{db: db}
}

func (r *LinkRepository) Create(ctx context.Context, link *model.Link) error {
	query := `INSERT INTO links (user_id, short_code, long_url) VALUES ($1, $2, $3) RETURNING id, clicks, created_at`
	var userID any
	if link.UserID != nil {
		userID = *link.UserID
	}
	err := r.db.QueryRow(ctx, query, userID, link.ShortURL, link.LongURL).Scan(&link.ID, &link.Clicks, &link.CreatedAt)
	if err != nil {
		return fmt.Errorf("не удалось создать ссылку: %w", err)
	}
	return nil
}

func (r *LinkRepository) FindByShortCode(ctx context.Context, code string) (*model.Link, error) {
	query := `SELECT id, user_id, short_code, long_url, clicks, created_at FROM links WHERE short_code = $1`
	link := &model.Link{}
	var userID *int64
	err := r.db.QueryRow(ctx, query, code).Scan(&link.ID, &userID, &link.ShortURL, &link.LongURL, &link.Clicks, &link.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("не удалось найти ссылку: %w", err)
	}
	link.UserID = userID
	return link, nil
}

func (r *LinkRepository) FindByID(ctx context.Context, id int64) (*model.Link, error) {
	query := `SELECT id, user_id, short_code, long_url, clicks, created_at FROM links WHERE id = $1`
	link := &model.Link{}
	var userID *int64
	err := r.db.QueryRow(ctx, query, id).Scan(&link.ID, &userID, &link.ShortURL, &link.LongURL, &link.Clicks, &link.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("не удалось найти ссылку: %w", err)
	}
	link.UserID = userID
	return link, nil
}

func (r *LinkRepository) FindByUserID(ctx context.Context, userID int64, search string) ([]model.Link, error) {
	var rows pgx.Rows
	var err error

	if search != "" {
		query := `SELECT id, user_id, short_code, long_url, clicks, created_at FROM links WHERE user_id = $1 AND long_url ILIKE $2 ORDER BY created_at DESC`
		rows, err = r.db.Query(ctx, query, userID, "%"+search+"%")
	} else {
		query := `SELECT id, user_id, short_code, long_url, clicks, created_at FROM links WHERE user_id = $1 ORDER BY created_at DESC`
		rows, err = r.db.Query(ctx, query, userID)
	}

	if err != nil {
		return nil, fmt.Errorf("не удалось получить ссылки: %w", err)
	}
	defer rows.Close()

	var links []model.Link
	for rows.Next() {
		var link model.Link
		var uid *int64
		if err := rows.Scan(&link.ID, &uid, &link.ShortURL, &link.LongURL, &link.Clicks, &link.CreatedAt); err != nil {
			return nil, fmt.Errorf("ошибка сканирования ссылки: %w", err)
		}
		link.UserID = uid
		links = append(links, link)
	}

	return links, nil
}

func (r *LinkRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM links WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("не удалось удалить ссылку: %w", err)
	}
	return nil
}

func (r *LinkRepository) IncrementClicks(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `UPDATE links SET clicks = clicks + 1 WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("не удалось увеличить счётчик: %w", err)
	}
	return nil
}

func (r *LinkRepository) SaveClickStat(ctx context.Context, linkID int64, browser, device, country, referer string) (int64, error) {
	var statID int64
	query := `INSERT INTO stats (link_id, browser, device_type, country, referer) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := r.db.QueryRow(ctx, query, linkID, browser, device, country, referer).Scan(&statID)
	if err != nil {
		return 0, fmt.Errorf("не удалось сохранить статистику: %w", err)
	}
	return statID, nil
}

func (r *LinkRepository) GetStats(ctx context.Context, linkID int64) (*model.LinkStats, error) {
	stats := &model.LinkStats{}

	err := r.db.QueryRow(ctx, `SELECT clicks FROM links WHERE id = $1`, linkID).Scan(&stats.TotalClicks)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("ссылка не найдена")
		}
		return nil, fmt.Errorf("не удалось получить счётчик: %w", err)
	}

	browsers, _ := r.getGroupedStats(ctx, linkID, "browser")
	stats.Browsers = browsers

	devices, _ := r.getGroupedStats(ctx, linkID, "device_type")
	stats.Devices = devices

	countries, _ := r.getGroupedStats(ctx, linkID, "country")
	stats.Countries = countries

	referrers, _ := r.getGroupedStats(ctx, linkID, "referer")
	stats.Referrers = referrers

	return stats, nil
}

func (r *LinkRepository) getGroupedStats(ctx context.Context, linkID int64, field string) ([]model.StatItem, error) {
	query := fmt.Sprintf(`
		SELECT %s, COUNT(*) as count
		FROM stats
		WHERE link_id = $1 AND %s IS NOT NULL
		GROUP BY %s
		ORDER BY count DESC
	`, field, field, field)

	rows, err := r.db.Query(ctx, query, linkID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.StatItem
	var total int
	for rows.Next() {
		var item model.StatItem
		if err := rows.Scan(&item.Name, &item.Count); err != nil {
			continue
		}
		items = append(items, item)
		total += item.Count
	}

	for i := range items {
		if total > 0 {
			items[i].Percentage = float64(items[i].Count) / float64(total) * 100
		}
	}

	return items, nil
}
