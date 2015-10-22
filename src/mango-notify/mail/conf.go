package mail

//Conf for sending mail.
type Conf struct {
	From     string
	Password string
	To       string

	// Connect to the SMTP server, valid example: smtp.gmail.com:465
	ServerName string
}
