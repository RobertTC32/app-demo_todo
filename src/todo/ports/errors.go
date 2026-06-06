package todo_ports

// Two types of errors exist in the ports:
// 1) error due general system problem and with no special handling,
// 2) error which is recognised and needs special handling (eg shown on screen)
//

// Sentinel errors:
// package level variables which indicate occurrence of unique events

// ERRORS for port_out-store

type StoreError struct {
	Message string
}

func NewStoreError(message string) StoreError {
	return StoreError{Message: message}
}

// implements "error" interface
func (this StoreError) Error() string {
	return this.Message
}

var (
	ErrNotFoundInStore           = NewStoreError("Not found")
	ErrMaxRowsExceededInStore    = NewStoreError("Max number exceeded")
	ErrAlreadyExistsInStore      = NewStoreError("Already exists")
	ErrConstraintViolatedInStore = NewStoreError("Constraint violation")
	ErrChangedByUserInStore      = NewStoreError("Changed by other user")
)

// ERRORS for port_in-ui and port_in-api

type AppError struct {
	Message string
}

func NewAppError(message string) AppError {
	return AppError{Message: message}
}

// implements "error" interface
func (this AppError) Error() string {
	return this.Message
}

var (
	ErrDummyInApp = NewAppError("Dummy failure")
)
