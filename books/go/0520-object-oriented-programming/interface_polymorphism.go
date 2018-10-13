package main

import (
	"fmt"
)

type Runner interface {
	Run()
}

type Admin struct {
	Username, Password string
}

func (admin Admin) Run() {
	fmt.Println("Admin ==> Run()")
}

type User struct {
	ID              uint64
	FullName, Email string
}

func (user User) Run() {
	fmt.Println("User ==> Run()")
}

// RunnerExample takes any type that fullfils the Runner interface
func RunnerExample(r Runner) {
	r.Run()
}

func main() {
	admin := Admin{
		"zola",
		"supersecretpassword",
	}

	user := User{
		1,
		"Zelalem Mekonen",
		"zola.mk.27@gmail.com",
	}

	RunnerExample(admin)

	RunnerExample(user)
}
