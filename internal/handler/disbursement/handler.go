package disbursementhandler

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/hansengotama/disbursement/internal/lib/httphelper"
	cacherepo "github.com/hansengotama/disbursement/internal/repository/cache"
	disbursementrepo "github.com/hansengotama/disbursement/internal/repository/disbursement"
	"github.com/hansengotama/disbursement/internal/repository/disbursementaccount"
	"github.com/hansengotama/disbursement/internal/repository/paymentprovider"
	walletrepo "github.com/hansengotama/disbursement/internal/repository/wallet"
	disbursementservice "github.com/hansengotama/disbursement/internal/service/disbursement"
	"io/ioutil"
	"net/http"
)

type DisbursementRequestBody struct {
	Amount                     float64 `json:"amount"`
	DisbursementAccountStrGUID string  `json:"disbursement_account_guid"`
}

func HandleRequestDisbursement(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	if userID <= 0 {
		httphelper.Response(w, httphelper.HTTPResponse{
			Code:       http.StatusUnauthorized,
			ErrMessage: httphelper.ErrUnauthorized.Error(),
		})
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httphelper.Response(w, httphelper.HTTPResponse{
			Code:       http.StatusBadRequest,
			ErrMessage: httphelper.ErrReadingRequestBody.Error(),
		})
		return
	}

	var request DisbursementRequestBody
	err = json.Unmarshal(body, &request)
	if err != nil {
		httphelper.Response(w, httphelper.HTTPResponse{
			Code:       http.StatusBadRequest,
			ErrMessage: httphelper.ErrParsingRequestBody.Error(),
		})
		return
	}

	if request.Amount <= 0 {
		httphelper.Response(w, httphelper.HTTPResponse{
			Code:       http.StatusUnprocessableEntity,
			ErrMessage: "amount must be more than 0",
		})
		return
	}

	if request.DisbursementAccountStrGUID == "" {
		httphelper.Response(w, httphelper.HTTPResponse{
			Code:       http.StatusUnprocessableEntity,
			ErrMessage: "disbursement_account_guid is required",
		})
		return
	}

	disbursementAccountGUID, err := uuid.Parse(request.DisbursementAccountStrGUID)
	if err != nil {
		httphelper.Response(w, httphelper.HTTPResponse{
			Code:       http.StatusUnprocessableEntity,
			ErrMessage: "Invalid disbursement_account_guid. Please provide a valid UUID.",
		})
		return
	}

	dep := disbursementservice.Dependency{
		WalletRepository:              walletrepo.WalletDB{},
		PaymentProviderRepository:     paymentproviderrepo.PaymentProviderDB{},
		DisbursementAccountRepository: disbursementaccountrepo.DisbursementAccountDB{},
		DisbursementRepository:        disbursementrepo.DisbursementDB{},
		CacheRepository:               cacherepo.NewCacheRedis(),
	}
	s := disbursementservice.NewDisbursementService(dep)
	err = s.RequestDisbursement(disbursementservice.RequestDisbursementParam{
		Context:                 r.Context(),
		UserID:                  userID,
		Amount:                  request.Amount,
		DisbursementAccountGUID: disbursementAccountGUID,
	})
	if err != nil {
		httphelper.Response(w, httphelper.HTTPResponse{
			Code:       http.StatusInternalServerError,
			ErrMessage: err.Error(),
		})
		return
	}

	httphelper.Response(w, httphelper.HTTPResponse{
		Code: http.StatusOK,
	})
}
