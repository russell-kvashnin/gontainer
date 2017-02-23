package svc

import "fmt"

type NewsService struct {
	mailer *Mailer `inject:"app.mailer" inject_type:"setter" inject_method:"SetMailer"`
}

func NewNewsService() *NewsService {
	inst := new(NewsService)

	return inst
}

func (ns *NewsService) SetMailer(mailer *Mailer) {
	ns.mailer = mailer
}

func (ns *NewsService) SendNewsletter() {
	fmt.Println(" ##### Doing some newsletter hardwork")

	ns.mailer.SendMail()

	fmt.Println(" ##### Newsletter sent")
}