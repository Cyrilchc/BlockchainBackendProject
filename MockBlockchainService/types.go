package main

type MockServiceExpectation struct {
	Blockchain string `json:"blockchain"`
	Pincode    string `json:"pin_code"`
}

type MockServiceAnswer struct {
	WalletAddress   string `json:"wallet_address"`
	CurrencyCode    string `json:"currency_code"`
	CurrencyBalance string `json:"currency_balance"`
}
