package requestHandlers_test

import (
	"HnH/requestHandlers"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

var GetVacanciesCorrectCases = []JsonTestCase{
	{
		requestBody:  "",
		statusCode:   http.StatusOK,
		responseBody: `[{"id":1,"name":"C++ developer","company_name":"VK","description":"Middle C++ developer in Mail.ru team","salary":250000},{"id":2,"name":"Go developer","company_name":"VK","description":"Golang junior developer without any experience","salary":100000},{"id":3,"name":"HR","company_name":"Yandex","description":"Human resources specialist","salary":70000},{"id":4,"name":"Frontend developer","company_name":"Google","description":"Middle Frontend developer, JavaScript, HTML, Figma","salary":500000},{"id":5,"name":"Project Manager","company_name":"VK","description":"Experienced specialist in IT-management","salary":200000}]`,
	},
}

func TestGetVacanciesCorrectInput(t *testing.T) {
	for _, testItem := range GetVacanciesCorrectCases {
		req := httptest.NewRequest("GET", vacanciesUrl, nil)
		w := httptest.NewRecorder()

		requestHandlers.GetVacancies(w, req)
		require.Equal(t, w.Code, testItem.statusCode)

		resp := w.Result()

		body, _ := io.ReadAll(resp.Body)
		require.Equal(t, string(body), testItem.responseBody)
	}
}
