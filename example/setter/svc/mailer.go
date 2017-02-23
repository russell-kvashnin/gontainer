package svc

import "fmt"

type Mailer struct {
}

func NewMailer() *Mailer {
	inst := new(Mailer)

	return inst
}

func (m *Mailer) SendMail() {
	fmt.Println("  - Mail sent")
}