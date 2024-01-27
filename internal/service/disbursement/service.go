package disbursementservice

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hansengotama/disbursement/internal/domain"
	"github.com/hansengotama/disbursement/internal/lib/postgres"
	cacherepo "github.com/hansengotama/disbursement/internal/repository/cache"
	disbursementrepo "github.com/hansengotama/disbursement/internal/repository/disbursement"
	"github.com/hansengotama/disbursement/internal/repository/disbursementaccount"
	"github.com/hansengotama/disbursement/internal/repository/paymentprovider"
	walletrepo "github.com/hansengotama/disbursement/internal/repository/wallet"
	"time"
)

type IDisbursementService interface {
	RequestDisbursement(param RequestDisbursementParam) error
}

type Dependency struct {
	WalletRepository              walletrepo.IWalletRepository
	PaymentProviderRepository     paymentproviderrepo.IPaymentProviderRepository
	DisbursementAccountRepository disbursementaccountrepo.IDisbursementAccountRepository
	DisbursementRepository        disbursementrepo.IDisbursementRepository
	CacheRepository               cacherepo.ICacheRepository
}

type DisbursementService struct {
	Dependency Dependency
}

type RequestDisbursementParam struct {
	Context                 context.Context
	UserID                  int
	Amount                  float64
	DisbursementAccountGUID uuid.UUID
}

func NewDisbursementService(dependency Dependency) IDisbursementService {
	return DisbursementService{
		Dependency: dependency,
	}
}

func (s DisbursementService) RequestDisbursement(param RequestDisbursementParam) error {
	ongoingRequestDisbursementKey := fmt.Sprintf("ongoing_request_disbursement_with_user_id_%v", param.UserID)
	ongoingRequestDisbursement, _ := s.Dependency.CacheRepository.Get(cacherepo.GetParam{
		Context: param.Context,
		Key:     ongoingRequestDisbursementKey,
	}).Result()

	if ongoingRequestDisbursement == "true" {
		return errors.New("only one ongoing request for disbursement is allowed. please try again later.")
	}

	s.Dependency.CacheRepository.Set(cacherepo.SetParam{
		Context: param.Context,
		Key:     ongoingRequestDisbursementKey,
		Value:   "true",
		TTL:     5 * time.Minute,
	})

	defer s.Dependency.CacheRepository.Set(cacherepo.SetParam{
		Context: param.Context,
		Key:     ongoingRequestDisbursementKey,
		Value:   "false",
		TTL:     5 * time.Minute,
	})

	conn := postgres.GetConnection()
	currentBalance, err := s.Dependency.WalletRepository.GetWalletBalanceByUserID(walletrepo.GetWalletBalanceByUserIDParam{
		Context:  param.Context,
		Executor: conn,
		UserID:   param.UserID,
	})
	if err != nil {
		return errors.New("error on get wallet balance")
	}

	disbursementAccount, err := s.Dependency.DisbursementAccountRepository.GetByGUID(disbursementaccountrepo.GetByGUIDParam{
		Context:  param.Context,
		Executor: conn,
		GUID:     param.DisbursementAccountGUID,
	})
	if err != nil {
		return errors.New("error on get disbursement account")
	}

	adminFee, err := s.Dependency.PaymentProviderRepository.GetAdminFeeByGUID(paymentproviderrepo.GetAdminFeeByGUIDParam{
		Context:  param.Context,
		Executor: conn,
		GUID:     disbursementAccount.PaymentProviderGUID,
	})

	amountWithFee := param.Amount + adminFee
	remainingBalance := currentBalance - amountWithFee
	if remainingBalance < 0 {
		return errors.New("insufficient funds: cannot disburse the specified amount")
	}

	tx, err := conn.Begin()
	if err != nil {
		return err
	}

	err = s.Dependency.WalletRepository.UpdateWalletBalanceByUserID(walletrepo.UpdateWalletBalanceByUserIDParam{
		Context:  param.Context,
		Executor: tx,
		UserID:   param.UserID,
		Balance:  remainingBalance,
	})
	if err != nil {
		tx.Rollback()
		return errors.New("error on update wallet balance")
	}

	err = s.Dependency.DisbursementRepository.Insert(disbursementrepo.InsertDisbursementParam{
		Context:                 param.Context,
		Executor:                tx,
		UserID:                  param.UserID,
		DisbursementAccountGUID: param.DisbursementAccountGUID,
		PaymentProviderGUID:     disbursementAccount.PaymentProviderGUID,
		AccountName:             disbursementAccount.Name,
		AccountNumber:           disbursementAccount.Number,
		AdminFee:                adminFee,
		Amount:                  param.Amount,
		AmountWithFee:           amountWithFee,
		Status:                  domain.DisbursementStatusPending,
	})
	if err != nil {
		tx.Rollback()
		return errors.New("error on insert disbursement")
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return errors.New("please try again")
	}

	return nil
}
