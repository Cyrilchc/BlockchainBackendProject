package main

type Player struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Pincode  string `json:"pin_code"`
}

type Wallet struct {
	Address string `json:"address"`
}
