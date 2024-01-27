package disbursementservice_test

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/hansengotama/disbursement/internal/lib/postgres"
	cacherepo "github.com/hansengotama/disbursement/internal/repository/cache"
	disbursementaccountrepo "github.com/hansengotama/disbursement/internal/repository/disbursementaccount"
	paymentproviderrepo "github.com/hansengotama/disbursement/internal/repository/paymentprovider"
	walletrepo "github.com/hansengotama/disbursement/internal/repository/wallet"
	disbursementservice "github.com/hansengotama/disbursement/internal/service/disbursement"
	"github.com/hansengotama/disbursement/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestDisbursementService_RequestDisbursement(t *testing.T) {
	disbursementAccountGUID := uuid.New()
	paymentProviderGUID := uuid.New()
	ctx := context.TODO()
	conn := postgres.GetConnection()

	testCases := []struct {
		description                 string
		param                       disbursementservice.RequestDisbursementParam
		mockWalletRepo              func() *mocks.IWalletRepository
		mockPaymentProviderRepo     func() *mocks.IPaymentProviderRepository
		mockDisbursementAccountRepo func() *mocks.IDisbursementAccountRepository
		mockDisbursementRepo        func() *mocks.IDisbursementRepository
		mockCacheRepo               func() *mocks.ICacheRepository
		expectedErr                 error
	}{
		{
			description: "when successfully request disbursement",
			param: disbursementservice.RequestDisbursementParam{
				Context:                 ctx,
				UserID:                  1,
				Amount:                  10000,
				DisbursementAccountGUID: disbursementAccountGUID,
			},
			mockWalletRepo: func() *mocks.IWalletRepository {
				mockWalletRepo := new(mocks.IWalletRepository)
				res := float64(20200)
				mockWalletRepo.On("GetWalletBalanceByUserID", walletrepo.GetWalletBalanceByUserIDParam{
					Context:  ctx,
					Executor: conn,
					UserID:   1,
				}).Return(res, nil).Once()

				mockWalletRepo.On("UpdateWalletBalanceByUserID", mock.AnythingOfType("UpdateWalletBalanceByUserIDParam")).Return(nil).Once()

				return mockWalletRepo
			},
			mockPaymentProviderRepo: func() *mocks.IPaymentProviderRepository {
				mockPaymentProviderRepo := new(mocks.IPaymentProviderRepository)
				res := float64(200)
				mockPaymentProviderRepo.On("GetAdminFeeByGUID", paymentproviderrepo.GetAdminFeeByGUIDParam{
					Context:  ctx,
					Executor: conn,
					GUID:     paymentProviderGUID,
				}).Return(res, nil).Once()

				return mockPaymentProviderRepo
			},
			mockDisbursementAccountRepo: func() *mocks.IDisbursementAccountRepository {
				mockDisbursementAccountRepo := new(mocks.IDisbursementAccountRepository)
				mockDisbursementAccountRepo.On("GetByGUID", disbursementaccountrepo.GetByGUIDParam{
					Context:  ctx,
					Executor: conn,
					GUID:     disbursementAccountGUID,
				}).Return(&disbursementaccountrepo.GetByGUIDParamRes{
					PaymentProviderGUID: paymentProviderGUID,
					Name:                "Hansen",
					Number:              "123",
				}, nil).Once()

				return mockDisbursementAccountRepo
			},
			mockDisbursementRepo: func() *mocks.IDisbursementRepository {
				mockDisbursementRepo := new(mocks.IDisbursementRepository)

				mockDisbursementRepo.On("Insert", mock.AnythingOfType("InsertDisbursementParam")).Return(nil).Once()

				return mockDisbursementRepo
			},
			mockCacheRepo: func() *mocks.ICacheRepository {
				mockCacheRepo := new(mocks.ICacheRepository)

				mockCacheRepo.On("Get", cacherepo.GetParam{
					Context: ctx,
					Key:     "ongoing_request_disbursement_with_user_id_1",
				}).Return("false", nil).Once()

				mockCacheRepo.On("Set", cacherepo.SetParam{
					Context: ctx,
					Key:     "ongoing_request_disbursement_with_user_id_1",
					Value:   "true",
					TTL:     5 * time.Minute,
				}).Return("OK", nil).Once()

				mockCacheRepo.On("Set", cacherepo.SetParam{
					Context: ctx,
					Key:     "ongoing_request_disbursement_with_user_id_1",
					Value:   "false",
					TTL:     5 * time.Minute,
				}).Return("OK", nil).Once()

				return mockCacheRepo
			},
			expectedErr: nil,
		},
		{
			description: "when failed request disbursement on has ongoing request disbursement",
			param: disbursementservice.RequestDisbursementParam{
				Context:                 ctx,
				UserID:                  1,
				Amount:                  10000,
				DisbursementAccountGUID: disbursementAccountGUID,
			},
			mockWalletRepo: func() *mocks.IWalletRepository {
				mockWalletRepo := new(mocks.IWalletRepository)

				return mockWalletRepo
			},
			mockPaymentProviderRepo: func() *mocks.IPaymentProviderRepository {
				mockPaymentProviderRepo := new(mocks.IPaymentProviderRepository)

				return mockPaymentProviderRepo
			},
			mockDisbursementAccountRepo: func() *mocks.IDisbursementAccountRepository {
				mockDisbursementAccountRepo := new(mocks.IDisbursementAccountRepository)

				return mockDisbursementAccountRepo
			},
			mockDisbursementRepo: func() *mocks.IDisbursementRepository {
				mockDisbursementRepo := new(mocks.IDisbursementRepository)

				return mockDisbursementRepo
			},
			mockCacheRepo: func() *mocks.ICacheRepository {
				mockCacheRepo := new(mocks.ICacheRepository)

				mockCacheRepo.On("Get", cacherepo.GetParam{
					Context: ctx,
					Key:     "ongoing_request_disbursement_with_user_id_1",
				}).Return("true", nil).Once()

				return mockCacheRepo
			},
			expectedErr: errors.New("only one ongoing request for disbursement is allowed. please try again later"),
		},
		{
			description: "when failed request disbursement on get wallet balance",
			param: disbursementservice.RequestDisbursementParam{
				Context:                 ctx,
				UserID:                  1,
				Amount:                  10000,
				DisbursementAccountGUID: disbursementAccountGUID,
			},
			mockWalletRepo: func() *mocks.IWalletRepository {
				mockWalletRepo := new(mocks.IWalletRepository)
				res := float64(0)
				mockWalletRepo.On("GetWalletBalanceByUserID", walletrepo.GetWalletBalanceByUserIDParam{
					Context:  ctx,
					Executor: conn,
					UserID:   1,
				}).Return(res, errors.New("error on get wallet balance by user id")).Once()

				return mockWalletRepo
			},
			mockPaymentProviderRepo: func() *mocks.IPaymentProviderRepository {
				mockPaymentProviderRepo := new(mocks.IPaymentProviderRepository)

				return mockPaymentProviderRepo
			},
			mockDisbursementAccountRepo: func() *mocks.IDisbursementAccountRepository {
				mockDisbursementAccountRepo := new(mocks.IDisbursementAccountRepository)

				return mockDisbursementAccountRepo
			},
			mockDisbursementRepo: func() *mocks.IDisbursementRepository {
				mockDisbursementRepo := new(mocks.IDisbursementRepository)

				return mockDisbursementRepo
			},
			mockCacheRepo: func() *mocks.ICacheRepository {
				mockCacheRepo := new(mocks.ICacheRepository)

				mockCacheRepo.On("Get", cacherepo.GetParam{
					Context: ctx,
					Key:     "ongoing_request_disbursement_with_user_id_1",
				}).Return("false", nil).Once()

				mockCacheRepo.On("Set", cacherepo.SetParam{
					Context: ctx,
					Key:     "ongoing_request_disbursement_with_user_id_1",
					Value:   "true",
					TTL:     5 * time.Minute,
				}).Return("OK", nil).Once()

				mockCacheRepo.On("Set", cacherepo.SetParam{
					Context: ctx,
					Key:     "ongoing_request_disbursement_with_user_id_1",
					Value:   "false",
					TTL:     5 * time.Minute,
				}).Return("OK", nil).Once()

				return mockCacheRepo
			},
			expectedErr: errors.New("error on get wallet balance"),
		},
		{
			description: "when failed request disbursement on get disbursement account",
			param: disbursementservice.RequestDisbursementParam{
				Context:                 ctx,
				UserID:                  1,
				Amount:                  10000,
				DisbursementAccountGUID: disbursementAccountGUID,
			},
			mockWalletRepo: func() *mocks.IWalletRepository {
				mockWalletRepo := new(mocks.IWalletRepository)
				res := float64(20200)
				mockWalletRepo.On("GetWalletBalanceByUserID", walletrepo.GetWalletBalanceByUserIDParam{
					Context:  ctx,
					Executor: conn,
					UserID:   1,
				}).Return(res, nil).Once()

				return mockWalletRepo
			},
			mockPaymentProviderRepo: func() *mocks.IPaymentProviderRepository {
				mockPaymentProviderRepo := new(mocks.IPaymentProviderRepository)

				return mockPaymentProviderRepo
			},
			mockDisbursementAccountRepo: func() *mocks.IDisbursementAccountRepository {
				mockDisbursementAccountRepo := new(mocks.IDisbursementAccountRepository)
				mockDisbursementAccountRepo.On("GetByGUID", disbursementaccountrepo.GetByGUIDParam{
					Context:  ctx,
					Executor: conn,
					GUID:     disbursementAccountGUID,
				}).Return(nil, errors.New("error on get disbursement account")).Once()

				return mockDisbursementAccountRepo
			},
			mockDisbursementRepo: func() *mocks.IDisbursementRepository {
				mockDisbursementRepo := new(mocks.IDisbursementRepository)

				return mockDisbursementRepo
			},
			mockCacheRepo: func() *mocks.ICacheRepository {
				mockCacheRepo := new(mocks.ICacheRepository)

				mockCacheRepo.On("Get", cacherepo.GetParam{
					Context: ctx,
					Key:     "ongoing_request_disbursement_with_user_id_1",
				}).Return("false", nil).Once()

				mockCacheRepo.On("Set", cacherepo.SetParam{
					Context: ctx,
					Key:     "ongoing_request_disbursement_with_user_id_1",
					Value:   "true",
					TTL:     5 * time.Minute,
				}).Return("OK", nil).Once()

				mockCacheRepo.On("Set", cacherepo.SetParam{
					Context: ctx,
					Key:     "ongoing_request_disbursement_with_user_id_1",
					Value:   "false",
					TTL:     5 * time.Minute,
				}).Return("OK", nil).Once()

				return mockCacheRepo
			},
			expectedErr: errors.New("error on get disbursement account"),
		},
		{
			description: "when failed request disbursement on get provider admin fee",
			param: disbursementservice.RequestDisbursementParam{
				Context:                 ctx,
				UserID:                  1,
				Amount:                  10000,
				DisbursementAccountGUID: disbursementAccountGUID,
			},
			mockWalletRepo: func() *mocks.IWalletRepository {
				mockWalletRepo := new(mocks.IWalletRepository)
				res := float64(20200)
				mockWalletRepo.On("GetWalletBalanceByUserID", walletrepo.GetWalletBalanceByUserIDParam{
					Context:  ctx,
					Executor: conn,
					UserID:   1,
				}).Return(res, nil).Once()

				return mockWalletRepo
			},
			mockPaymentProviderRepo: func() *mocks.IPaymentProviderRepository {
				mockPaymentProviderRepo := new(mocks.IPaymentProviderRepository)
				res := float64(0)
				mockPaymentProviderRepo.On("GetAdminFeeByGUID", paymentproviderrepo.GetAdminFeeByGUIDParam{
					Context:  ctx,
					Executor: conn,
					GUID:     paymentProviderGUID,
				}).Return(res, errors.New("error on get provider admin fee")).Once()

				return mockPaymentProviderRepo
			},
			mockDisbursementAccountRepo: func() *mocks.IDisbursementAccountRepository {
				mockDisbursementAccountRepo := new(mocks.IDisbursementAccountRepository)
				mockDisbursementAccountRepo.On("GetByGUID", disbursementaccountrepo.GetByGUIDParam{
					Context:  ctx,
					Executor: conn,
					GUID:     disbursementAccountGUID,
				}).Return(&disbursementaccountrepo.GetByGUIDParamRes{
					PaymentProviderGUID: paymentProviderGUID,
					Name:                "Hansen",
					Number:              "123",
				}, nil).Once()

				return mockDisbursementAccountRepo
			},
			mockDisbursementRepo: func() *mocks.IDisbursementRepository {
				mockDisbursementRepo := new(mocks.IDisbursementRepository)

				return mockDisbursementRepo
			},
			mockCacheRepo: func() *mocks.ICacheRepository {
				mockCacheRepo := new(mocks.ICacheRepository)

				mockCacheRepo.On("Get", cacherepo.GetParam{
					Context: ctx,
					Key:     "ongoing_request_disbursement_with_user_id_1",
				}).Return("false", nil).Once()

				mockCacheRepo.On("Set", cacherepo.SetParam{
					Context: ctx,
					Key:     "ongoing_request_disbursement_with_user_id_1",
					Value:   "true",
					TTL:     5 * time.Minute,
				}).Return("OK", nil).Once()

				mockCacheRepo.On("Set", cacherepo.SetParam{
					Context: ctx,
					Key:     "ongoing_request_disbursement_with_user_id_1",
					Value:   "false",
					TTL:     5 * time.Minute,
				}).Return("OK", nil).Once()

				return mockCacheRepo
			},
			expectedErr: errors.New("error on get provider admin fee"),
		},
		{
			description: "when failed request disbursement on insufficient funds",
			param: disbursementservice.RequestDisbursementParam{
				Context:                 ctx,
				UserID:                  1,
				Amount:                  30000,
				DisbursementAccountGUID: disbursementAccountGUID,
			},
			mockWalletRepo: func() *mocks.IWalletRepository {
				mockWalletRepo := new(mocks.IWalletRepository)
				res := float64(20200)
				mockWalletRepo.On("GetWalletBalanceByUserID", walletrepo.GetWalletBalanceByUserIDParam{
					Context:  ctx,
					Executor: conn,
					UserID:   1,
				}).Return(res, nil).Once()

				return mockWalletRepo
			},
			mockPaymentProviderRepo: func() *mocks.IPaymentProviderRepository {
				mockPaymentProviderRepo := new(mocks.IPaymentProviderRepository)
				res := float64(200)
				mockPaymentProviderRepo.On("GetAdminFeeByGUID", paymentproviderrepo.GetAdminFeeByGUIDParam{
					Context:  ctx,
					Executor: conn,
					GUID:     paymentProviderGUID,
				}).Return(res, nil).Once()

				return mockPaymentProviderRepo
			},
			mockDisbursementAccountRepo: func() *mocks.IDisbursementAccountRepository {
				mockDisbursementAccountRepo := new(mocks.IDisbursementAccountRepository)
				mockDisbursementAccountRepo.On("GetByGUID", disbursementaccountrepo.GetByGUIDParam{
					Context:  ctx,
					Executor: conn,
					GUID:     disbursementAccountGUID,
				}).Return(&disbursementaccountrepo.GetByGUIDParamRes{
					PaymentProviderGUID: paymentProviderGUID,
					Name:                "Hansen",
					Number:              "123",
				}, nil).Once()

				return mockDisbursementAccountRepo
			},
			mockDisbursementRepo: func() *mocks.IDisbursementRepository {
				mockDisbursementRepo := new(mocks.IDisbursementRepository)

				return mockDisbursementRepo
			},
			mockCacheRepo: func() *mocks.ICacheRepository {
				mockCacheRepo := new(mocks.ICacheRepository)

				mockCacheRepo.On("Get", cacherepo.GetParam{
					Context: ctx,
					Key:     "ongoing_request_disbursement_with_user_id_1",
				}).Return("false", nil).Once()

				mockCacheRepo.On("Set", cacherepo.SetParam{
					Context: ctx,
					Key:     "ongoing_request_disbursement_with_user_id_1",
					Value:   "true",
					TTL:     5 * time.Minute,
				}).Return("OK", nil).Once()

				mockCacheRepo.On("Set", cacherepo.SetParam{
					Context: ctx,
					Key:     "ongoing_request_disbursement_with_user_id_1",
					Value:   "false",
					TTL:     5 * time.Minute,
				}).Return("OK", nil).Once()

				return mockCacheRepo
			},
			expectedErr: errors.New("insufficient funds. Please ensure your account balance is sufficient"),
		},
		{
			description: "when failed request disbursement on update wallet balance",
			param: disbursementservice.RequestDisbursementParam{
				Context:                 ctx,
				UserID:                  1,
				Amount:                  10000,
				DisbursementAccountGUID: disbursementAccountGUID,
			},
			mockWalletRepo: func() *mocks.IWalletRepository {
				mockWalletRepo := new(mocks.IWalletRepository)
				res := float64(20200)
				mockWalletRepo.On("GetWalletBalanceByUserID", walletrepo.GetWalletBalanceByUserIDParam{
					Context:  ctx,
					Executor: conn,
					UserID:   1,
				}).Return(res, nil).Once()

				mockWalletRepo.On("UpdateWalletBalanceByUserID", mock.AnythingOfType("UpdateWalletBalanceByUserIDParam")).Return(errors.New("error on update wallet balance")).Once()

				return mockWalletRepo
			},
			mockPaymentProviderRepo: func() *mocks.IPaymentProviderRepository {
				mockPaymentProviderRepo := new(mocks.IPaymentProviderRepository)
				res := float64(200)
				mockPaymentProviderRepo.On("GetAdminFeeByGUID", paymentproviderrepo.GetAdminFeeByGUIDParam{
					Context:  ctx,
					Executor: conn,
					GUID:     paymentProviderGUID,
				}).Return(res, nil).Once()

				return mockPaymentProviderRepo
			},
			mockDisbursementAccountRepo: func() *mocks.IDisbursementAccountRepository {
				mockDisbursementAccountRepo := new(mocks.IDisbursementAccountRepository)
				mockDisbursementAccountRepo.On("GetByGUID", disbursementaccountrepo.GetByGUIDParam{
					Context:  ctx,
					Executor: conn,
					GUID:     disbursementAccountGUID,
				}).Return(&disbursementaccountrepo.GetByGUIDParamRes{
					PaymentProviderGUID: paymentProviderGUID,
					Name:                "Hansen",
					Number:              "123",
				}, nil).Once()

				return mockDisbursementAccountRepo
			},
			mockDisbursementRepo: func() *mocks.IDisbursementRepository {
				mockDisbursementRepo := new(mocks.IDisbursementRepository)

				return mockDisbursementRepo
			},
			mockCacheRepo: func() *mocks.ICacheRepository {
				mockCacheRepo := new(mocks.ICacheRepository)

				mockCacheRepo.On("Get", cacherepo.GetParam{
					Context: ctx,
					Key:     "ongoing_request_disbursement_with_user_id_1",
				}).Return("false", nil).Once()

				mockCacheRepo.On("Set", cacherepo.SetParam{
					Context: ctx,
					Key:     "ongoing_request_disbursement_with_user_id_1",
					Value:   "true",
					TTL:     5 * time.Minute,
				}).Return("OK", nil).Once()

				mockCacheRepo.On("Set", cacherepo.SetParam{
					Context: ctx,
					Key:     "ongoing_request_disbursement_with_user_id_1",
					Value:   "false",
					TTL:     5 * time.Minute,
				}).Return("OK", nil).Once()

				return mockCacheRepo
			},
			expectedErr: errors.New("error on update wallet balance"),
		},
		{
			description: "when failed request disbursement on insert disbursement",
			param: disbursementservice.RequestDisbursementParam{
				Context:                 ctx,
				UserID:                  1,
				Amount:                  10000,
				DisbursementAccountGUID: disbursementAccountGUID,
			},
			mockWalletRepo: func() *mocks.IWalletRepository {
				mockWalletRepo := new(mocks.IWalletRepository)
				res := float64(20200)
				mockWalletRepo.On("GetWalletBalanceByUserID", walletrepo.GetWalletBalanceByUserIDParam{
					Context:  ctx,
					Executor: conn,
					UserID:   1,
				}).Return(res, nil).Once()

				mockWalletRepo.On("UpdateWalletBalanceByUserID", mock.AnythingOfType("UpdateWalletBalanceByUserIDParam")).Return(nil).Once()

				return mockWalletRepo
			},
			mockPaymentProviderRepo: func() *mocks.IPaymentProviderRepository {
				mockPaymentProviderRepo := new(mocks.IPaymentProviderRepository)
				res := float64(200)
				mockPaymentProviderRepo.On("GetAdminFeeByGUID", paymentproviderrepo.GetAdminFeeByGUIDParam{
					Context:  ctx,
					Executor: conn,
					GUID:     paymentProviderGUID,
				}).Return(res, nil).Once()

				return mockPaymentProviderRepo
			},
			mockDisbursementAccountRepo: func() *mocks.IDisbursementAccountRepository {
				mockDisbursementAccountRepo := new(mocks.IDisbursementAccountRepository)
				mockDisbursementAccountRepo.On("GetByGUID", disbursementaccountrepo.GetByGUIDParam{
					Context:  ctx,
					Executor: conn,
					GUID:     disbursementAccountGUID,
				}).Return(&disbursementaccountrepo.GetByGUIDParamRes{
					PaymentProviderGUID: paymentProviderGUID,
					Name:                "Hansen",
					Number:              "123",
				}, nil).Once()

				return mockDisbursementAccountRepo
			},
			mockDisbursementRepo: func() *mocks.IDisbursementRepository {
				mockDisbursementRepo := new(mocks.IDisbursementRepository)

				mockDisbursementRepo.On("Insert", mock.AnythingOfType("InsertDisbursementParam")).Return(errors.New("error on insert disbursement")).Once()

				return mockDisbursementRepo
			},
			mockCacheRepo: func() *mocks.ICacheRepository {
				mockCacheRepo := new(mocks.ICacheRepository)

				mockCacheRepo.On("Get", cacherepo.GetParam{
					Context: ctx,
					Key:     "ongoing_request_disbursement_with_user_id_1",
				}).Return("false", nil).Once()

				mockCacheRepo.On("Set", cacherepo.SetParam{
					Context: ctx,
					Key:     "ongoing_request_disbursement_with_user_id_1",
					Value:   "true",
					TTL:     5 * time.Minute,
				}).Return("OK", nil).Once()

				mockCacheRepo.On("Set", cacherepo.SetParam{
					Context: ctx,
					Key:     "ongoing_request_disbursement_with_user_id_1",
					Value:   "false",
					TTL:     5 * time.Minute,
				}).Return("OK", nil).Once()

				return mockCacheRepo
			},
			expectedErr: errors.New("error on insert disbursement"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			s := disbursementservice.NewDisbursementService(disbursementservice.Dependency{
				WalletRepository:              tc.mockWalletRepo(),
				PaymentProviderRepository:     tc.mockPaymentProviderRepo(),
				DisbursementAccountRepository: tc.mockDisbursementAccountRepo(),
				DisbursementRepository:        tc.mockDisbursementRepo(),
				CacheRepository:               tc.mockCacheRepo(),
			})

			err := s.RequestDisbursement(tc.param)
			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr.Error())
				return
			}

			assert.NoError(t, err)
			assert.Nil(t, err)
		})
	}
}
