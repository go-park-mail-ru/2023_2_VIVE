package test

import (
	"HnH/internal/requestHandlers"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

var LoginCorrectCases = []JsonTestCase{
	{
		requestBody:  `{"email":"vive@mail.ru", "password":"Vive2023top~", "role":"applicant"}`,
		statusCode:   http.StatusOK,
		responseBody: "",
	},
	{
		requestBody:  `{"email":"vk_ed@mail.ru", "password":"Technopark2023!", "role":"applicant"}`,
		statusCode:   http.StatusOK,
		responseBody: "",
	},
	{
		requestBody:  `{"email":"petr98@mail.ru", "password":"PetyaMolodec23!", "role":"employer"}`,
		statusCode:   http.StatusOK,
		responseBody: "",
	},
	{
		requestBody:  `{"email":"golang@gmail.com", "password":"GolangEnjoyer2002?", "role":"applicant"}`,
		statusCode:   http.StatusOK,
		responseBody: "",
	},
	{
		requestBody:  `{"email":"katya1729@empire.ru", "password":"TheEmpress29#", "role":"employer"}`,
		statusCode:   http.StatusOK,
		responseBody: "",
	},
}

var TestToCookie = sync.Map{}

func TestLoginLogoutCorrectInput(t *testing.T) {
	for num, testItem := range LoginCorrectCases {
		req := httptest.NewRequest("POST", sessionUrl, strings.NewReader(testItem.requestBody))
		w := httptest.NewRecorder()

		requestHandlers.Login(w, req)
		require.Equal(t, w.Code, testItem.statusCode)

		resp := w.Result()
		cookie := resp.Cookies()
		require.Equal(t, len(cookie), 1)

		TestToCookie.Store(num, cookie[0])

		body, _ := io.ReadAll(resp.Body)
		require.Equal(t, string(body), testItem.responseBody)
	}

	TestToCookie.Range(func(testNum, cookie interface{}) bool {
		req := httptest.NewRequest("GET", sessionUrl, nil)

		sessionCookie, _ := cookie.(*http.Cookie)
		req.AddCookie(sessionCookie)
		w := httptest.NewRecorder()

		requestHandlers.CheckLogin(w, req)
		require.Equal(t, w.Code, http.StatusOK)

		return true
	})

	TestToCookie.Range(func(testNum, cookie interface{}) bool {
		req := httptest.NewRequest("DELETE", sessionUrl, nil)

		sessionCookie, _ := cookie.(*http.Cookie)
		req.AddCookie(sessionCookie)
		w := httptest.NewRecorder()

		requestHandlers.Logout(w, req)
		require.Equal(t, w.Code, http.StatusOK)

		return true
	})
}

var LoginIncorrectCases = []JsonTestCase{
	{
		requestBody:  `{"email":"unknown@mail.ru", "password":"vive2023top", "role":"applicant"}`,
		statusCode:   http.StatusUnauthorized,
		responseBody: NO_DATA_FOUND,
	},
	{
		requestBody:  `{"email":"vk_ed@mail.ru", "password":"itIsNotTheTechnoparkAnymore", "role":"employer"}`,
		statusCode:   http.StatusUnauthorized,
		responseBody: INCORRECT_CREDENTIALS,
	},
	{
		requestBody:  `{"email":"", "password":"petyamolodec", "role":"employer"}`,
		statusCode:   http.StatusUnauthorized,
		responseBody: INCORRECT_CREDENTIALS,
	},
	{
		requestBody:  `{"email":"golang@gmail.com", "password":"", "role":"employer"}`,
		statusCode:   http.StatusUnauthorized,
		responseBody: INCORRECT_CREDENTIALS,
	},
	{
		requestBody:  `{:"katya1729@empire.ru", "password":"theempress", "role":"employer"}`,
		statusCode:   http.StatusBadRequest,
		responseBody: MISSED_FIELD_JSON,
	},
	{
		requestBody:  `{"email":"golang@gmail.com", "password":"GolangEnjoyer2002?", "role":"employer"}`,
		statusCode:   http.StatusUnauthorized,
		responseBody: INCORRECT_ROLE,
	},
}

func TestLoginIncorrectInput(t *testing.T) {
	for _, testItem := range LoginIncorrectCases {
		req := httptest.NewRequest("POST", sessionUrl, strings.NewReader(testItem.requestBody))
		w := httptest.NewRecorder()

		requestHandlers.Login(w, req)
		require.Equal(t, w.Code, testItem.statusCode)

		resp := w.Result()
		cookie := resp.Cookies()
		require.Equal(t, len(cookie), 0)

		body, _ := io.ReadAll(resp.Body)
		require.Equal(t, string(body), testItem.responseBody)
	}
}

var CheckLoginIncorrectCases = []CookieTestCase{
	{
		cookie: &http.Cookie{
			Name:  "session",
			Value: "123456abc",
		},
		expectedStatusCode: http.StatusUnauthorized,
		expectedError:      INVALID_COOKIE,
	},
	{
		cookie: &http.Cookie{
			Name:  "incorrectName",
			Value: "12s3f-fa541-34gdf",
		},
		expectedStatusCode: http.StatusUnauthorized,
		expectedError:      NO_COOKIE,
	},
	{
		cookie:             nil,
		expectedStatusCode: http.StatusUnauthorized,
		expectedError:      NO_COOKIE,
	},
}

func TestCheckLoginIncorrectInput(t *testing.T) {
	for _, testItem := range CheckLoginIncorrectCases {
		req := httptest.NewRequest("GET", sessionUrl, nil)

		sessionCookie := testItem.cookie

		if sessionCookie != nil {
			req.AddCookie(sessionCookie)
		}

		w := httptest.NewRecorder()

		requestHandlers.CheckLogin(w, req)
		require.Equal(t, w.Code, testItem.expectedStatusCode)

		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)
		require.Equal(t, string(body), testItem.expectedError)
	}
}

var LogoutIncorrectCases = []CookieTestCase{
	{
		cookie: &http.Cookie{
			Name:  "session",
			Value: "82s7e-4a54o-g48n7",
		},
		expectedStatusCode: http.StatusUnauthorized,
		expectedError:      AUTH_REQUIRED,
	},
	{
		cookie: &http.Cookie{
			Name:  "incorrectName",
			Value: "12s3f-fa541-34gdf",
		},
		expectedStatusCode: http.StatusUnauthorized,
		expectedError:      NO_COOKIE,
	},
	{
		cookie:             nil,
		expectedStatusCode: http.StatusUnauthorized,
		expectedError:      NO_COOKIE,
	},
}

func TestLogoutIncorrectInput(t *testing.T) {
	for _, testItem := range LogoutIncorrectCases {
		req := httptest.NewRequest("DELETE", sessionUrl, nil)

		sessionCookie := testItem.cookie

		if sessionCookie != nil {
			req.AddCookie(sessionCookie)
		}

		w := httptest.NewRecorder()

		requestHandlers.Logout(w, req)
		require.Equal(t, w.Code, testItem.expectedStatusCode)

		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)
		require.Equal(t, string(body), testItem.expectedError)
	}
}
