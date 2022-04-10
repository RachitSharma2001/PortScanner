package error

func NoError(err error) bool {
	return err == nil
}
