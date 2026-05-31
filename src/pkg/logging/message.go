package logging

const (
	// Database
	MsgDBQueryError     = "db query error"
	MsgDBExecError      = "db exec error"
	MsgDBConnFailed     = "database connection failed"
	MsgDBRecordNotFound = "db record not found"

	// Cache
	MsgCacheMiss     = "cache miss"
	MsgCacheSetError = "cache set error"
	MsgCacheGetError = "cache get error"

	// HTTP
	MsgHTTPRequestFailed   = "http request failed"
	MsgHTTPInvalidResponse = "invalid http response"
	MsgExternalTimeout     = "external service timeout"

	// Validation
	MsgInvalidInput     = "invalid input"
	MsgUserNotFound     = "user not found"
	MsgPermissionDenied = "permission denied"
	MsgDuplicateEntry   = "duplicate entry"
	MsgUnauthorized     = "unauthorized access"

	// Basic
	MsgUnexpectedError   = "unexpected error"
	MsgOperationCanceled = "operation canceled"
	MsgPanicRecovered    = "panic recovered"
)
