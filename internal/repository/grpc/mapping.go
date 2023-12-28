package grpc

import (
	"HnH/pkg/authUtils"
	"HnH/pkg/serverErrors"
	notifOpts "HnH/services/notifications/pkg/searchOptions"
	notifErr "HnH/services/notifications/pkg/serviceErrors"
	searchErr "HnH/services/searchEngineService/pkg/searchOptions"

	"fmt"
)

var strToErr = map[string]error{
	authUtils.ERROR_WHILE_WRITING.Error():      authUtils.ERROR_WHILE_WRITING,
	authUtils.ERROR_WHILE_DELETING.Error():     authUtils.ERROR_WHILE_DELETING,
	serverErrors.NO_SESSION.Error():            serverErrors.NO_SESSION,
	serverErrors.INTERNAL_SERVER_ERROR.Error(): serverErrors.INTERNAL_SERVER_ERROR,

	serverErrors.ErrEntityNotFound.Error():    serverErrors.ErrEntityNotFound,
	serverErrors.ErrQuestionsNotFound.Error(): serverErrors.ErrQuestionsNotFound,
	serverErrors.ErrAnswerNotFound.Error():    serverErrors.ErrAnswerNotFound,
	serverErrors.ErrNoLastUpdate.Error():      serverErrors.ErrNoLastUpdate,

	notifOpts.ErrNoOption.Error():         notifOpts.ErrNoOption,
	notifOpts.ErrWrongValueFormat.Error(): notifOpts.ErrWrongValueFormat,

	notifErr.ErrOpenConn.Error():          notifErr.ErrOpenConn,
	notifErr.ErrInvalidUserID.Error():     notifErr.ErrInvalidUserID,
	notifErr.ErrConnAlreadyExists.Error(): notifErr.ErrConnAlreadyExists,
	notifErr.ErrNoConn.Error():            notifErr.ErrNoConn,
	notifErr.ErrInvalidConnection.Error(): notifErr.ErrInvalidConnection,

	searchErr.ErrNoOption.Error():         searchErr.ErrNoOption,
	searchErr.ErrWrongValueFormat.Error(): searchErr.ErrWrongValueFormat,
}

func GetErrByMessage(message string) error {
	errToReturn, ok := strToErr[message]
	if !ok {
		return fmt.Errorf(message)
	}

	return errToReturn
}
