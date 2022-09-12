package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"io"
	"log"
	"net/http"
	"regexp"
)

const (
	SERVER_PORT           = "5000"
	EXTERNAL_SERVICE_PORT = "5001"
	HOST                  = "localhost"
	PORT                  = 5432
	USER                  = "postgres"
	PASSWORD              = "cyril"
	DBNAME                = "blockchainbackend"
)

var regexUsername = regexp.MustCompile(`^[a-z0-9_]{3,100}$`)
var regexPassword = regexp.MustCompile(`^.{6,32}$`)
var regexPinCode = regexp.MustCompile(`^\d{6}$`)

func main() {
	http.HandleFunc("/", createUserAndWallet)
	//http.HandleFunc("/wallets/create", mockCreateWallet)
	err := http.ListenAndServe(fmt.Sprintf(":%s", SERVER_PORT), nil)
	if err != nil {
		fmt.Print(err)
	}
}

/*
*
Creates a user and a wallet
*/
func createUserAndWallet(w http.ResponseWriter, r *http.Request) {
	db, err := connectDatabase()
	if err != nil {
		sendHttpError(http.StatusInternalServerError, "Unable to connect to database", w, err)
		return
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Print("Unable to close database connection")
		}
	}(db)

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
			log.Print("IO error, Unable to close http request body : %s", err)
		}
	}(r.Body)

	// Map body to struct
	player := Player{}
	err = json.Unmarshal(body, &player)
	if err != nil {
		sendHttpError(http.StatusInternalServerError, "Unable to serialize body", w, err)
		return
	}

	// Uniqueness of username will be checked by database constraints

	// Username naming policy
	if !regexUsername.MatchString(player.Username) {
		sendHttpError(http.StatusBadRequest, "Username does not respect the naming policy", w, err)
		return
	}

	// Password policy
	if !regexPassword.MatchString(player.Password) {
		sendHttpError(http.StatusBadRequest, "Password does not respect the complexity policy", w, err)
		return
	}

	// Pincode policy
	if !regexPinCode.MatchString(player.Pincode) {
		sendHttpError(http.StatusBadRequest, "Pincode must contains 6 digits", w, err)
		return
	}

	// Prepare external service call
	MockServiceExpectation := MockServiceExpectation{
		Pincode:    player.Pincode,
		Blockchain: "ETH",
	}

	data, err := json.Marshal(MockServiceExpectation)
	if err != nil {
		sendHttpError(http.StatusInternalServerError, "Unable to serialize MockServiceExpectation", w, err)
		return
	}

	// Send request to external service
	response, err := http.Post(fmt.Sprintf("http://%s:%s/wallets/create", HOST, EXTERNAL_SERVICE_PORT), "application/json", bytes.NewBuffer(data))
	if err != nil {
		sendHttpError(http.StatusInternalServerError, "External service returned an error", w, err)
		return
	}

	// Read body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		sendHttpError(http.StatusInternalServerError, "Unable to read body", w, err)
		return
	}

	// Check error
	if response.StatusCode != 200 {
		sendHttpError(http.StatusInternalServerError, "External service returned an error", w, errors.New(string(responseBody)))
		return
	}

	// Handle body closure
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Print("IO error, Unable to close http request body : %s", err)
		}
	}(response.Body)

	// Map body to struct
	mockServiceAnswer := MockServiceAnswer{}
	err = json.Unmarshal(responseBody, &mockServiceAnswer)
	if err != nil {
		sendHttpError(http.StatusInternalServerError, "Unable to serialize body", w, err)
		return
	}

	// Create wallet object
	wallet := Wallet{
		Address: mockServiceAnswer.WalletAddress,
	}

	// Insert player in database
	playerId, err := insertPlayer(&player, db)
	if err != nil {
		sendHttpError(http.StatusInternalServerError, "Unable to insert player", w, err)
		return
	}

	// Insert Wallet in database
	err = insertWallet(&wallet, &playerId, db)
	if err != nil {
		sendHttpError(http.StatusInternalServerError, "Unable to insert wallet", w, err)
		return
	}

	// Build response
	playerInformation := PlayerInformation{
		Username:      player.Username,
		WalletAddress: wallet.Address,
	}

	playerInformationJson, err := json.Marshal(playerInformation)
	if err != nil {
		sendHttpError(http.StatusInternalServerError, "Unable to serialize response", w, err)
		return
	}

	// Send successful response
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(playerInformationJson)
	if err != nil {
		log.Print("Unable to write response : %s", err)
	}
}
