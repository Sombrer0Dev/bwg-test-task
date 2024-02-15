package invoice

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/Sombrer0Dev/bwg-test-task/internal/utils/api/response"
	"github.com/Sombrer0Dev/bwg-test-task/internal/utils/logger/sl"
)

type Request struct {
	Currency string    `json:"currency" validate:"required,iso4217"`
	Wallet   uuid.UUID `json:"wallet"  validate:"required,uuid"`
	Amount   float64   `json:"amount" validate:"required,gte=0"`
}

type Response struct {
	response.Response
	Error   string  `json:"error,omitempty"`
	Balance float64 `json:"wallet,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=AccountInvoicer
type AccountInvoicer interface {
	Invoice(currency string, wallet uuid.UUID, amount float64) error
}

func New(log *slog.Logger, Account AccountInvoicer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, response.Error("Failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err = validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, response.ValidationError(validateErr))

			return
		}

		err = Account.Invoice(req.Currency, req.Wallet, req.Amount)
		if err != nil {
			log.Error("failed to get invoice", sl.Err(err))

			render.JSON(w, r, response.Error("failed to get invoice"))
		}

		log.Info("invoice started")

		responseOK(w, r)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, Response{
		Response: response.OK(),
	})
}
