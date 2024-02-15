package add_test

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
	"github.com/Sombrer0Dev/bwg-test-task/internal/http-server/handlers/account/add/mocks"
	"github.com/Sombrer0Dev/bwg-test-task/internal/utils/logger/handlers/slogdiscard"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name      string
		currency  string
		respError string
		mockError error
	}{
		{
			name:     "Success",
			currency: "USD",
		},
		{
			name:      "Empty Currency",
			currency:  "",
			respError: "field Currency is a required field",
		},
		{
			name:      "Invalid Currency",
			currency:  "some invalid currency",
			respError: "field Currency is not valid",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc := tc
			t.Parallel()

			AccountMock := mocks.NewAccount(t)

			if tc.respError == "" || tc.mockError != nil {
				AccountMock.On("Add", tc.currency, mock.AnythingOfType("string")).Return(uuid.New(), int64(1), tc.mockError).Once()
			}

			handler := add.New(slogdiscard.NewDiscardLogger(), AccountMock)
			input := fmt.Sprintf(`{"currency": "%s"}`, tc.currency)

			req, err := http.NewRequest(http.MethodPost, "/account/add", bytes.NewReader([]byte(input)))
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
