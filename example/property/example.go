package property

import (
	"fmt"

	"github.com/russell-kvashnin/gontainer"
	"github.com/russell-kvashnin/gontainer/example/property/svc"
)

func PropertyInjectionExamples() {
	fmt.Println("Gontainer property injection examples \n ")

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
