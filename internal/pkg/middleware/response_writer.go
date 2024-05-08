package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type response struct {
	ErrorCode    int         `json:"error_code"`
	ErrorMessage string      `json:"error_message"`
	Data         interface{} `json:"data"`
}

func SuccessResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	payload := response{
		Data: data,
	}

	out, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("cannot marshal success response payload")
		writeServerErrorResponse(w)
		return
	}

	_, err = w.Write(out)
	if err != nil {
		fmt.Println("cannot write json success response")
		writeServerErrorResponse(w)
		return
	}
}

func ErrorResponse(w http.ResponseWriter, httpStatus int, err error) {
	// Printing the error
	fmt.Println(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)

	payload := response{
		ErrorCode:    httpStatus,
		ErrorMessage: err.Error(),
	}

	out, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("error occurred while marshaling response payload")
		writeServerErrorResponse(w)
		return
	}

	_, err = w.Write(out)
	if err != nil {
		fmt.Println("error occurred while writing response")
		writeServerErrorResponse(w)
		return
	}
}

func writeServerErrorResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	_, err := w.Write([]byte(fmt.Sprintf("{\"message\":%s}", "internal server error")))
	if err != nil {
		fmt.Println("error occurred while writing response")
	}
}
