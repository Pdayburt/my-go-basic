package control

import (
	"fmt"
	"testing"
)

func TestLoop(t *testing.T) {
	users := []User{
		{
			Name: "Jack",
		},
		{
			Name: "Rose",
		},
	}
	m := make(map[string]*User)
	for _, user := range users {
		fmt.Printf("user :%p\n", &user)
		m[user.Name] = &user
	}
}

type User struct {
	Name string
}
