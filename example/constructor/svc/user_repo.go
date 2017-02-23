package svc

import "fmt"

type UserRepository struct {
	db *DB
}

func NewUserRepository(db *DB) *UserRepository {
	inst := new(UserRepository)
	inst.db = db

	return inst
}

func (ur *UserRepository) CreateUser() {
	ur.db.Connect()

	fmt.Println("  - User created")
}
