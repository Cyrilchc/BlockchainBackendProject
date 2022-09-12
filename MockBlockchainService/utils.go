package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func sendHttpError(httpErrorCode int, message string, w http.ResponseWriter, err error) {
	log.Print(err)
	w.WriteHeader(httpErrorCode)
	_, err = w.Write(buildHttpErrorMessage(message))
	if err != nil {
		log.Printf("%s : %s", message, err)
	}
}

func checkHttpMethod(method string, w http.ResponseWriter, r *http.Request) {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, err := w.Write(buildHttpErrorMessage("Only POST method is allowed"))
		if err != nil {
			log.Printf("Unable to write to http response : %s", err)
		}

		return
	}
}

func buildHttpErrorMessage(message string) []byte {
	resp := make(map[string]string)
	resp["error"] = message
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Print(err)
	}

	return jsonResp
}
