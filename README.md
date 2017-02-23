# Gontainer, DI container in Golang
###### Code written by lobotomized squirrel squad under tons of cocaine
###### Don't use it in your application, to avoid getting awkwards 

## How to
Basic usage, for more look in example folder
```go
package main

import (
	"fmt"
	"github.com/russell-kvashnin/gontainer"
)

type DB struct {
	config map[string]string
}

func NewDB(config map[string]string) *DB {
	inst := new(DB)
	inst.config = config

	return inst
}

func (db *DB) Connect() {
	fmt.Println("  - Database connected")
}

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


func main() {
	dbConfig := map[string]string{
		"host": "localhost",
		"username": "root",
		"password": "toor",
	}
	dbDef := gontainer.ServiceDefinition{
		Name: "app.db",
		Factory: gontainer.Factory{
			Constructor: NewDB,
			Args: gontainer.ConstructorArguments{
				dbConfig,
			},
		},
	}
	userRepoDef := gontainer.ServiceDefinition{
		Name: "app.repo.user",
		Factory: gontainer.Factory{
			Constructor: NewUserRepository,
			Args: gontainer.ConstructorArguments{
				gontainer.Injection("app.db"),
			},
		},
	}

	defs := gontainer.ServiceDefinitions{
		dbDef,
		userRepoDef,
	}

	container := gontainer.NewContainer()
	errs := container.Compile(defs)
	if errs != nil {
		fmt.Println(errs)
	}

	userRepo, err := container.Get("app.repo.user")
	if err != nil {
		fmt.Println(err.Error())
	}

	userRepo.(*UserRepository).CreateUser()
}
```

