package constructor

import (
	"fmt"

	"github.com/russell-kvashnin/gontainer"
	"github.com/russell-kvashnin/gontainer/example/constructor/svc"
)

func ConstructorInjectionExamples() {
	fmt.Println("Gontainer constructor injection examples \n ")

	defs := getServicesDefinitions()

	container := gontainer.NewContainer()
	errs := container.Compile(defs)
	if errs != nil {
		fmt.Println(errs)
	}

	us, err := container.Get("app.service.user")
	if err != nil {
		fmt.Println(err.Error())
	}

	us.(*svc.UserService).RegisterUser()

	fmt.Println("")
}

func getServicesDefinitions() gontainer.ServiceDefinitions {

	dbConfig := map[string]string{
		"host": "localhost",
		"username": "root",
		"password": "toor",
	}
	db := gontainer.ServiceDefinition{
		Name: "app.db",
		Factory: gontainer.Factory{
			Constructor: svc.NewDB,
			Args: gontainer.ConstructorArguments{
				dbConfig,
			},
		},
	}

	mailerConfig := map[string]string{
		"transport": "smtp",
	}
	mailer := gontainer.ServiceDefinition{
		Name: "app.mailer",
		Factory: gontainer.Factory{
			Constructor: svc.NewMailer,
			Args: gontainer.ConstructorArguments{
				mailerConfig,
			},
		},
	}

	userRepo := gontainer.ServiceDefinition{
		Name: "app.repo.user",
		Factory: gontainer.Factory{
			Constructor: svc.NewUserRepository,
			Args: gontainer.ConstructorArguments{
				gontainer.Injection("app.db"),
			},
		},
	}

	userService := gontainer.ServiceDefinition{
		Name: "app.service.user",
		Factory: gontainer.Factory{
			Constructor: svc.NewUserService,
			Args: gontainer.ConstructorArguments{
				gontainer.Injection("app.repo.user"),
				gontainer.Injection("app.mailer"),
			},
		},
	}

	defs := gontainer.ServiceDefinitions{
		db,
		mailer,
		userRepo,
		userService,
	}

	return defs
}
