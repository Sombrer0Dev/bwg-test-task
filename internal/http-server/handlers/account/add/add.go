package add

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
	Currency string `json:"currency" validate:"required,currency"`
}

type Response struct {
	response.Response
	Error  string    `json:"error,omitempty"`
	Wallet uuid.UUID `json:"Wallet,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=Account
type Account interface {
	Add(currency string) (uuid.UUID, int64, error)
}

func New(log *slog.Logger, Account Account) http.HandlerFunc {
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

		wallet, id, err := Account.Add(req.Currency)
		if err != nil {
			log.Error("failed to add account", sl.Err(err))

			render.JSON(w, r, response.Error("failed to add account"))
		}

		log.Info("account added", slog.Int64("id", id))

		responseOK(w, r, wallet)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, wallet uuid.UUID) {
	render.JSON(w, r, Response{
		Response: response.OK(),
		Wallet:   wallet,
	})
}
