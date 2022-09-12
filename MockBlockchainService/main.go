package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
)

const (
	SERVER_PORT = "5001"
)

func main() {
	http.HandleFunc("/wallets/create", mockCreateWallet)
	err := http.ListenAndServe(fmt.Sprintf(":%s", SERVER_PORT), nil)
	if err != nil {
		fmt.Print(err)
	}
}

/*
Mocks a blockchain service to create a wallet
*/
func mockCreateWallet(w http.ResponseWriter, r *http.Request) {
	// Check method
	checkHttpMethod("POST", w, r)

	// Read body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		sendHttpError(http.StatusInternalServerError, "Unable to read body", w, err)
		return
	}

	// Handle body closure
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("IO error, Unable to close http request body : %s", err)
		}
	}(r.Body)

	// Map body to struct
	mockServiceExpectation := MockServiceExpectation{}
	err = json.Unmarshal(body, &mockServiceExpectation)
	if err != nil {
		sendHttpError(http.StatusInternalServerError, "Unable to serialize body", w, err)
		return
	}

	// The service should do something here with the provided pincode and blockchain to create the wallet, once it's done, the new wallet address is sent back
	mockServiceAnswer := MockServiceAnswer{
		WalletAddress:   uuid.NewString(),
		CurrencyCode:    "ETH",
		CurrencyBalance: "0",
	}

	// Deserialize response to json
	mockServiceAnswerJson, err := json.Marshal(mockServiceAnswer)
	if err != nil {
		sendHttpError(http.StatusInternalServerError, "Unable to deserialize response", w, err)
		return
	}

	// Send back the response
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(mockServiceAnswerJson)
	if err != nil {
		log.Print("Unable to write response : %s", err)
	}
}
