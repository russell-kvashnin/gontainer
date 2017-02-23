package setter

import (
	"fmt"

	"github.com/russell-kvashnin/gontainer"
	"github.com/russell-kvashnin/gontainer/example/setter/svc"
)

func SetterInjectionExamples() {
	fmt.Println("Gontainer setter injection examples \n ")

	defs := getServicesDefinitions()

	container := gontainer.NewContainer()
	errs := container.Compile(defs)

	if errs != nil {
		fmt.Println(errs)
	}

	nlService, err := container.Get("app.service.newsletter")
	if err != nil {
		fmt.Println(err.Error())
	}

	nlService.(*svc.NewsService).SendNewsletter()

	fmt.Println("")
}

func getServicesDefinitions() gontainer.ServiceDefinitions {
	mailer := gontainer.ServiceDefinition{
		Name: "app.mailer",
		Factory: gontainer.Factory{
			Constructor: svc.NewMailer,
		},
	}

	nlService := gontainer.ServiceDefinition{
		Name: "app.service.newsletter",
		Factory: gontainer.Factory{
			Constructor: svc.NewNewsService,
		},
	}

	svcs := gontainer.ServiceDefinitions{
		mailer,
		nlService,
	}

	return svcs
}