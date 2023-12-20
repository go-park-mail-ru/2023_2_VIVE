package grpc

import (
	"HnH/pkg/authUtils"
	"HnH/pkg/serverErrors"
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
}

func GetErrByMessage(message string) error {
	errToReturn, ok := strToErr[message]
	if !ok {
		return serverErrors.INTERNAL_SERVER_ERROR
	}

	return errToReturn
}
