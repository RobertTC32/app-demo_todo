package commons

import "fmt"

// Error Wrapping could be used for this:
// "errors.New()" can not wrap errors;
// how to wrap and unwrap errors?
// - Wrapping errors:
//   Wrap errors to add additional information;
//   To wrap another error,
//   the wrapping error must implement "Unwrap() error" or "Unwrap() []error" method;
//   Errors can be wrapped using "fmt.Errorf()" (with %w to wrap errors),
//   or "errors.Join()" (to wrap multiple provided errors, discarding nil errors);
//   Errorf implements the "Unwrap() error", and Join the "Unwrap() []error" method;
// - Unwrapping errors:
//   "Unwrap()" can be used to unwrap errors,
//   but "errors.Is()" and "errors.As()" is recommended;
//   "Unwrap" calls the "Unwrap" method implemented on the error,
//   and errors created by "fmt.Errorf()" implement the "Unwrap" method;
//   Errors created by "errors.Join()" do not implement the "Unwrap" method;
//
// How to use in code?
// - Inside called function:
//   logic return error err1;
//   return fmt.Errorf("my-function-as-err2: %w", err1)
// - Logic outside calling function:
//   function call returns error errF;
//   errors.Unwrap(errF) will return err1;
//   how to get first part of err2?
//

type WrappingError struct {
	Current error
	Wrapped error
}

func NewWrappingError(currentErr error, wrappedErr error) WrappingError {
	return WrappingError{Current: currentErr, Wrapped: wrappedErr}
}

// implement "error" interface
func (this WrappingError) Error() string {
	return this.Current.Error()
}

// implement "unwrap" interface
func (this WrappingError) Unwrap() error {
	return this.Wrapped
}

func (this WrappingError) CurrentErr() error {
	return this.Current
}

func (this WrappingError) FullError() string {
	return fmt.Sprintf("%v: %v", this.Current, this.Wrapped)
}
