package email

import (
	"context"
	"fmt"
	"log/slog"
	"net/smtp"
	"strings"

	"github.com/tartoide/stori/stori-challenge/internal/domain"
	"github.com/tartoide/stori/stori-challenge/internal/services"
)

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

type smtpEmailService struct {
	config SMTPConfig
	logger *slog.Logger
}

// NewSMTPEmailService creates a new SMTP email service
func NewSMTPEmailService(config SMTPConfig, logger *slog.Logger) services.EmailService {
	return &smtpEmailService{
		config: config,
		logger: logger,
	}
}

// SendSummary sends an email summary to the recipient
func (s *smtpEmailService) SendSummary(ctx context.Context, recipient string, summary *domain.Summary) error {
	s.logger.Info("sending email summary", "recipient", recipient)

	htmlBody, err := s.RenderTemplate(summary)
	if err != nil {
		return fmt.Errorf("rendering email template: %w", err)
	}

	message := s.createEmailMessage(recipient, "Stori - Your Account Summary", htmlBody)

	if err := s.sendEmail(recipient, message); err != nil {
		return fmt.Errorf("sending email to %s: %w", recipient, err)
	}

	s.logger.Info("email summary sent successfully", "recipient", recipient)
	return nil
}

// RenderTemplate renders the email template with summary data
func (s *smtpEmailService) RenderTemplate(summary *domain.Summary) (string, error) {
	return RenderEmailTemplate(summary)
}

func (s *smtpEmailService) createEmailMessage(to, subject, htmlBody string) string {
	headers := make(map[string]string)
	headers["From"] = s.config.From
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=utf-8"

	var message strings.Builder
	for key, value := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	message.WriteString("\r\n")
	message.WriteString(htmlBody)

	return message.String()
}

func (s *smtpEmailService) sendEmail(recipient, message string) error {
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)

	addr := s.config.Host + ":" + s.config.Port
	recipients := []string{recipient}

	s.logger.Debug("connecting to SMTP server", "host", s.config.Host, "port", s.config.Port)

	err := smtp.SendMail(addr, auth, s.config.From, recipients, []byte(message))
	if err != nil {
		s.logger.Error("SMTP send failed", "error", err, "recipient", recipient)
		return domain.ErrEmailDeliveryFailed
	}

	return nil
}
