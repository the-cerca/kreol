package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type Play struct {
	ID          string
	Word        string
	Translation string
	Random      string
}
type Theme struct {
	ID   string
	Name string
}
type Word struct {
	ID   string
	Word string
}
type Language struct {
	ID     string
	Name   string
	Status string
}
type LanguageManager struct {
	DB *sql.DB
}

func (lm *LanguageManager) QueriesLanguages(ctx context.Context) (*[]Language, error) {
	var l []Language
	rows, err := lm.DB.QueryContext(ctx, `select id, name, status from languages`)
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("failed to close rows : %s", err)
		}
	}(rows)
	if err != nil {
		if errors.Is(rows.Err(), sql.ErrNoRows) {
			return nil, errors.New("no rows in languages")
		}
	}

	for rows.Next() {
		var lan Language
		if err := rows.Scan(&lan.ID, &lan.Name, &lan.Status); err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		l = append(l, lan)
	}
	if err := rows.Err(); err != nil {
		return nil, err //
	}

	return &l, nil
}

func (lm *LanguageManager) SubscribeLanguage(ctx context.Context, languageID string) error {
	v, err := GetUserByContext(ctx)
	if err != nil {
		return err
	}
	_, err = lm.DB.ExecContext(ctx, `insert into user_languages (user_id, language_id) values ($1, $2)`, v.ID, languageID)
	if err != nil {
		return err
	}
	return nil
}
func (lm *LanguageManager) QueryNameLanguageByID(ctx context.Context, ID string) (string, error) {
	var n string
	if err := lm.DB.QueryRowContext(ctx, `select name from languages where id=$1`, ID).Scan(&n); err != nil {
		return "", fmt.Errorf("Language not find %w", err)
	}
	return n, nil
}

func (lm *LanguageManager) QueryIdLanguageByName(ctx context.Context, name string) (string, error) {
	var id string
	if err := lm.DB.QueryRowContext(ctx, `select id from languages where name=$1`, name).Scan(&id); err != nil {
		return "", fmt.Errorf("Language not find %w", err)
	}
	return id, nil
}
func (lc *LanguageManager) SubscribedLanguage(ctx context.Context) (*Language, error) {
	v, err := GetUserByContext(ctx)
	if err != nil {
		return nil, err
	}
	query := `SELECT l.id, l.name FROM user_languages ul JOIN languages l ON ul.language_id = l.id WHERE ul.user_id = $1`
	lc.DB.QueryRowContext(ctx, query, v.ID)
	return nil, nil
}
func (lc *LanguageManager) QueriesThemes(ctx context.Context) (*[]Theme, error) {
	var themes []Theme
	rows, err := lc.DB.QueryContext(ctx, `select id, name from themes order by themes.created_at`)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	defer func(r *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("failed to close rows ")
		}
	}(rows)
	for rows.Next() {
		var th Theme
		if err := rows.Scan(&th.ID, &th.Name); err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		themes = append(themes, th)
	}
	return &themes, nil
}
func (lm *LanguageManager) QueryTheme(ctx context.Context, name string) (string, error) {
	var id string
	if err := lm.DB.QueryRowContext(ctx, `select id from themes where name = $1`, name).Scan(&id); err != nil {
		return "", fmt.Errorf("%w", err)
	}
	return id, nil

}
func (lm *LanguageManager) QueryAllUnseenTheme(ctx context.Context, langId string) (*[]Theme, error) {
	u, err := GetUserByContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("user not found %w", err)
	}
	query := `
	select
	  t.id,
	  t.name
	from
	  themes t
	where
	  t.published = true
	  and t.language_id = $1
	  and not exists (
		select
		  1
		from
		  words w
		  join words_views wv on w.id = wv.word_id
		where
		  w.theme_id = t.id
		  and wv.user_id = $2
	  )
	order by
	  t.organize
	`
	var themes []Theme
	rows, err := lm.DB.QueryContext(ctx, query, langId, u.ID)
	if err != nil {
		return nil, fmt.Errorf("query unseen themes : %w", err)
	}
	defer func(r *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("failed to close rows ")
		}
	}(rows)

	for rows.Next() {
		var th Theme
		if err := rows.Scan(&th.ID, &th.Name); err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		themes = append(themes, th)
	}
	return &themes, nil
}
func (lm *LanguageManager) QueryWordByTheme(ctx context.Context, themeId string) (*Play, error) {
	u, err := GetUserByContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("user not found %w", err)
	}
	query := `WITH UnseenWord AS (
		SELECT
			w.id,
			w.word,
			w.language_id,
			w.theme_id
		FROM
			words w
			LEFT JOIN words_views wv ON w.id = wv.word_id AND wv.user_id = $1
			WHERE
			wv.id IS NULL
			AND w.theme_id = $2 
		ORDER BY
			w.created_at
		LIMIT 1
	),
	WordTranslation AS (
		SELECT
			t.translation
		FROM
			translations t
			JOIN UnseenWord uw ON t.word_id = uw.id
	),
	RandomWord AS (
		SELECT
			t.id,
			t.translation
		FROM
			words w
			JOIN translations t ON w.id = t.word_id
		WHERE
			w.theme_id = (SELECT theme_id FROM UnseenWord)
			AND w.language_id = (SELECT language_id FROM UnseenWord)
			AND t.translation != (SELECT translation FROM WordTranslation)
		ORDER BY
			RANDOM()
		LIMIT 1
	)
	SELECT
		uw.id AS word_id,
		uw.word,
		wt.translation AS translation,
		rw.translation AS random_translation
	FROM
		UnseenWord uw
		CROSS JOIN WordTranslation wt
		CROSS JOIN RandomWord rw;
	`
	var p Play
	if err := lm.DB.QueryRowContext(ctx, query, u.ID, themeId).Scan(&p.ID, &p.Word, &p.Translation, &p.Random); err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return &p, nil
}
