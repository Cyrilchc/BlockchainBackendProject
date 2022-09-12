package main

type Player struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Pincode  string `json:"pin_code"`
}

type Wallet struct {
	Address string `json:"address"`
}

type MockServiceExpectation struct {
	Blockchain string `json:"blockchain"`
	Pincode    string `json:"pin_code"`
}

type MockServiceAnswer struct {
	WalletAddress   string `json:"wallet_address"`
	CurrencyCode    string `json:"currency_code"`
	CurrencyBalance string `json:"currency_balance"`
}

type PlayerInformation struct {
	Username      string `json:"username"`
	WalletAddress string `json:"wallet_address"`
}
