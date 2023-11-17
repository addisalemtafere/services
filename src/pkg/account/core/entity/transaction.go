package entity

import (
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	REPLENISHMENT TransactionType = "REPLENISHMENT"
	P2P           TransactionType = "P2P"
	SALE          TransactionType = "SALE"
	BILL_PAYMENT  TransactionType = "BILL_PAYMENT"
	SETTLEMENT    TransactionType = "SETTLEMENT"
)

type TransactionMedium string

const (
	LAKIPAY     TransactionMedium = "LAKIPAY"
	CYBERSOURCE TransactionMedium = "CYBERSOURCE"
	ETHSWITCH   TransactionMedium = "ETHSWITCH"
)

type Transaction struct {
	Id        uuid.UUID
	From      Account
	To        Account
	Type      TransactionType
	Medium    TransactionMedium
	Reference string
	Comment   string
	Tag       Tag
	Verified  bool
	TTL       int64
	Details   interface{}
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Replenishment struct {
	Amount float64
}

type BatchTransaction struct {
	Id   uuid.UUID
	From []struct {
		Account Account
		Amount  float64
	}
	To []struct {
		Account Account
		Amount  float64
	}
	Amount       float64
	Transactions []Transaction
	Verified     bool
	TTL          int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
