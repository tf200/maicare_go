package leave

import "errors"

var (
	ErrLeaveRequestInvalidRequest = errors.New("invalid leave request")
	ErrLeaveRequestNotFound       = errors.New("leave request not found")
	ErrLeaveRequestStateInvalid   = errors.New("leave request is not in an editable state")
	ErrLeaveRequestForbidden      = errors.New("leave request is not accessible by the actor")
)
