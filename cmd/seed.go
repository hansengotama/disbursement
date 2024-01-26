package seed

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/hansengotama/disbursement/internal/lib/postgres"
	accesstokenrepo "github.com/hansengotama/disbursement/internal/repository/accesstoken"
	"github.com/hansengotama/disbursement/internal/repository/disbursementaccount"
	"github.com/hansengotama/disbursement/internal/repository/paymentprovider"
	userrepo "github.com/hansengotama/disbursement/internal/repository/user"
	walletrepo "github.com/hansengotama/disbursement/internal/repository/wallet"
	"time"
)

const (
	userID                               = 1
	accessToken                          = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
	danaPaymentProviderStrGUID           = "8c2db1f4-0e63-4f47-8f77-2ac99acdbbc7"
	hansenDanaDisbursementAccountStrGUID = "6c2db1f4-0e63-4f47-8f77-2ac99acdbbc7"
)

func Execute() {
	conn := postgres.GetConnection()
	tx, err := conn.Begin()
	if err != nil {
		panic(err)
	}

	defer func() {
		if r := recover(); r != nil {
			// If a panic occurs, roll back the transaction
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				panic(rollbackErr)
			}

			// Re-panic after rolling back to propagate the original error
			panic(r)
		}
	}()

	seedUser(tx)
	seedAccessToken(tx)
	seedWallet(tx)
	seedPaymentProvider(tx)
	seedDisbursementAccount(tx)

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully seeding!")
}

func seedUser(executor postgres.SQLExecutor) {
	r := userrepo.UserDB{}

	err := r.Insert(userrepo.InsertUserParam{
		Context:  context.TODO(),
		Executor: executor,
		Name:     "Hansen",
	})
	if err != nil {
		panic(err)
	}
}

func seedAccessToken(executor postgres.SQLExecutor) {
	r := accesstokenrepo.AccessTokenDB{}

	oneDay := time.Hour * 24
	err := r.Insert(accesstokenrepo.InsertAccessTokenParam{
		Context:        context.TODO(),
		Executor:       executor,
		Token:          accessToken,
		UserID:         userID,
		ExpirationTime: time.Now().Add(oneDay),
	})
	if err != nil {
		panic(err)
	}
}

func seedWallet(executor postgres.SQLExecutor) {
	r := walletrepo.WalletDB{}

	err := r.Insert(walletrepo.InsertWalletParam{
		Context:  context.TODO(),
		Executor: executor,
		UserID:   userID,
		Balance:  10000,
	})
	if err != nil {
		panic(err)
	}
}

func seedPaymentProvider(executor postgres.SQLExecutor) {
	r := paymentproviderrepo.PaymentProviderDB{}

	danaPaymentProviderGUID, err := uuid.Parse(danaPaymentProviderStrGUID)
	if err != nil {
		panic(err)
	}

	err = r.Insert(paymentproviderrepo.InsertPaymentProviderParam{
		Context:  context.TODO(),
		Executor: executor,
		GUID:     danaPaymentProviderGUID,
		Name:     "Dana",
		AdminFee: 200,
		Type:     "ewallet",
	})
	if err != nil {
		panic(err)
	}
}

func seedDisbursementAccount(executor postgres.SQLExecutor) {
	r := disbursementaccountrepo.DisbursementAccountDB{}

	danaPaymentProviderGUID, err := uuid.Parse(danaPaymentProviderStrGUID)
	if err != nil {
		panic(err)
	}

	hansenDanaDisbursementAccountGUID, err := uuid.Parse(hansenDanaDisbursementAccountStrGUID)
	if err != nil {
		panic(err)
	}

	err = r.Insert(disbursementaccountrepo.InsertDisbursementAccountParam{
		Context:             context.TODO(),
		Executor:            executor,
		GUID:                hansenDanaDisbursementAccountGUID,
		UserID:              userID,
		PaymentProviderGUID: danaPaymentProviderGUID,
		Name:                "Hansen",
		Number:              "628111814032",
	})
	if err != nil {
		panic(err)
	}
}
