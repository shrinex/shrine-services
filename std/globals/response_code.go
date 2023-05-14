package globals

type ResponseCode int32

const (
	SessionExpired  ResponseCode = 625
	SessionReplaced ResponseCode = 626
	SessionOverflow ResponseCode = 627
)
