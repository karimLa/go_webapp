package controllers

type parseError string

func (e parseError) Error() string {
	return string(e)
}

func (e parseError) Public() string {
	return string(e)
}
