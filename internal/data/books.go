package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

// By default, the keys in the JSON object are equal to the field names in the struct ( ID,
// CreatedAt, Title and so on).
type Book struct {
	ID int64 `json:"id"` // Unique integer ID for the movie

	Title    string   `json:"title"` // Movie title
	Author   string   `json:"author"`
	Year     int32    `json:"year,omitempty"` // Movie release year, "omitempty" - hide from response if empty
	Language string   `json:"language"`
	Genres   []string `json:"genres,omitempty"` // Slice of genres for the movie (romance, comedy, etc.)
	Price    int32    `json:"price"`
	Quantity int32    `json:"quantity"`
	Version  int32    `json:"version"`
	// The version number starts at 1 and will be incremented each
	// time the movie information is updated
}

// Идея для покупки будет,человек пишет айди и количество после это книга обновиться,после чего  у него пропадет определенное
// кол-во количеств и сделаем валидатор  что бы проверяел,если книга колво будет равно или меньше 0,то высветиться ошибка
type BookModel struct {
	DB *sql.DB
}

func (m BookModel) AddMovieInDB(book *Book) error {
	query := `
		INSERT INTO books (title, author, year, language,genres,price,quantity)
		VALUES ($1, $2, $3, $4,$5 ,$6,$7)
		RETURNING id,  version`
	args := []any{book.Title, book.Author, book.Year, book.Language, pq.Array(book.Genres), book.Price, book.Quantity}
	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Use QueryRowContext() and pass the context as the first argument.
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&book.ID, &book.Version)
}

// 1 из 2/3 вещей что может сделать обычный юзер,просто смотреть есть ли такая книга
func (m BookModel) GetInfo(id int64) (*Book, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	// Remove the pg_sleep(10) clause.
	query := `
	SELECT id,title, author,year, language, genres, price,quantity,version
	FROM books
	WHERE id = $1`
	var book Book
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Remove &[]byte{} from the first Scan() destination.
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.Year,
		&book.Language,
		pq.Array(&book.Genres),
		&book.Price,
		&book.Quantity,
		&book.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &book, nil
}

func (m BookModel) Update(book *Book) error {
	query := `
	UPDATE books
	SET title = $1, author = $2, year = $3, genres = $4, quantity=$5,version = version + 1
	WHERE id = $6 AND version = $7
	RETURNING version`
	args := []any{
		book.Title,
		book.Author,
		book.Year,
		pq.Array(book.Genres),
		book.Quantity,
		book.ID,
		book.Version,
	}
	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// method for deleting a specific record from the movies table.
	// Use QueryRowContext() and pass the context as the first argument.
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&book.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m BookModel) Buying(book *Book) error {
	query := `
	UPDATE books
	SET  quantity=$1,version = version + 1
	WHERE id = $2 AND version = $3
	RETURNING version`
	args := []any{
		book.Title,
		book.Author,
		book.Year,
		pq.Array(book.Genres),
		book.Quantity,
		book.ID,
		book.Version,
	}
	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// method for deleting a specific record from the movies table.
	// Use QueryRowContext() and pass the context as the first argument.
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&book.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m BookModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `
DELETE FROM books
WHERE id = $1`
	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Use ExecContext() and pass the context as the first argument.
	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
