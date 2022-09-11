package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
)

const (
	SERVER_PORT = "5000"
	HOST        = "localhost"
	PORT        = 5432
	USER        = "postgres"
	PASSWORD    = "cyril"
	DBNAME      = "blockchainbackend"
)

func main() {
	http.HandleFunc("/", createUserAndWallet)
	http.HandleFunc("/mock/wallets/create", mockCreateWallet)
	err := http.ListenAndServe(fmt.Sprintf(":%s", SERVER_PORT), nil)
	if err != nil {
		log.Fatal(err)
	}
}

/*
*
Creates a user and a wallet
*/
func createUserAndWallet(w http.ResponseWriter, r *http.Request) {
	db, err := connectDatabase(HOST, PORT, USER, PASSWORD, DBNAME)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal("Unable to close database connection")
		}
	}(db)

	regexUsername := regexp.MustCompile(`[a-z0-9_]{3,100}`)
	regexPassword := regexp.MustCompile(`.{6,32}`)
	regexPinCode := regexp.MustCompile(`\d{6}`)

	// Check method
	checkHttpMethod("POST", w, r)

	// Read body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write(buildHttpErrorMessage("Unable to read body"))
		if err != nil {
			log.Fatalf("Unable to read body : %s", err)
		}

		return
	}

	// Handle body closure
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("IO error, Unable to close http request body : %s", err)
		}
	}(r.Body)

	// Map body to struct
	player := Player{}
	err = json.Unmarshal(body, &player)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write(buildHttpErrorMessage("Unable to serialize body"))
		if err != nil {
			log.Fatalf("Unable to serialize body : %s", err)
		}

		return
	}

	// Uniqueness of username will be checked by database constraints

	// Username naming policy
	if !regexUsername.MatchString(player.Username) {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write(buildHttpErrorMessage("Username does not respect the naming policy"))
		if err != nil {
			log.Fatalf("Unable to write response : %s", err)
		}

		return
	}

	// Password policy
	if !regexPassword.MatchString(player.Password) {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write(buildHttpErrorMessage("Password does not respect the complexity policy"))
		if err != nil {
			log.Fatalf("Unable to write response : %s", err)
		}

		return
	}

	// Pincode policy
	if !regexPinCode.MatchString(player.Pincode) {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write(buildHttpErrorMessage("Pincode must contains 6 digits"))
		if err != nil {
			log.Fatalf("Unable to write response : %s", err)
		}

		return
	}

	// Prepare external service call
	MockServiceExpectation := MockServiceExpectation{
		Pincode:    player.Pincode,
		Blockchain: "ETH",
	}

	data, err := json.Marshal(MockServiceExpectation)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write(buildHttpErrorMessage("Unable to serialize MockServiceExpectation"))
		if err != nil {
			log.Fatalf("Unable to serialize MockServiceExpectation : %s", err)
		}

		return
	}

	// Send request to external service
	response, err := http.Post(fmt.Sprintf("http://%s:%s/mock/wallets/create", HOST, SERVER_PORT), "application/json", bytes.NewBuffer(data))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write(buildHttpErrorMessage("External service returned an error"))
		if err != nil {
			log.Fatalf("External service returned an error : %s", err)
		}

		return
	}

	// Process response
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write(buildHttpErrorMessage("Unable to read body"))
		if err != nil {
			log.Fatalf("Unable to read body : %s", err)
		}

		return
	}

	// Handle body closure
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("IO error, Unable to close http request body : %s", err)
		}
	}(response.Body)

	// Map body to struct
	mockServiceAnswer := MockServiceAnswer{}
	err = json.Unmarshal(responseBody, &mockServiceAnswer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write(buildHttpErrorMessage("Unable to serialize body"))
		if err != nil {
			log.Fatalf("Unable to serialize body : %s", err)
		}

		return
	}

	// Create wallet object
	wallet := Wallet{
		Address: mockServiceAnswer.WalletAddress,
	}

	// Insert player in database
	result, err := insertPlayer(&player, db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write(buildHttpErrorMessage("Unable to insert player"))
		if err != nil {
			log.Fatalf("Unable to insert player : %s", err)
		}

		return
	}

	// Insert Wallet in database
	playerId, err := result.LastInsertId()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write(buildHttpErrorMessage("Unable to get player's id"))
		if err != nil {
			log.Fatalf("Unable to get player's id : %s", err)
		}

		return
	}

	err = insertWallet(&wallet, &playerId, db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write(buildHttpErrorMessage("Unable to insert wallet"))
		if err != nil {
			log.Fatalf("Unable to insert wallet : %s", err)
		}

		return
	}

	// Return response
}
