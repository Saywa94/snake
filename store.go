package main

import "database/sql"

type Score struct {
	ID    int
	Name  string
	Score int
}

type Store struct {
	conn *sql.DB
}

func (s *Store) Init() error {
	var err error
	s.conn, err = sql.Open("sqlite3", "./snake.db")
	if err != nil {
		return err
	}

	createTableStmt := `
	CREATE TABLE IF NOT EXISTS scores (
		id integer not null primary key,
		name text not null,
		score integer not null
	);`

	if _, err = s.conn.Exec(createTableStmt); err != nil {
		return err
	}

	return nil
}

func (s *Store) GetScores() ([]Score, error) {
	rows, err := s.conn.Query("SELECT * FROM scores ORDER BY score DESC LIMIT 10")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var scores []Score
	for rows.Next() {
		var score Score
		if err := rows.Scan(&score.ID, &score.Name, &score.Score); err != nil {
			return nil, err
		}
		scores = append(scores, score)
	}

	return scores, nil
}
