package notification

import (
	"birthday-service/internal/config"
	errMsg "birthday-service/internal/err"
	empHandlers "birthday-service/internal/handlers/emp"
	subHandlers "birthday-service/internal/handlers/subs"
	"context"
	"fmt"
	"log/slog"
	"net/smtp"
	"strings"
)

func SendEmail(cfg *config.ConfigSMTP, to []string, subject, body string) error {
	auth := smtp.PlainAuth("", cfg.SMTPUsername, cfg.SMTPPassword, cfg.SMTPHost)
	msg := []byte("To: " + strings.Join(to, ",") + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	addr := fmt.Sprintf("%s:%d", cfg.SMTPHost, cfg.SMTPPort)
	return smtp.SendMail(addr, auth, cfg.SMTPUsername, to, msg)
}

func SendBirthdayNotifications(ctx context.Context, subRepository subHandlers.Sub, empRepository empHandlers.Employee, cfg *config.ConfigSMTP, log *slog.Logger) {

	employees, err := empRepository.GetUpcomingBirthdays(ctx)
	if err != nil {
		log.Error("failed to get upcoming birthdays", errMsg.Err(err))
		return
	}
	for _, employee := range employees {
		users, err := subRepository.GetSubs(ctx, employee.ID)
		if err != nil {
			log.Error("failed to get subscribers", errMsg.Err(err))
			continue
		}
		if len(users) == 0 {
			continue
		}

		var emails []string
		for _, user := range users {
			emails = append(emails, user.Email)
		}

		subject := fmt.Sprintf("It's %s's birthday soon!", employee.Name)
		body := fmt.Sprintf("Don't forget to congratulate %s on %s!", employee.Name, employee.Birthday.Format("02 January"))

		err = SendEmail(cfg, emails, subject, body)
		if err != nil {
			log.Error("failed to send email", errMsg.Err(err))
		} else {
			fmt.Println("email sent")
		}
	}
}
