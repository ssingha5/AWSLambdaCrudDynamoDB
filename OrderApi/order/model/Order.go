package model

type Order struct {
	Id       int    `json:"id"`
	ItemName string `json:"itemName"`
	Quantity int    `json:"quantity"`
}
