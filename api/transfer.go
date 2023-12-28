package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/TruongCongToan/simple_bank/token"

	db "github.com/TruongCongToan/simple_bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountId int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountId   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	fromAccount, valid := server.isValidAccount(ctx, req.FromAccountId, req.Currency)
	if !valid {
		return
	}

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if fromAccount.Owner != payload.Username {
		err := errors.New("from account does not belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	}

	_, valid = server.isValidAccount(ctx, req.ToAccountId, req.Currency)
	if !valid {
		return
	}

	arg := db.TransferTxParams{
		FromAccountId: req.FromAccountId,
		ToAccountId:   req.ToAccountId,
		Amount:        req.Amount,
	}

	result, err := server.store.TranferTx(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (server *Server) isValidAccount(ctx *gin.Context, accountId int64, currency string) (db.Account, bool) {

	account, error := server.store.GetAccount(ctx, accountId)
	if error != nil {
		if error == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, error.Error())
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, error.Error())
		return account, false
	}

	if account.Currency != currency {
		error := fmt.Sprintf("Account %d is not a valid currency: %s vs %s", account.ID, account.Currency, currency)

		ctx.JSON(http.StatusBadRequest, error)
		return account, false
	}
	return account, true
}
