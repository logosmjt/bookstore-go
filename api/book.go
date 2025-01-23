package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/logosmjt/bookstore-go/db/sqlc"
	"github.com/logosmjt/bookstore-go/token"
)

type createBookRequest struct {
	Title         string    `json:"title" binding:"required"`
	Author        string    `json:"author" binding:"required,alphanum"`
	Price         int64     `json:"price" binding:"required"`
	Description   string    `json:"description" binding:"required"`
	CoverImageUrl string    `json:"cover_image_url" binding:"required"`
	PublishedDate time.Time `json:"published_date" binding:"required"`
}

func (server *Server) createBook(ctx *gin.Context) {
	var req createBookRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	args := db.CreateBookParams{
		Title:         req.Title,
		Author:        req.Author,
		Price:         req.Price,
		Description:   req.Description,
		CoverImageUrl: req.CoverImageUrl,
		PublishedDate: req.PublishedDate,
		UserID: pgtype.Int8{
			Int64: authPayload.Userid,
			Valid: true,
		},
	}
	book, err := server.store.CreateBook(ctx, args)
	if err != nil {
		errCode := db.ErrorCode(err)
		if errCode == db.ForeignKeyViolation || errCode == db.UniqueViolation {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, book)
}

type listBookRequest struct {
	PageNo   int32 `form:"page_no" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listBooks(ctx *gin.Context) {
	var req listBookRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	args := db.ListBooksParams{
		UserID: pgtype.Int8{
			Int64: authPayload.Userid,
			Valid: true,
		},
		Limit:  req.PageSize,
		Offset: (req.PageNo - 1) * req.PageSize,
	}

	books, err := server.store.ListBooks(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, books)
}
