package mail

import "fmt"

//Body of the mail content.
type Body struct {
	Content string
}

//Build body part of the mail.
func (b *Body) Build() string {
	body := make(map[string]string)
	body["Content-Type"] = "text/html"
	body["Content-Transfer-Encoding"] = "8bit"

	mBody := ""
	for k, v := range body {
		mBody += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	mBody += fmt.Sprintf("\r\n%s\r\n--%s", b.Content, marker)

	return mBody
}
