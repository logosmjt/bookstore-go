package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
	mockdb "github.com/logosmjt/bookstore-go/db/mock"
	db "github.com/logosmjt/bookstore-go/db/sqlc"
	"github.com/logosmjt/bookstore-go/token"
	"github.com/logosmjt/bookstore-go/util"
	"github.com/stretchr/testify/require"
)

func TestCreateBookAPI(t *testing.T) {
	user, _ := randomUser(t)
	book := randomBook(user.ID)

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"title":           book.Title,
				"author":          book.Author,
				"price":           book.Price,
				"description":     book.Description,
				"cover_image_url": book.CoverImageUrl,
				"published_date":  book.PublishedDate,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Name, user.ID, user.Role, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.CreateBookParams{
					Title:         book.Title,
					Author:        book.Author,
					Price:         book.Price,
					Description:   book.Description,
					CoverImageUrl: book.CoverImageUrl,
					PublishedDate: book.PublishedDate,
					UserID:        book.UserID,
				}
				store.EXPECT().
					CreateBook(gomock.Any(), gomock.Eq(args)).
					Times(1).
					Return(book, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			body: gin.H{
				"title":           book.Title,
				"author":          book.Author,
				"price":           book.Price,
				"description":     book.Description,
				"cover_image_url": book.CoverImageUrl,
				"published_date":  book.PublishedDate,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateBook(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"title":           book.Title,
				"author":          book.Author,
				"price":           book.Price,
				"description":     book.Description,
				"cover_image_url": book.CoverImageUrl,
				"published_date":  book.PublishedDate,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Name, user.ID, user.Role, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateBook(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Book{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/books"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}

}

func TestListBooksAPI(t *testing.T) {
	user, _ := randomUser(t)

	n := 5
	books := make([]db.Book, n)
	for i := 0; i < n; i++ {
		books[i] = randomBook(user.ID)
	}

	type Query struct {
		pageNo   int
		pageSize int
	}

	testCases := []struct {
		name          string
		query         Query
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{{
		name: "OK",
		query: Query{
			pageNo:   1,
			pageSize: n,
		},
		setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Name, user.ID, user.Role, time.Minute)
		},
		buildStubs: func(store *mockdb.MockStore) {
			arg := db.ListBooksParams{
				UserID: pgtype.Int8{Int64: user.ID},
				Limit:  int32(n),
				Offset: 0,
			}

			store.EXPECT().
				ListBooks(gomock.Any(), gomock.Eq(arg)).
				Times(1).
				Return(books, nil)
		},
		checkResponse: func(recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusOK, recorder.Code)
			requireBodyMatchBooks(t, recorder.Body, books)
		},
	},
		{
			name: "NoAuthorization",
			query: Query{
				pageNo:   1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListBooks(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			query: Query{
				pageNo:   1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Name, user.ID, user.Role, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListBooks(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Book{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidPageID",
			query: Query{
				pageNo:   -1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Name, user.ID, user.Role, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListBooks(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageSize",
			query: Query{
				pageNo:   1,
				pageSize: 100000,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Name, user.ID, user.Role, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListBooks(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		}}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := "/books"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add query parameters to request URL
			q := request.URL.Query()
			q.Add("page_no", fmt.Sprintf("%d", tc.query.pageNo))
			q.Add("page_size", fmt.Sprintf("%d", tc.query.pageSize))
			request.URL.RawQuery = q.Encode()

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}

}

func randomBook(userid int64) db.Book {
	return db.Book{
		ID:            util.RandomInt(1, 1000),
		Title:         util.RandomString(6),
		Author:        util.RandomUserName(),
		Price:         util.RandomInt(1, 1000),
		Description:   util.RandomString(10),
		CoverImageUrl: util.RandomString(10),
		PublishedDate: util.RandomTime(),
		UserID:        pgtype.Int8{Int64: userid, Valid: true},
	}
}

func requireBodyMatchBooks(t *testing.T, body *bytes.Buffer, accounts []db.Book) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotBooks []db.Book
	err = json.Unmarshal(data, &gotBooks)
	require.NoError(t, err)
	require.Equal(t, accounts, gotBooks)
}
