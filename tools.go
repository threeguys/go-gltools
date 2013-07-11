package gltools

type GLToolsError struct {
	msg string
}

func (err GLToolsError) Error() string {
	return err.msg
}

