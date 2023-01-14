package types

type User struct {
	Id       int
	Login    string
	Password string
	Balance  int
}

type UserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (ur *UserRequest) IsValid() bool {
	return !(len(ur.Login) == 0 || len(ur.Password) == 0)
}
