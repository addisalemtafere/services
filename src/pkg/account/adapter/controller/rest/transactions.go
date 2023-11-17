package rest

import (
	"auth/src/pkg/account/core/entity"
	"auth/src/pkg/auth/adapter/controller/procedure"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	Id     uuid.UUID `json:"id"`
	From   Account   `json:"from"`
	To     Account   `json:"to"`
	Amount float64   `json:"amount"`
	Type   string    `json:"type"`
	Date   time.Time `json:"date"`
}

func NewTransactionFromEntity(i entity.Transaction) Transaction {
	return Transaction{
		Id:   i.Id,
		From: NewAccountFromEntity(i.From),
		To:   NewAccountFromEntity(i.To),
		Type: string(i.Type),
		Date: i.CreatedAt,
	}
}

func (controller Controller) GetRequestTransaction(w http.ResponseWriter, r *http.Request) {

	controller.log.Println("Adding Transaction")
	// Authenticate (AuthN)

	if len(strings.Split(r.Header.Get("Authorization"), " ")) != 2 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Please provide an authentication token in header",
			},
		}, http.StatusUnauthorized)
		return
	}

	// Validate token
	token := strings.Split(r.Header.Get("Authorization"), " ")[1]

	_, err := controller.auth.GetCheckAuth(token)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(procedure.Error).Type,
				Message: err.(procedure.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}

	// Check access (AuthZ)

	// Parse
	type Request struct {
		From       uuid.UUID `json:"from"`
		Recipients []struct {
			Account uuid.UUID `json:"account"`
			Amount  float64   `json:"amount"`
		} `json:"recipients"`
		Type   string `json:"type"`
		Medium string `json:"medium"`
	}

	var req Request

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	// Usecase

	// Response

	// Usecase [CREATE TRANSACTION]
	// txn, err := controller.interactor.CreateTransaction(session.User.Id, req.From, req.To, req.Amount, req.Type, req.Details.Token)
	// if err != nil {
	// 	SendJSONResponse(w, Response{
	// 		Success: false,
	// 		Error: &Error{
	// 			Type:    err.(usecase.Error).Type,
	// 			Message: err.(usecase.Error).Message,
	// 		},
	// 	}, http.StatusBadRequest)
	// 	return
	// }

	// SendJSONResponse(w, Response{
	// 	Success: true,
	// 	Data:    NewTransactionFromEntity(*txn),
	// }, http.StatusOK)

	w.Write([]byte("Received Request"))
}

// func (controller Controller) GetTransactions(w http.ResponseWriter, r *http.Request) {
// 	// Request
// 	type Request struct {
// 		Token string
// 	}

// 	req := Request{}

// 	token := r.Header.Get("Authorization")

// 	if len(strings.Split(token, " ")) == 2 {
// 		req.Token = strings.Split(token, " ")[1]
// 	}

// 	// Authenticate user
// 	session, err := controller.auth.GetCheckAuth(req.Token)
// 	if err != nil {
// 		SendJSONResponse(w, Response{
// 			Success: false,
// 			Error: &Error{
// 				Type:    err.(auth.Error).Type,
// 				Message: err.(auth.Error).Message,
// 			},
// 		}, http.StatusUnauthorized)
// 		return
// 	}

// 	// Get user id
// 	userId := session.User.Id

// 	// Get Transactions
// 	txns, err := controller.interactor.GetUserTransactions(userId)
// 	if err != nil {
// 		SendJSONResponse(w, Response{
// 			Success: false,
// 			Error: &Error{
// 				Type:    err.(usecase.Error).Type,
// 				Message: err.(usecase.Error).Message,
// 			},
// 		}, http.StatusBadRequest)
// 		return
// 	}

// 	// Map transactions to present
// 	var _txns []Transaction = make([]Transaction, 0)
// 	for i := 0; i < len(txns); i++ {
// 		_txns = append(_txns, NewTransactionFromEntity(txns[i]))
// 	}

// 	SendJSONResponse(w, Response{
// 		Success: true,
// 		Data:    _txns,
// 	}, http.StatusOK)
// }

// func (controller Controller) GetVerifyTransaction(w http.ResponseWriter, r *http.Request) {
// 	// Authenticate request

// 	controller.log.Println("Adding Transaction")

// 	// Check header
// 	if len(strings.Split(r.Header.Get("Authorization"), " ")) != 2 {
// 		SendJSONResponse(w, Response{
// 			Success: false,
// 			Error: &Error{
// 				Type:    "UNAUTHORIZED",
// 				Message: "Please provide an authentication token in header",
// 			},
// 		}, http.StatusUnauthorized)
// 		return
// 	}

// 	// Validate token
// 	token := strings.Split(r.Header.Get("Authorization"), " ")[1]

// 	_, err := controller.auth.GetCheckAuth(token)
// 	if err != nil {
// 		SendJSONResponse(w, Response{
// 			Success: false,
// 			Error: &Error{
// 				Type:    err.(procedure.Error).Type,
// 				Message: err.(procedure.Error).Message,
// 			},
// 		}, http.StatusUnauthorized)
// 		return
// 	}

// 	// Parse request
// 	type Request struct {
// 		Transaction uuid.UUID `json:"transaction"`
// 		Amount      float64   `json:"amount"`
// 	}

// 	var req Request

// 	decoder := json.NewDecoder(r.Body)
// 	err = decoder.Decode(&req)
// 	if err != nil {
// 		SendJSONResponse(w, Response{
// 			Success: false,
// 			Error: &Error{
// 				Type:    "INVALID_REQUEST",
// 				Message: err.Error(),
// 			},
// 		}, http.StatusBadRequest)
// 		return
// 	}

// 	defer r.Body.Close()

// 	// Usecase
// }
