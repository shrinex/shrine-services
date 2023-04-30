package authx

import "time"

type (
	FlushMode int
)

const (
	FlushModeOnSave FlushMode = iota
	FlushModeImmediate
)

var (
	nowFunc = time.Now
)

const (
	tokenKey          = "__tokenKey"
	timeoutKey        = "__timeoutKey"
	platformKey       = "__platformKey"
	startTimeKey      = "__startTimeKey"
	idleTimeoutKey    = "__idleTimeoutKey"
	lastAccessTimeKey = "__lastAccessTimeKey"

	sessionKeyPrefix          = "shrine:session"
	sessionAttributeKeyPrefix = "shrine:session:attribute"
	sessionPrincipalKeyPrefix = "shrine:session:principal"
)
