package svc

import "fmt"

type NewsService struct {
	Mailer *Mailer `inject:"app.mailer" inject_type:"property"`
}

func NewNewsService() *NewsService {
	inst := new(NewsService)

	return inst
}

func (ns *NewsService) SendNewsletter() {
	fmt.Println(" ##### Doing some newsletter hardwork")

	ns.Mailer.SendMail()

	fmt.Println(" ##### Newsletter sent")
}
