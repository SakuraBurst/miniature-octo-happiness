package types

type User struct {
	Id       int
	Login    string
	Password string
	Balance  int64
}

type UserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserBalance struct {
	CurrentBalance float64 `json:"current"`
	Withdraw       float64 `json:"withdraw"`
}

func (ur *UserRequest) IsValid() bool {
	return !(len(ur.Login) == 0 || len(ur.Password) == 0)
}
