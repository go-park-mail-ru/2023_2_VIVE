package appErrors

import (
	"HnH/internal/repository/psql"
	"HnH/internal/usecase"
	"HnH/pkg/authUtils"
	"HnH/pkg/serverErrors"
	"errors"

	"net/http"
)

var errToCode = map[error]int{
	serverErrors.INCORRECT_CREDENTIALS:  http.StatusBadRequest,
	serverErrors.INVALID_EMAIL:          http.StatusBadRequest,
	serverErrors.INVALID_PASSWORD:       http.StatusBadRequest,
	serverErrors.INVALID_ROLE:           http.StatusBadRequest,
	serverErrors.INCORRECT_ROLE:         http.StatusNotFound,
	serverErrors.NO_ACCOUNT_DATA_FOUND:  http.StatusNotFound,
	serverErrors.NO_DATA_FOUND:          http.StatusNotFound,
	serverErrors.ACCOUNT_ALREADY_EXISTS: http.StatusConflict,
	serverErrors.SESSION_ALREADY_EXISTS: http.StatusConflict,
	serverErrors.INVALID_COOKIE:         http.StatusUnauthorized,
	serverErrors.NO_SESSION:             http.StatusUnauthorized,
	serverErrors.NO_COOKIE:              http.StatusBadRequest,
	serverErrors.AUTH_REQUIRED:          http.StatusUnauthorized,
	serverErrors.FORBIDDEN:              http.StatusForbidden,
	serverErrors.INTERNAL_SERVER_ERROR:  http.StatusInternalServerError,
	serverErrors.ErrNoLastUpdate:        http.StatusNotFound,
	serverErrors.ErrEntityNotFound:      http.StatusNotFound,
	serverErrors.ErrQuestionsNotFound:   http.StatusNotFound,
	serverErrors.ErrAnswerNotFound:      http.StatusNotFound,

	authUtils.INCORRECT_CREDENTIALS: http.StatusBadRequest,
	authUtils.INVALID_EMAIL:         http.StatusBadRequest,
	authUtils.EMPTY_EMAIL:           http.StatusBadRequest,
	authUtils.EMPTY_PASSWORD:        http.StatusBadRequest,
	authUtils.INVALID_PASSWORD:      http.StatusBadRequest,
	authUtils.ENTITY_NOT_FOUND:      http.StatusNotFound,
	authUtils.ERROR_WHILE_WRITING:   http.StatusNotModified,
	authUtils.ERROR_WHILE_DELETING:  http.StatusNotModified,

	usecase.ErrInapropriateRole: http.StatusForbidden,
	usecase.ErrForbidden:        http.StatusForbidden,
	usecase.ErrReadAvatar:       http.StatusInternalServerError,
	usecase.BadAvatarSize:       http.StatusBadRequest,
	usecase.BadAvatarType:       http.StatusBadRequest,

	psql.ErrEntityNotFound:     http.StatusNotFound,
	psql.ErrNotInserted:        http.StatusNotModified,
	psql.ErrNoRowsUpdated:      http.StatusNotModified,
	psql.ErrNoRowsDeleted:      http.StatusNotModified,
	psql.IncorrectUserID:       http.StatusBadRequest,
	psql.ErrRecordAlredyExists: http.StatusConflict,
}

func GetErrAndCodeToSend(err error) (error, int) {
	var source error
	for err != nil {
		source = err
		err = errors.Unwrap(err)
	}

	code, ok := errToCode[source]
	if !ok {
		return serverErrors.INTERNAL_SERVER_ERROR, http.StatusInternalServerError
	}

	return source, code
}
