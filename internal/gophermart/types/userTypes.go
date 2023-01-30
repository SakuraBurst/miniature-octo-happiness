package types

type User struct {
	ID       int
	Login    string
	Password string
	Balance  float64
	Withdraw float64
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

func (u *User) AddBalance(sum float64) {
	u.Balance = u.Balance + sum
}

func (u *User) WithdrawBalance(requestedSum float64) {
	u.Balance = u.Balance - requestedSum
	u.Withdraw = u.Withdraw + requestedSum
}
