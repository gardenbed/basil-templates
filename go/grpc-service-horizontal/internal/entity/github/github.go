package github

import "fmt"

// User is the entity for a GitHub user.
type User struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// String implements the fmt.Stringer interface.
func (u *User) String() string {
	return fmt.Sprintf("User{id=%d login=%s email=%s name=%s}", u.ID, u.Login, u.Email, u.Name)
}
