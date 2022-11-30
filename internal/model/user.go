package model

type User struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
	Balance    int    `json:"balance"`
}

type CreateTransactionRequest struct {
	ID     string `json:"id"`
	Amount int    `json:"amount"`
}
