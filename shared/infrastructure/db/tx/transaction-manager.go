package tx

import (
	"api-your-accounts/shared/domain/transaction"
	"errors"

	"gorm.io/gorm"
)

var (
	ErrInvalidTX = errors.New("invalid type TX")
)

type gormTransaction struct {
	tx *gorm.DB
}

func (m *gormTransaction) Set(tx any) error {
	var ok bool
	if m.tx, ok = tx.(*gorm.DB); !ok {
		return ErrInvalidTX
	}

	return nil
}

func (m *gormTransaction) Get() any {
	return m.tx
}

type gormTransactionManager struct {
	db *gorm.DB
}

func (tm *gormTransactionManager) Transaction(fc func(tx transaction.Transaction) error) error {
	return tm.db.Transaction(func(db *gorm.DB) error {
		tx := &gormTransaction{}
		tx.Set(db)
		return fc(tx)
	})
}

func NewTransactionManager(db *gorm.DB) transaction.TransactionManager {
	return &gormTransactionManager{db}
}
