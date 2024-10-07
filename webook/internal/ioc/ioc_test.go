package ioc

import (
	"example.com/mod/webook/internal/repository/dao"
	"fmt"
	"testing"
)

func TestIOC(t *testing.T) {
	db := InitDbB()
	var users []dao.User
	fmt.Printf("before========users: %p,&users: %p", users, &users)
	db.Model(&dao.User{}).Order("utime desc").Find(&users)
	for _, user := range users {
		fmt.Println(user)
	}
	fmt.Printf("after========users: %p,&users: %p", users, &users)

}
