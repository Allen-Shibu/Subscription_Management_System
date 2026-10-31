// Package email sends transactional email over SMTP.
package email

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"strconv"
	"time"
)

// Message is one transactional email. Attachment may be nil.
type Message struct {
	To             string
	Subject        string
	Text           string
	HTML           string
	AttachmentName string
	Attachment     []byte
}

// Sender delivers email messages.
type Sender interface {
	Send(ctx context.Context, msg Message) error
}

// Config holds SMTP connection settings.
type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	Timeout  time.Duration
}

type SMTPSender struct {
	cfg Config
}

// NoopSender discards messages, suitable for tests or when email delivery is disabled.
type NoopSender struct{}

func (n *NoopSender) Send(ctx context.Context, msg Message) error {
	return nil
}

// NewSMTP returns an SMTP-backed Sender. Port defaults to 587 and the
// timeout defaults to 15 seconds.
func NewSMTP(cfg Config) *SMTPSender {
	if cfg.Port == 0 {
		cfg.Port = 587
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 15 * time.Second
	}
	return &SMTPSender{cfg: cfg}
}

func (s *SMTPSender) Send(ctx context.Context, msg Message) error {
	if s.cfg.Host == "" {
		return errors.New("SMTP is not configured")
	}
	if msg.To == "" {
		return errors.New("email recipient is required")
	}
	if s.cfg.From == "" {
		return errors.New("email sender is required")
	}

	raw, err := buildMessage(s.cfg.From, msg)
	if err != nil {
		return err
	}

	addr := net.JoinHostPort(s.cfg.Host, strconv.Itoa(s.cfg.Port))
	dialer := &net.Dialer{Timeout: s.cfg.Timeout}
	implicitTLS := s.cfg.Port == 465
	var conn net.Conn
	if implicitTLS {
		tlsDialer := &tls.Dialer{
			NetDialer: dialer,
			Config: &tls.Config{
				ServerName: s.cfg.Host,
				MinVersion: tls.VersionTLS12,
			},
		}
		conn, err = tlsDialer.DialContext(ctx, "tcp", addr)
	} else {
		conn, err = dialer.DialContext(ctx, "tcp", addr)
	}
	if err != nil {
		return fmt.Errorf("connect to smtp server: %w", err)
	}
	if deadline, ok := ctx.Deadline(); ok {
		_ = conn.SetDeadline(deadline)
	} else {
		_ = conn.SetDeadline(time.Now().Add(s.cfg.Timeout))
	}

	client, err := smtp.NewClient(conn, s.cfg.Host)
	if err != nil {
		_ = conn.Close()
		return fmt.Errorf("start smtp session: %w", err)
	}
	defer func() { _ = client.Close() }()

	if implicitTLS {
		// Already wrapped in TLS above.
	} else if ok, _ := client.Extension("STARTTLS"); ok {
		if err := client.StartTLS(&tls.Config{ServerName: s.cfg.Host, MinVersion: tls.VersionTLS12}); err != nil {
			return fmt.Errorf("start tls: %w", err)
		}
	}

	if s.cfg.Username != "" {
		auth := smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("smtp authentication: %w", err)
		}
	}

	from := envelopeAddress(s.cfg.From)
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("set sender: %w", err)
	}
	to := envelopeAddress(msg.To)
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("set recipient: %w", err)
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("start message body: %w", err)
	}
	if _, err := writer.Write(raw); err != nil {
		_ = writer.Close()
		return fmt.Errorf("write message body: %w", err)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("finish message body: %w", err)
	}
	if err := client.Quit(); err != nil {
		return fmt.Errorf("close smtp session: %w", err)
	}
	return nil
}

func envelopeAddress(from string) string {
	addr, err := mail.ParseAddress(from)
	if err != nil {
		return from
	}
	return addr.Address
}

func buildMessage(from string, msg Message) ([]byte, error) {
	var body bytes.Buffer
	mw := multipart.NewWriter(&body)

	if msg.Text != "" {
		part, err := mw.CreatePart(textHeader(`text/plain; charset="utf-8"`))
		if err != nil {
			return nil, err
		}
		if err := writeQuotedPrintable(part, msg.Text); err != nil {
			return nil, err
		}
	}
	if msg.HTML != "" {
		part, err := mw.CreatePart(textHeader(`text/html; charset="utf-8"`))
		if err != nil {
			return nil, err
		}
		if err := writeQuotedPrintable(part, msg.HTML); err != nil {
			return nil, err
		}
	}
	if len(msg.Attachment) > 0 {
		name := msg.AttachmentName
		if name == "" {
			name = "attachment"
		}
		header := textproto.MIMEHeader{}
		header.Set("Content-Type", "application/pdf")
		header.Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": name}))
		header.Set("Content-Transfer-Encoding", "base64")
		part, err := mw.CreatePart(header)
		if err != nil {
			return nil, err
		}
		if err := writeBase64(part, msg.Attachment); err != nil {
			return nil, err
		}
	}

	if err := mw.Close(); err != nil {
		return nil, err
	}

	var out bytes.Buffer
	fmt.Fprintf(&out, "From: %s\r\n", from)
	fmt.Fprintf(&out, "To: %s\r\n", msg.To)
	fmt.Fprintf(&out, "Subject: %s\r\n", mime.QEncoding.Encode("utf-8", msg.Subject))
	fmt.Fprintf(&out, "Date: %s\r\n", time.Now().UTC().Format(time.RFC1123Z))
	out.WriteString("MIME-Version: 1.0\r\n")
	fmt.Fprintf(&out, "Content-Type: multipart/mixed; boundary=%q\r\n", mw.Boundary())
	out.WriteString("\r\n")
	out.Write(body.Bytes())
	return out.Bytes(), nil
}

func textHeader(contentType string) textproto.MIMEHeader {
	header := textproto.MIMEHeader{}
	header.Set("Content-Type", contentType)
	header.Set("Content-Transfer-Encoding", "quoted-printable")
	return header
}

func writeQuotedPrintable(w io.Writer, value string) error {
	qp := quotedprintable.NewWriter(w)
	if _, err := qp.Write([]byte(value)); err != nil {
		return err
	}
	return qp.Close()
}

func writeBase64(w io.Writer, data []byte) error {
	encoded := base64.StdEncoding.EncodeToString(data)
	for len(encoded) > 76 {
		if _, err := io.WriteString(w, encoded[:76]+"\r\n"); err != nil {
			return err
		}
		encoded = encoded[76:]
	}
	if encoded != "" {
		if _, err := io.WriteString(w, encoded+"\r\n"); err != nil {
			return err
		}
	}
	return nil
}
