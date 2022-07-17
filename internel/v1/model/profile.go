package model

type Profile struct {
	ID      uint64 `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
}
