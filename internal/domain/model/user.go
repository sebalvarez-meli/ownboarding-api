package model

type User struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Token Token  `json:"token"`
}
