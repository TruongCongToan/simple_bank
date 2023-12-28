package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/TruongCongToan/simple_bank/db/sqlc"
	"github.com/TruongCongToan/simple_bank/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) CreateAccount(ctx *gin.Context) {
	var req createAccountRequest

	if error := ctx.ShouldBindJSON(&req); error != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(error))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			switch pqError.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) GetAccount(ctx *gin.Context) {
	var req getAccountRequest
	if error := ctx.ShouldBindUri(&req); error != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(error))
		return
	}

	account, error := server.store.GetAccount(ctx, req.ID)

	if error != nil {
		if error == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(error))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(error))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Username != account.Owner {
		err := errors.New(" account does not belong to authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getListAccountRequest struct {
	PageId   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1,max=10"`
}

func (server *Server) GetListAccount(ctx *gin.Context) {
	var req getListAccountRequest
	if error := ctx.ShouldBindQuery(&req); error != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(error))
		return
	}

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.ListAccountsParams{
		Owner:  payload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageId - 1) * req.PageSize,
	}
	accounts, error := server.store.ListAccounts(ctx, arg)

	if error != nil {
		if error == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(error))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(error))
		return
	}
	ctx.JSON(http.StatusOK, accounts)
}
