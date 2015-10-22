package mail

import (
	"crypto/tls"
	"log"
	"net"
	"net/smtp"
)

//Mail to be sent.
type Mail struct {
	Subject string
	//for HTML emails just put HTML in the body
	Body       Body
	Attachment Attachment
	Conf       Conf
}

//Send should validate and send mail.
func (m *Mail) Send() {
	//Setup Headers
	message := makeHeaders(m) + m.Body.Build() + m.Attachment.Build()

	host, _, _ := net.SplitHostPort(m.Conf.ServerName)
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	conn, err := tls.Dial("tcp", m.Conf.ServerName, tlsconfig)
	if err != nil {
		log.Panic(err)
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		log.Panic(err)
	}

	if err = c.Auth(smtp.PlainAuth("", m.Conf.From, m.Conf.Password, host)); err != nil {
		log.Panic(err)
	}

	if err = c.Mail(m.Conf.From); err != nil {
		log.Panic(err)
	}

	if err = c.Rcpt(m.Conf.To); err != nil {
		log.Panic(err)
	}

	w, err := c.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	c.Quit()
}
