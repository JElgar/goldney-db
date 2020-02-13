package errors

type ApiError struct {
    // This is the error thown by go
   Err error
    // This is the human readable error
   Message string
    // This is the error code
   Code int
}

func (e ApiError) Error() string {
   return e.Message
}
