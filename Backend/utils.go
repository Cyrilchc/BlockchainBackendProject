package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func checkUsername(username string) error {
	if !regexUsername.MatchString(username) {
		return errors.New("username does not respect the naming policy")
	}

	return nil
}

func checkPassword(password string) error {
	if !regexPassword.MatchString(password) {
		return errors.New("password does not respect the complexity policy")
	}

	return nil
}

func checkPinCode(pincode string) error {
	if !regexPinCode.MatchString(pincode) {
		return errors.New("pin code must contains 6 digits")
	}

	return nil
}

func buildHttpErrorMessage(message string) ([]byte, error) {
	resp := make(map[string]string)
	resp["error"] = message
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	return jsonResp, nil
}

func checkHttpMethod(method string, w http.ResponseWriter, r *http.Request) error {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		message, err := buildHttpErrorMessage("Only POST method is allowed")
		if err != nil {
			return err
		}

		_, err = w.Write(message)
		if err != nil {
			return err
		}

		return errors.New("method not allowed")
	}

	return nil
}

func connectDatabase() (*sql.DB, error) {
	con := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", DB_HOST, PORT, USER, PASSWORD, DBNAME)
	return sql.Open("postgres", con)
}

func sendHttpError(httpErrorCode int, message string, w http.ResponseWriter, err error) error {
	w.WriteHeader(httpErrorCode)
	errorMessage, err := buildHttpErrorMessage(message)
	if err != nil {
		return err
	}

	_, err = w.Write(errorMessage)
	return err
}

func insertPlayer(player *Player, db *sql.DB) (int, error) {
	query := `insert into "players"("username","password","pincode","jsondata") values ($1, $2, $3, $4) RETURNING id`
	jsonPlayer, err := json.Marshal(player)
	if err != nil {
		return 0, nil
	}

	insertedId := 0
	err = db.QueryRow(query, player.Username, player.Password, player.Pincode, jsonPlayer).Scan(&insertedId)
	if err != nil {
		return 0, err
	}

	return insertedId, nil
}

func insertWallet(wallet *Wallet, id *int, db *sql.DB) error {
	query := `insert into "wallets"("address","id_player","jsondata") values ($1, $2, $3)`
	jsonWallet, err := json.Marshal(wallet)
	if err != nil {
		return err
	}

	_, err = db.Exec(query, wallet.Address, id, jsonWallet)
	return err
}
