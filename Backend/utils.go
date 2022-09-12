package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func buildHttpErrorMessage(message string) []byte {
	resp := make(map[string]string)
	resp["error"] = message
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Print(err)
	}

	return jsonResp
}

func checkHttpMethod(method string, w http.ResponseWriter, r *http.Request) error {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, err := w.Write(buildHttpErrorMessage("Only POST method is allowed"))
		if err != nil {
			log.Printf("Unable to write to http response : %s", err)
		}

		return errors.New("Method not allowed")
	}

	return nil
}

func connectDatabase() (*sql.DB, error) {
	con := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", DB_HOST, PORT, USER, PASSWORD, DBNAME)
	return sql.Open("postgres", con)
}

func sendHttpError(httpErrorCode int, message string, w http.ResponseWriter, err error) {
	log.Print(err)
	w.WriteHeader(httpErrorCode)
	_, err = w.Write(buildHttpErrorMessage(message))
	if err != nil {
		log.Printf("%s : %s", message, err)
	}
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
