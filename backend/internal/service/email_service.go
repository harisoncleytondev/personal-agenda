package service

import (
	"bytes"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/harisoncleytondev/personal-agenda/config"
)

type EmailService struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

func NewEmailService() *EmailService {
	return &EmailService{
		Host:     config.GetSMTPHost(),
		Port:     config.GetSMTPPort(),
		Username: config.GetSMTPUser(),
		Password: config.GetSMTPPass(),
		From:     config.GetSMTPFrom(),
	}
}

func (e *EmailService) SendAppointmentAlert(toEmail, userName, taskName, taskDate, timeStart, description string) error {
	subject := fmt.Sprintf("Lembrete: %s", taskName)

	htmlBody := fmt.Sprintf(`
	<!DOCTYPE html>
	<html lang="pt-BR">
	<head>
		<meta charset="UTF-8">
		<style>
			body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background-color: #f3f4f6; margin: 0; padding: 0; }
			.container { max-width: 600px; margin: 40px auto; background-color: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1); }
			.header { background-color: #4f46e5; color: #ffffff; padding: 30px 20px; text-align: center; }
			.header h1 { margin: 0; font-size: 24px; font-weight: 600; }
			.content { padding: 30px; color: #374151; line-height: 1.6; }
			.greeting { font-size: 18px; font-weight: 600; margin-bottom: 20px; color: #111827; }
			.card { background-color: #f8fafc; border-left: 4px solid #4f46e5; padding: 20px; margin: 20px 0; border-radius: 0 8px 8px 0; }
			.item { margin-bottom: 10px; }
			.label { font-weight: 600; color: #4b5563; }
			.footer { background-color: #f9fafb; padding: 20px; text-align: center; font-size: 13px; color: #6b7280; border-top: 1px solid #e5e7eb; }
		</style>
	</head>
	<body>
		<div class="container">
			<div class="header">
				<h1>Lembrete de Agendamento</h1>
			</div>
			<div class="content">
				<div class="greeting">Olá, %s!</div>
				<p>Este é um aviso automático de que você possui um compromisso chegando.</p>
				
				<div class="card">
					<div class="item"><span class="label">Compromisso:</span> %s</div>
					<div class="item"><span class="label">Data:</span> %s</div>
					<div class="item"><span class="label">Horário:</span> %s</div>
					<div class="item"><span class="label">Descrição:</span> %s</div>
				</div>
				
				<p>Acesse sua Personal Agenda para gerenciar seus compromissos.</p>
			</div>
			<div class="footer">
				Este é um e-mail automático, por favor não responda.<br>
				&copy; Sua agenda pessoal
			</div>
		</div>
	</body>
	</html>
	`, userName, taskName, taskDate, timeStart, description)

	var body bytes.Buffer
	body.WriteString(fmt.Sprintf("From: Personal Agenda <%s>\r\n", e.From))
	body.WriteString(fmt.Sprintf("To: %s\r\n", toEmail))
	body.WriteString(fmt.Sprintf("Subject: =?utf-8?q?%s?=\r\n", strings.ReplaceAll(subject, " ", "_")))
	body.WriteString("MIME-version: 1.0;\r\n")
	body.WriteString("Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n")
	body.WriteString(htmlBody)

	auth := smtp.PlainAuth("", e.Username, e.Password, e.Host)
	addr := fmt.Sprintf("%s:%s", e.Host, e.Port)

	err := smtp.SendMail(addr, auth, e.From, []string{toEmail}, body.Bytes())
	return err
}