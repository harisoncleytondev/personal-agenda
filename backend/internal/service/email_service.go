package service

import (
	"bytes"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/harisoncleytondev/personal-agenda/config"
	"github.com/harisoncleytondev/personal-agenda/internal/model"
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

func (e *EmailService) SendAppointmentAlert(toEmail, userName string, ap model.AlertInfo) error {
	subject := "Sua agenda pessoal: Novo Compromisso"

	hora := "O dia todo"
	if ap.TimeStart != nil {
		hora = *ap.TimeStart
	}
	desc := "Sem descrição adicional."
	if ap.Description != nil && *ap.Description != "" {
		desc = *ap.Description
	}
	dataFormatada := ap.TaskDate.Format("02/01/2006")

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
			.section-title { font-size: 16px; font-weight: 600; margin-top: 20px; color: #4f46e5; border-bottom: 2px solid #e5e7eb; padding-bottom: 5px; }
			.card { background-color: #f8fafc; border-left: 4px solid #4f46e5; padding: 15px; margin: 10px 0; border-radius: 0 8px 8px 0; }
			.item { margin-bottom: 5px; }
			.label { font-weight: 600; color: #4b5563; }
			.footer { background-color: #f9fafb; padding: 20px; text-align: center; font-size: 13px; color: #6b7280; border-top: 1px solid #e5e7eb; }
		</style>
	</head>
	<body>
		<div class="container">
			<div class="header">
				<h1>Novo Compromisso</h1>
			</div>
			<div class="content">
				<div class="greeting">Olá, %s!</div>
				<p>Você criou um novo compromisso.</p>
				<div class="section-title">Detalhes do compromisso:</div>
				<div class="card">
					<div class="item"><span class="label">Compromisso:</span> %s</div>
					<div class="item"><span class="label">Data:</span> %s</div>
					<div class="item"><span class="label">Horário:</span> %s</div>
					<div class="item"><span class="label">Descrição:</span> %s</div>
				</div>
				<p style="margin-top: 20px;">Acesse Sua agenda pessoal para gerenciar seus compromissos.</p>
			</div>
			<div class="footer">
				Este é um e-mail automático, por favor não responda.<br>
				&copy; Sua agenda pessoal
			</div>
		</div>
	</body>
	</html>`, userName, ap.TaskName, dataFormatada, hora, desc)

	var body bytes.Buffer
	body.WriteString(fmt.Sprintf("From: Sua agenda pessoal <%s>\r\n", e.From))
	body.WriteString(fmt.Sprintf("To: %s\r\n", toEmail))
	body.WriteString(fmt.Sprintf("Subject: =?utf-8?q?%s?=\r\n", strings.ReplaceAll(subject, " ", "_")))
	body.WriteString("MIME-version: 1.0;\r\n")
	body.WriteString("Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n")
	body.WriteString(htmlBody)

	auth := smtp.PlainAuth("", e.Username, e.Password, e.Host)
	addr := fmt.Sprintf("%s:%s", e.Host, e.Port)

	return smtp.SendMail(addr, auth, e.From, []string{toEmail}, body.Bytes())
}

func (e *EmailService) SendDailySummaryAlert(toEmail, userName string, todayTasks, reminderTasks []model.AlertInfo) error {
	subject := "Sua agenda pessoal: Resumo Diário"

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
			.section-title { font-size: 16px; font-weight: 600; margin-top: 20px; color: #4f46e5; border-bottom: 2px solid #e5e7eb; padding-bottom: 5px; }
			.card { background-color: #f8fafc; border-left: 4px solid #4f46e5; padding: 15px; margin: 10px 0; border-radius: 0 8px 8px 0; }
			.item { margin-bottom: 5px; }
			.label { font-weight: 600; color: #4b5563; }
			.footer { background-color: #f9fafb; padding: 20px; text-align: center; font-size: 13px; color: #6b7280; border-top: 1px solid #e5e7eb; }
		</style>
	</head>
	<body>
		<div class="container">
			<div class="header">
				<h1>Resumo da sua agenda</h1>
			</div>
			<div class="content">
				<div class="greeting">Olá, %s!</div>
				<p>Aqui está o seu resumo de compromissos e lembretes.</p>`, userName)

	if len(todayTasks) > 0 {
		htmlBody += `<div class="section-title">Compromissos pra hoje:</div>`
		for _, task := range todayTasks {
			hora := "O dia todo"
			if task.TimeStart != nil {
				hora = *task.TimeStart
			}
			desc := "Sem descrição adicional."
			if task.Description != nil && *task.Description != "" {
				desc = *task.Description
			}
			htmlBody += fmt.Sprintf(`
				<div class="card">
					<div class="item"><span class="label">Compromisso:</span> %s</div>
					<div class="item"><span class="label">Horário:</span> %s</div>
					<div class="item"><span class="label">Descrição:</span> %s</div>
				</div>`, task.TaskName, hora, desc)
		}
	} else {
		htmlBody += `<div class="section-title">Compromissos pra hoje:</div><p>Nenhum compromisso agendado para hoje.</p>`
	}

	if len(reminderTasks) > 0 {
		htmlBody += `<div class="section-title">Lembretes:</div>`
		for _, task := range reminderTasks {
			hora := "O dia todo"
			if task.TimeStart != nil {
				hora = *task.TimeStart
			}
			desc := "Sem descrição adicional."
			if task.Description != nil && *task.Description != "" {
				desc = *task.Description
			}
			dataFormatada := task.TaskDate.Format("02/01/2006")
			htmlBody += fmt.Sprintf(`
				<div class="card">
					<div class="item"><span class="label">Compromisso:</span> %s</div>
					<div class="item"><span class="label">Data:</span> %s</div>
					<div class="item"><span class="label">Horário:</span> %s</div>
					<div class="item"><span class="label">Descrição:</span> %s</div>
				</div>`, task.TaskName, dataFormatada, hora, desc)
		}
	}

	htmlBody += `
				<p style="margin-top: 20px;">Acesse Sua agenda pessoal para gerenciar seus compromissos.</p>
			</div>
			<div class="footer">
				Este é um e-mail automático, por favor não responda.<br>
				&copy; Sua agenda pessoal
			</div>
		</div>
	</body>
	</html>`

	var body bytes.Buffer
	body.WriteString(fmt.Sprintf("From: Sua agenda pessoal <%s>\r\n", e.From))
	body.WriteString(fmt.Sprintf("To: %s\r\n", toEmail))
	body.WriteString(fmt.Sprintf("Subject: =?utf-8?q?%s?=\r\n", strings.ReplaceAll(subject, " ", "_")))
	body.WriteString("MIME-version: 1.0;\r\n")
	body.WriteString("Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n")
	body.WriteString(htmlBody)

	auth := smtp.PlainAuth("", e.Username, e.Password, e.Host)
	addr := fmt.Sprintf("%s:%s", e.Host, e.Port)

	return smtp.SendMail(addr, auth, e.From, []string{toEmail}, body.Bytes())
}