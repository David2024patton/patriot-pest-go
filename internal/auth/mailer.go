package auth

// Mailer — outbound email for login codes and campaign sends. Port of the
// PHP app/Auth/Mailer.php: two transports, dev log (APP_ENV=local or no SMTP
// host) and a minimal SSL SMTP client (EHLO / AUTH LOGIN / MAIL / RCPT / DATA).
// Dev mode writes storage/logs/mail-YYYY-MM-DD.log so OTP codes stay readable
// without a mailbox; production never writes that file.

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Mailer holds SMTP credentials. Dev true when APP_ENV=local or host is empty.
type Mailer struct {
	Host     string
	Port     int
	User     string
	Pass     string
	From     string // MAIL_FROM_ADDRESS
	FromName string
	Dev      bool
}

// NewMailer builds a Mailer from config values. Dev mode is decided by the
// caller (config.Env == "local").
func NewMailer(host string, port int, user, pass, from, fromName string) *Mailer {
	if port == 0 {
		port = 465
	}
	if from == "" {
		from = "no-reply@patriotpest.pro"
	}
	if fromName == "" {
		fromName = "Patriot Pest Control"
	}
	return &Mailer{Host: host, Port: port, User: user, Pass: pass, From: from, FromName: fromName, Dev: host == "" || os.Getenv("APP_ENV") == "local"}
}

// Send delivers the message (or logs it in dev mode). Returns true on success.
func (m *Mailer) Send(to, subject, bodyHTML string) bool {
	if m.Dev {
		m.log(to, subject, bodyHTML)
		return true
	}
	return m.smtp(to, subject, bodyHTML)
}

// log appends the message to today's mail log (dev/debug affordance).
func (m *Mailer) log(to, subject, body string) {
	dir := filepath.Join("storage", "logs")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return
	}
	entry := strings.Repeat("=", 70) + "\n" +
		"[" + time.Now().Format(time.RFC3339) + "] TO: " + to + "\nSUBJECT: " + subject + "\n" +
		strings.Repeat("-", 70) + "\n" + body + "\n\n"
	if err := appendFile(filepath.Join(dir, "mail-"+time.Now().Format("2006-01-02")+".log"), entry); err != nil {
		fmt.Printf("mailer: log write failed: %v\n", err)
	}
}

func appendFile(path, s string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, _ = f.WriteString(s)
	return nil
}

// smtp sends over SSL (port 465 implicit TLS). Mirrors the PHP client exactly:
// greeting, EHLO, optional AUTH LOGIN, MAIL/RCPT/DATA, QUIT.
func (m *Mailer) smtp(to, subject, body string) bool {
	host := m.Host
	port := m.Port
	remote := net.JoinHostPort(host, fmt.Sprintf("%d", port))

	conn, err := net.Dial("tcp", remote)
	if err != nil {
		fmt.Printf("mailer: connect failed %v\n", err)
		return false
	}
	tlsConn := tls.Client(conn, &tls.Config{ServerName: host})
	if err := tlsConn.Handshake(); err != nil {
		conn.Close()
		fmt.Printf("mailer: TLS handshake failed %v\n", err)
		return false
	}
	defer conn.Close()

	read := func() string {
		br := bufio.NewReader(tlsConn)
		var data string
		for {
			line, err := br.ReadString('\n')
			if err != nil {
				break
			}
			data += line
			if len(line) < 4 || !strings.Contains(line[3:], "-") {
				break
			}
		}
		return data
	}
	send := func(cmd string) string {
		_, _ = tlsConn.Write([]byte(cmd + "\r\n"))
		return read()
	}

	read() // greeting
	send("EHLO patriotpest.pro")
	if m.User != "" {
		send("AUTH LOGIN")
		send(base64.StdEncoding.EncodeToString([]byte(m.User)))
		send(base64.StdEncoding.EncodeToString([]byte(m.Pass)))
	}
	send("MAIL FROM:<" + m.From + ">")
	send("RCPT TO:<" + to + ">")
	send("DATA")

	headers := "From: " + m.FromName + " <" + m.From + ">\r\n" +
		"To: <" + to + ">\r\n" +
		"Subject: " + subject + "\r\n" +
		"Date: " + time.Now().Format(time.RFC1123) + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n"
	_, _ = tlsConn.Write([]byte(headers + "\r\n" + body + "\r\n.\r\n"))
	resp := read()
	send("QUIT")

	if ok := strings.HasPrefix(strings.TrimSpace(resp), "250"); !ok {
		fmt.Printf("mailer: DATA rejected: %s\n", strings.TrimSpace(resp))
		return false
	}
	return true
}

// MailTemplate wraps a message body in the branded email template
// (port of Mailer::template). Used for OTP + campaign sends.
func MailTemplate(heading, innerHTML string) string {
	unsub := ""
	return `<div style="font-family:Arial,Helvetica,sans-serif;background:#12140d;padding:24px">` +
		`<div style="max-width:560px;margin:0 auto;background:#1b1e14;border:1px solid #3a3f2c;border-radius:6px;padding:32px;color:#e8e6da">` +
		`<div style="font-size:13px;letter-spacing:2px;color:#c8a24a;font-weight:bold;text-transform:uppercase">★ Patriot Pest Control</div>` +
		`<h1 style="font-size:22px;color:#f2efe2;margin:18px 0 10px">` + heading + `</h1>` +
		innerHTML + unsub +
		`<p style="font-size:12px;color:#8a8f7a;margin-top:28px">Veteran-owned · WA / ID / OR / AZ · (509) 471-5767</p>` +
		`</div></div>`
}
