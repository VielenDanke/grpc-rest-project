package u

func (x *Error) Error() string {
	return x.Msg
}

func (x *ErrorNotFound) Error() string {
	return x.Msg
}
