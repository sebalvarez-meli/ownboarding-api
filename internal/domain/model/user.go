package model

type User struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Token Token  `json:"token"`
}
