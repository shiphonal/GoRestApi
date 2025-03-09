package sqlite

import (
	"GoServise/internal/storage"
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}

	// stmt - подготовительный запрос, на данном этапе мы его подготавливаем, после будем выполнять
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS url (
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL); 
	CREATE INDEX IF NOT EXISTS url_alias_idx ON url(alias);`)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s : %w", op, err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
		}
	}(stmt)

	// Exec - выполнит подготовительный запрос и подставит данные
	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		// Проверка является ли ошибка - ошибкой повторения псевдонимов (если alias повторяется, то для него у нас есть оформленная ошибка)
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return 0, fmt.Errorf("%s : %w", op, storage.ErrURLExists)
			}
		}
		/*if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s : %w", op, storage.ErrURLExists)
		}*/
		return 0, fmt.Errorf("%s : %w", op, err)
	}

	// берём id
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s : %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s, %w", op, err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
		}
	}(stmt)

	var url string
	err = stmt.QueryRow(alias).Scan(&url)
	if err != nil {
		if errors.Is(err, storage.ErrURLExists) {
			return "", storage.ErrURLNotFound
		}
		return "", fmt.Errorf("%s : %w", op, err)
	}

	return url, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.sqlite.DeleteURL"

	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias = ?")
	if err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
		}
	}(stmt)

	res, err := stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}

	if affected == 0 {
		return fmt.Errorf("%s : %w", op, storage.ErrURLNotFound)
	}
	return nil
}

func (s *Storage) PatchURL(oldAlias string, newAlias string) (bool, error) {
	const op = "storage.sqlite.PatchURL"

	stmt, err := s.db.Prepare("UPDATE url SET alias = ? WHERE alias = ?")
	if err != nil {
		return false, fmt.Errorf("%s : %w", op, err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
		}
	}(stmt)

	res, err := stmt.Exec(newAlias, oldAlias)
	if err != nil {
		return false, fmt.Errorf("%s : %w", op, err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("%s : %w", op, err)
	}
	if affected == 0 {
		return false, fmt.Errorf("%s : %s", op, "patch the same alias")
	}
	return true, nil
}
