package statuses

type Status int

const (
	OK                    Status = 200
	CREATED               Status = 201
	INVALID_REQUEST       Status = 400
	UNAUTHORIZED          Status = 401
	FORBIDDEN             Status = 403
	NOT_FOUND             Status = 404
	CONFLICT              Status = 409
	INTERNAL_SERVER_ERROR Status = 500
)
