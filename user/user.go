package user

type Repository interface {
	Register(u *User) string

	Find(id string) *User
}

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	NickName  string `json:"nickName"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Country   string `json:"country"`
}
