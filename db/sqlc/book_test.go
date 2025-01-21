package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/logosmjt/bookstore-go/util"
	"github.com/stretchr/testify/require"
)

func createRandomBook(t *testing.T) Book {
	user := createRandomUser(t)
	arg := CreateBookParams{
		Title:         util.RandomString(6),
		Author:        util.RandomUserName(),
		Price:         util.RandomInt(0, 1000),
		Description:   util.RandomString(20),
		CoverImageUrl: util.RandomString(10),
		PublishedDate: util.RandomTime(),
		UserID:        pgtype.Int8{Int64: user.ID, Valid: true},
	}

	book, err := testStore.CreateBook(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Title, book.Title)
	require.Equal(t, arg.Author, book.Author)
	require.Equal(t, arg.Price, book.Price)
	require.Equal(t, arg.Description, book.Description)
	require.Equal(t, arg.CoverImageUrl, book.CoverImageUrl)
	require.Equal(t, arg.PublishedDate.UTC(), book.PublishedDate.UTC())
	require.Equal(t, arg.UserID.Int64, user.ID)
	require.NotZero(t, book.CreatedAt)

	return book
}

func TestCreateBook(t *testing.T) {
	createRandomBook(t)
}

func TestUpdateBookTitle(t *testing.T) {
	book := createRandomBook(t)
	newTitle := util.RandomString(6)

	arg := UpdateBookParams{
		ID: book.ID,
		Title: pgtype.Text{
			String: newTitle,
			Valid:  true,
		},
	}
	updatedBook, err := testStore.UpdateBook(context.Background(), arg)

	require.NoError(t, err)
	require.NotEqual(t, book.Title, updatedBook.Title)
	require.Equal(t, newTitle, updatedBook.Title)
	require.Equal(t, book.Author, updatedBook.Author)
	require.Equal(t, book.Price, updatedBook.Price)
	require.Equal(t, book.Description, updatedBook.Description)
	require.Equal(t, book.CoverImageUrl, updatedBook.CoverImageUrl)
	require.Equal(t, book.PublishedDate.UTC(), updatedBook.PublishedDate.UTC())
	require.Equal(t, book.UserID, updatedBook.UserID)
}

func TestUpdateBookAuthor(t *testing.T) {
	book := createRandomBook(t)
	newAuthor := util.RandomUserName()

	arg := UpdateBookParams{
		ID: book.ID,
		Author: pgtype.Text{
			String: newAuthor,
			Valid:  true,
		},
	}
	updatedBook, err := testStore.UpdateBook(context.Background(), arg)

	require.NoError(t, err)
	require.NotEqual(t, book.Author, updatedBook.Author)
	require.Equal(t, newAuthor, updatedBook.Author)
	require.Equal(t, book.Title, updatedBook.Title)
	require.Equal(t, book.Price, updatedBook.Price)
	require.Equal(t, book.Description, updatedBook.Description)
	require.Equal(t, book.CoverImageUrl, updatedBook.CoverImageUrl)
	require.Equal(t, book.PublishedDate.UTC(), updatedBook.PublishedDate.UTC())
	require.Equal(t, book.UserID, updatedBook.UserID)
}

func TestUpdateBook(t *testing.T) {
	book := createRandomBook(t)
	newTitle := util.RandomString(6)
	newAuthor := util.RandomUserName()
	newPrice := util.RandomInt(0, 1000)
	newDescription := util.RandomString(20)
	newCoverImageUrl := util.RandomString(20)
	newPublishedDate := util.RandomTime()

	arg := UpdateBookParams{
		ID: book.ID,
		Title: pgtype.Text{
			String: newTitle,
			Valid:  true,
		},
		Author: pgtype.Text{
			String: newAuthor,
			Valid:  true,
		},
		Price: pgtype.Int8{
			Int64: newPrice,
			Valid: true,
		},
		Description: pgtype.Text{
			String: newDescription,
			Valid:  true,
		},
		CoverImageUrl: pgtype.Text{
			String: newCoverImageUrl,
			Valid:  true,
		},
		PublishedDate: pgtype.Timestamptz{
			Time:  newPublishedDate,
			Valid: true,
		},
	}
	updatedBook, err := testStore.UpdateBook(context.Background(), arg)

	require.NoError(t, err)
	require.NotEqual(t, book.Title, updatedBook.Title)
	require.Equal(t, newTitle, updatedBook.Title)
	require.NotEqual(t, book.Author, updatedBook.Author)
	require.Equal(t, newAuthor, updatedBook.Author)
	require.NotEqual(t, book.Price, updatedBook.Price)
	require.Equal(t, newPrice, updatedBook.Price)
	require.NotEqual(t, book.Description, updatedBook.Description)
	require.Equal(t, newDescription, updatedBook.Description)
	require.NotEqual(t, book.CoverImageUrl, updatedBook.CoverImageUrl)
	require.Equal(t, newCoverImageUrl, updatedBook.CoverImageUrl)
	require.NotEqual(t, book.PublishedDate.UTC(), updatedBook.PublishedDate.UTC())
	require.Equal(t, newPublishedDate.UTC(), updatedBook.PublishedDate.UTC())

	require.Equal(t, book.UserID, updatedBook.UserID)
}
