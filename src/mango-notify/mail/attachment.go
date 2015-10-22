package mail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mango-notify/models"
)

//Attachment in mail.
type Attachment struct {
	Lines models.Lines
}

//Build mail attachment to string format.
func (a *Attachment) Build() string {
	if a.Lines == nil || len(a.Lines) == 0 {
		return ""
	}

	//read and encode attachment
	var content []byte
	for _, l := range a.Lines {
		for _, b := range []byte(l.Content + "\n") {
			content = append(content, b)
		}
	}
	encoded := base64.StdEncoding.EncodeToString(content)

	//split the encoded file in lines (doesn't matter, but low enough not to hit a max limit)
	lineMaxLength := 500
	nbrLines := len(encoded) / lineMaxLength

	var buf bytes.Buffer
	//append lines to buffer
	for i := 0; i < nbrLines; i++ {
		buf.WriteString(encoded[i*lineMaxLength:(i+1)*lineMaxLength] + "\n")
	} //for

	//append last line in buffer
	buf.WriteString(encoded[nbrLines*lineMaxLength:])

	//part 3 will be the attachment
	headers := makeAttachmentHeaders()
	return fmt.Sprintf("\r\n%s\r\n\r\n%s\r\n--%s--", headers, buf.String(), marker)
}
