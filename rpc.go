package couchcampaign

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, err error) {
	code := Code(err)
	if code == http.StatusInternalServerError {
		log.Println(err)
		Respond(w, code, "")
		return
	}
	Respond(w, code, err.Error())
}

func Respond(w http.ResponseWriter, code int, format string, args ...interface{}) {
	w.WriteHeader(code)
	fmt.Fprintf(w, format, args...)
}

func Error(code int) error {
	return Errorf(code, "")
}

func Errorf(code int, format string, args ...interface{}) error {
	e := statusCodeError(code)
	m := fmt.Sprintf(format, args...)
	return fmt.Errorf("%w: %s", e, m)
}

type statusCodeError int

func (e statusCodeError) Error() string {
	return http.StatusText(int(e))
}

func Code(err error) int {
	var e statusCodeError
	if errors.As(err, &e) {
		return int(e)
	}
	return http.StatusInternalServerError
}
