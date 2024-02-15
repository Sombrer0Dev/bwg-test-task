package invoice_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/Sombrer0Dev/bwg-test-task/internal/http-server/handlers/account/add"
	"github.com/Sombrer0Dev/bwg-test-task/internal/http-server/handlers/account/invoice"
	"github.com/Sombrer0Dev/bwg-test-task/internal/http-server/handlers/account/invoice/mocks"
	"github.com/Sombrer0Dev/bwg-test-task/internal/utils/logger/handlers/slogdiscard"
)

func TestInvoiceHandler(t *testing.T) {
	cases := []struct {
		name      string
		currency  string
		amount    float64
		wallet    uuid.UUID
		respError string
		mockError error
	}{
		{
			name:     "Success",
			currency: "USD",
			wallet:   uuid.New(),
			amount:   1.1,
		},
		{
			name:      "Invalid Currency",
			currency:  "some invalid currency",
			wallet:    uuid.New(),
			amount:    1.1,
			respError: "field Currency is not valid",
		},
		{
			name:      "Invalid Amount",
			currency:  "USD",
			wallet:    uuid.New(),
			amount:    -10,
			respError: "field Amount is not valid",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc := tc
			t.Parallel()

			AccountMock := mocks.NewAccountInvoicer(t)

			if tc.respError == "" || tc.mockError != nil {
				AccountMock.On("Invoice", tc.currency, tc.wallet, tc.amount, mock.AnythingOfType("string")).Return(tc.mockError).Once()
			}

			handler := invoice.New(slogdiscard.NewDiscardLogger(), AccountMock)
			input := fmt.Sprintf(`{"currency": "%s", "wallet": "%s", "amount": %f}`, tc.currency, tc.wallet, tc.amount)

			req, err := http.NewRequest(http.MethodPost, "/add", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp add.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
