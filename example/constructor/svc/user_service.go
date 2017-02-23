package svc

import "fmt"

type UserService struct {
	mailer *Mailer
	repo *UserRepository
}

func NewUserService(repo *UserRepository, mailer *Mailer) *UserService {
	inst := new(UserService)
	inst.repo = repo
	inst.mailer = mailer

	return inst
}

func (us UserService) RegisterUser() {
	fmt.Println(" ##### Starting user registration")

	us.repo.CreateUser()
	us.mailer.SendMail()

	fmt.Println(" ##### User registered ")
}