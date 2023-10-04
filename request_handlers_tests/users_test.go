package requestHandlers_test

import (
	"HnH/requestHandlers"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var SignupCorrectCases = []JsonTestCase{
	{
		requestBody:  `{"email":"vive@mail.ru", "password":"Vive2023top~", "first_name":"Vladimir", "last_name":"Borozenets", "role":"applicant"}`,
		statusCode:   http.StatusOK,
		responseBody: "",
	},
	{
		requestBody:  `{"email":"vk_ed@mail.ru", "password":"Technopark2023!", "first_name":"Ivan", "last_name":"Pupkin", "role":"applicant"}`,
		statusCode:   http.StatusOK,
		responseBody: "",
	},
	{
		requestBody:  `{"email":"petr98@mail.ru", "password":"PetyaMolodec23!", "first_name":"Petr", "last_name":"Ivanov", "role":"employer"}`,
		statusCode:   http.StatusOK,
		responseBody: "",
	},
	{
		requestBody:  `{"email":"golang@gmail.com", "password":"GolangEnjoyer2002?", "first_name":"Sergey", "last_name":"Alekseev", "role":"applicant"}`,
		statusCode:   http.StatusOK,
		responseBody: "",
	},
	{
		requestBody:  `{"email":"katya1729@empire.ru", "password":"TheEmpress29#", "first_name":"Ekaterina", "last_name":"Vtoraya", "role":"employer"}`,
		statusCode:   http.StatusOK,
		responseBody: "",
	},
}

func TestMain(m *testing.M) {
	for _, testItem := range SignupCorrectCases {
		req := httptest.NewRequest("POST", usersUrl, strings.NewReader(testItem.requestBody))
		w := httptest.NewRecorder()

		requestHandlers.SignUp(w, req)
	}

	os.Exit(m.Run())
}

var NewSignupCorrectCases = []JsonTestCase{
	{
		requestBody:  `{"email":"fn11_73b@bmstu.ru", "password":"Moskva^22", "first_name":"Ruslan", "last_name":"Novikov", "role":"applicant"}`,
		statusCode:   http.StatusOK,
		responseBody: "",
	},
	{
		requestBody:  `{"email":"government@nsk.ru", "password":"NSK154_top%", "first_name":"Anatoliy", "last_name":"Lokot", "role":"employer"}`,
		statusCode:   http.StatusOK,
		responseBody: "",
	},
}

func TestSignupCorrectInput(t *testing.T) {
	for _, testItem := range NewSignupCorrectCases {
		req := httptest.NewRequest("POST", usersUrl, strings.NewReader(testItem.requestBody))
		w := httptest.NewRecorder()

		requestHandlers.SignUp(w, req)
		require.Equal(t, w.Code, testItem.statusCode)

		resp := w.Result()
		cookie := resp.Cookies()
		require.Equal(t, len(cookie), 1)

		body, _ := io.ReadAll(resp.Body)
		require.Equal(t, string(body), testItem.responseBody)
	}
}

var SignupIncorrectCases = []JsonTestCase{
	{
		requestBody:  `{"email":"vive@mail.ru", "password":"New_password2", "first_name":"Vladimir", "last_name":"Borozenets", "role":"applicant"}`,
		statusCode:   http.StatusUnauthorized,
		responseBody: ACCOUNT_ALREADY_EXISTS,
	},
	{
		requestBody:  `{"email":"someMail@mail.ru", "password":"", "first_name":"Pichai", "last_name":"Sundararajan", "role":"employer"}`,
		statusCode:   http.StatusUnauthorized,
		responseBody: INCORRECT_CREDENTIALS,
	},
	{
		requestBody:  `{"email":"specialist@mail.ru", "password":"masterOfScience@man11", "first_name":"Anton", "last_name":"Umnov", "role":"genius"}`,
		statusCode:   http.StatusUnauthorized,
		responseBody: INVALID_ROLE,
	},
	{
		requestBody:  `{:"tech@gmail.com", "password":"hiTech_69", "first_name":"Pasha", "last_name":"Technik", "role":"applicant"}`,
		statusCode:   http.StatusBadRequest,
		responseBody: MISSED_FIELD_JSON,
	},
	{
		requestBody:  `{"email":"vk.com", "password":"!mailRuGroup2011", "first_name":"Alexey", "last_name":"Loginov", "role":"employer"}`,
		statusCode:   http.StatusUnauthorized,
		responseBody: INVALID_EMAIL,
	},
	{
		requestBody:  `{"email":"andrey02@mail.ru", "password":"Znatok_2002", "first_name":"Andrey", "last_name":"Dolgikh", "role":""}`,
		statusCode:   http.StatusUnauthorized,
		responseBody: INVALID_ROLE,
	},
	{
		requestBody:  `{"email":"hardwork@mail.ru", "password":"Trudogolik74", "first_name":"Nikolay", "last_name":"Belaz", "role":"employer"}`,
		statusCode:   http.StatusUnauthorized,
		responseBody: INVALID_PASSWORD,
	},
	{
		requestBody:  `{"email":"novichok@gmail.com", "password":"short", "first_name":"Evgeniy", "last_name":"Novikov", "role":"applicant"}`,
		statusCode:   http.StatusUnauthorized,
		responseBody: INVALID_PASSWORD,
	},
}

func TestSignupIncorrectInput(t *testing.T) {
	for _, testItem := range SignupIncorrectCases {
		req := httptest.NewRequest("POST", usersUrl, strings.NewReader(testItem.requestBody))
		w := httptest.NewRecorder()

		requestHandlers.SignUp(w, req)
		require.Equal(t, w.Code, testItem.statusCode)

		resp := w.Result()
		cookie := resp.Cookies()
		require.Equal(t, len(cookie), 0)

		body, _ := io.ReadAll(resp.Body)
		require.Equal(t, string(body), testItem.responseBody)
	}
}

var GetInfoCorrectCases = []GetUserTestCase{
	{
		authData:           `{"email":"vive@mail.ru", "password":"Vive2023top~"}`,
		expectedMessage:    `{"id":1,"email":"vive@mail.ru","first_name":"Vladimir","last_name":"Borozenets","role":"applicant"}`,
		expectedStatusCode: http.StatusOK,
	},
	{
		authData:           `{"email":"petr98@mail.ru", "password":"PetyaMolodec23!"}`,
		expectedMessage:    `{"id":3,"email":"petr98@mail.ru","first_name":"Petr","last_name":"Ivanov","role":"employer"}`,
		expectedStatusCode: http.StatusOK,
	},
	{
		authData:           `{"email":"katya1729@empire.ru", "password":"TheEmpress29#"}`,
		expectedMessage:    `{"id":5,"email":"katya1729@empire.ru","first_name":"Ekaterina","last_name":"Vtoraya","role":"employer"}`,
		expectedStatusCode: http.StatusOK,
	},
}

func TestGetInfoCorrectInput(t *testing.T) {
	for _, testItem := range GetInfoCorrectCases {
		authReq := httptest.NewRequest("POST", sessionUrl, strings.NewReader(testItem.authData))
		authWriter := httptest.NewRecorder()

		requestHandlers.Login(authWriter, authReq)

		authResp := authWriter.Result()
		cookie := authResp.Cookies()[0]

		req := httptest.NewRequest("GET", currentUserUrl, nil)
		req.AddCookie(cookie)
		w := httptest.NewRecorder()

		requestHandlers.GetInfo(w, req)
		require.Equal(t, w.Code, testItem.expectedStatusCode)

		resp := w.Result()

		body, _ := io.ReadAll(resp.Body)
		require.Equal(t, string(body), testItem.expectedMessage)
	}
}

var GetInfoIncorrectCases = []CookieTestCase{
	{
		cookie: &http.Cookie{
			Name:  "session",
			Value: "df56g-f5hg4-gd5h4",
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

func TestGetInfoIncorrectInput(t *testing.T) {
	for _, testItem := range GetInfoIncorrectCases {
		req := httptest.NewRequest("GET", currentUserUrl, nil)

		sessionCookie := testItem.cookie

		if sessionCookie != nil {
			req.AddCookie(sessionCookie)
		}

		w := httptest.NewRecorder()

		requestHandlers.GetInfo(w, req)
		require.Equal(t, w.Code, testItem.expectedStatusCode)

		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)
		require.Equal(t, string(body), testItem.expectedError)
	}
}
