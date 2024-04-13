package model

type AuthRequestDto struct {
	Password string `json:"password"`
}

type AuthSuccessDto struct {
	Token string `json:"token"`
}

type AuthFailedDto struct {
	Error string `json:"error"`
}
