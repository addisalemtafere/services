package psql

import (
	"auth/src/pkg/account/core/entity"
	"context"

	"github.com/google/uuid"
)

func (repo PsqlRepo) StoreTransaction(txn entity.Transaction) error {
	tx, err := repo.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	// Store Transaction
	_, err = tx.Exec(`
	INSERT INTO accounts.transactions (id, "from", "to", "type", "reference", verified, created_at)
	VALUES ($1::UUID, $2::UUID, $3::UUID, $4, $5, $6, $7)
	`, txn.Id, txn.From.Id, txn.To.Id, txn.Type, txn.Reference, txn.Verified, txn.CreatedAt)

	if err != nil {
		tx.Rollback()
		return err
	}

	// Store transaction details
	switch txn.Type {
	case entity.REPLENISHMENT:
		{
			txnDetail := txn.Details.(entity.Replenishment)
			_, err = tx.Exec(`
			INSERT INTO accounts.a2a_transactions (transaction_id, amount)
			VALUES ($1::UUID,$2)
			`, txn.Id, txnDetail.Amount)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return err
}

func (repo PsqlRepo) FindTransactionsByUserId(id uuid.UUID) ([]entity.Transaction, error) {
	var txns []entity.Transaction = make([]entity.Transaction, 0)

	rows, err := repo.db.Query(`
	SELECT 
		transactions.id, transactions.type, transactions.created_at, transactions.updated_at,
		"tag".id, "tag".name, "tag".color,
		"from".id, "from".title, "from".type, "from".default, "from".user,
		"to".id, "to".title, "to".type, "to".default, "to".user
	FROM accounts.transactions
	LEFT JOIN accounts.accounts as "from" ON "from".id = transactions.from
	LEFT JOIN accounts.accounts as "to" ON "to".id = transactions.to
	LEFT JOIN accounts.tags as "tag" ON "tag".id = transactions.tag
	WHERE "from".user = $1::UUID OR "to".user = $1::UUID;
	`, id)

	if err != nil {
		repo.log.Println(err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var txn entity.Transaction
		err := rows.Scan(&txn.Id, &txn.Type, &txn.CreatedAt, &txn.UpdatedAt,
			&txn.From.Id, &txn.From.Title, &txn.From.Type, &txn.From.Default, &txn.From.User.Id,
			&txn.To.Id, &txn.To.Title, &txn.To.Type, &txn.To.Default, &txn.To.User.Id,
		)

		if err == nil {
			// Fetch txn details
			switch txn.Type {
			case entity.REPLENISHMENT:
				{
					txn.Details = nil
				}
			}

			txns = append(txns, txn)
		}
	}

	return txns, nil
}

func (repo PsqlRepo) FindTransactionById(id uuid.UUID) (*entity.Transaction, error) {
	var txn entity.Transaction

	err := repo.db.QueryRow(`
	SELECT 
		transactions.id, transactions.type, transactions.created_at, transactions.updated_at,
		"tag".id, "tag".name, "tag".color,
		"from".id, "from".title, "from".type, "from".default, "from".user,
		"to".id, "to".title, "to".type, "to".default, "to".user
	FROM accounts.transactions
	LEFT JOIN accounts.accounts as "from" ON "from".id = transactions.from
	LEFT JOIN accounts.accounts as "to" ON "to".id = transactions.to
	LEFT JOIN accounts.tags as "tag" ON "tag".id = transactions.tag
	WHERE "from".user = $1::UUID OR "to".user = $1::UUID;
	`).Scan(
		&txn.Id, &txn.Type, &txn.CreatedAt, &txn.UpdatedAt,
		&txn.From.Id, &txn.From.Title, &txn.From.Type, &txn.From.Default, &txn.From.User.Id,
		&txn.To.Id, &txn.To.Title, &txn.To.Type, &txn.To.Default, &txn.To.User.Id,
	)

	return &txn, err
}
