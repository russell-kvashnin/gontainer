package svc

import "fmt"

type Mailer struct {
	config map[string]string
}

func NewMailer(config map[string]string) *Mailer {
	inst := new(Mailer)
	inst.config = config

	return inst
}

func (m *Mailer) SendMail() {
	fmt.Println("  - Mail sent")
}