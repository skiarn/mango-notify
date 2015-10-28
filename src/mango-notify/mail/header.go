package mail

import (
	"fmt"
	"net/mail"
	"path/filepath"
)

func makeHeaders(m *Mail) string {
	// Setup headers
	from := mail.Address{"", m.Conf.From}
	to := mail.Address{"", m.Conf.To}

	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = m.Subject
	headers["MIME-Version"] = "1.0"
	marker := "FILEWATCHBOUNDARY"
	headers["Content-Type"] = "multipart/mixed; boundary=" + marker

	// Setup headers
	mHeaders := ""
	for k, v := range headers {
		mHeaders += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	mHeaders += fmt.Sprintf("\r\n--%s", marker)
	mHeaders += "\r\n"

	return mHeaders
}

func makeAttachmentHeaders() string {
	var filename = "report.mgofy"
	var extension = filepath.Ext(filename)
	var extNoDot = extension[1:len(extension)]

	headers := make(map[string]string)
	headers["Content-Type"] = "application/" + extNoDot + fmt.Sprintf("; name=\"%s\"", filename)
	headers["Content-Transfer-Encoding"] = "text/plain"

	headers["Content-Disposition"] = fmt.Sprintf(" attachment; filename=\"%s\"", filename)
	//part 3 will be the attachment

	mHeaders := ""
	for k, v := range headers {
		mHeaders += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	mHeaders += "\r\n"

	return mHeaders
}
