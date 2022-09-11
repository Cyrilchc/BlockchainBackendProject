package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func buildHttpErrorMessage(message string) []byte {
	resp := make(map[string]string)
	resp["error"] = message
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatal(err)
	}

	return jsonResp
}

func checkHttpMethod(method string, w http.ResponseWriter, r *http.Request) {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, err := w.Write(buildHttpErrorMessage("Only POST method is allowed"))
		if err != nil {
			log.Fatalf("Unable to write to http response : %s", err)
		}

		return
	}
}

func connectDatabase(host string, port int, user string, password string, dbname string) (*sql.DB, error) {
	con := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", HOST, PORT, USER, PASSWORD, DBNAME)
	return sql.Open("postgres", con)
}

func insertPlayer(player *Player, db *sql.DB) (sql.Result, error) {
	query := `insert into "players"("username","password","pincode","jsondata") values ($1, $2, $3, $4)`
	jsonPlayer, err := json.Marshal(player)
	if err != nil {
		return nil, err
	}

	result, err := db.Exec(query, player.Username, player.Password, player.Pincode, jsonPlayer)
	return result, err
}

func insertWallet(wallet *Wallet, id *int64, db *sql.DB) error {
	query := `insert into "wallets"("address","id_player","jsondata") values ($1, $2, $3)`
	jsonWallet, err := json.Marshal(wallet)
	if err != nil {
		return err
	}

	_, err = db.Exec(query, wallet.Address, id, jsonWallet)
	return err
}
