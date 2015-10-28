package mail

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"mango-notify/models"
)

//Attachment in mail.
type Attachment struct {
	Lines models.Lines
	//EncryptionKey should be base64 encoded string.
	EncryptionKey string
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
	var encoded string
	// if encryptonkey exist then use it to encrypt message.
	if a.EncryptionKey != "" {
		key, err := base64.URLEncoding.DecodeString(a.EncryptionKey)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Plaintext: ", content)
		fmt.Println("Key: ", key)
		ciphertext, err := encrypt(key, content)
		if err != nil {
			log.Fatal(err)
		}

		encoded = base64.URLEncoding.EncodeToString(ciphertext)
		//adding padding to make it work for android device..
		//if l := len(encoded) % 4; l > 0 {
		//	encoded += string([]byte{'=', '=', '='}[3-l:]) // or strings.Repeat("=", 4-l)
		//}
		fmt.Println("Encoded Ciphertext: ", encoded)
	} else {
		encoded = base64.URLEncoding.EncodeToString(content)
	}

	//part 3 will be the attachment
	headers := makeAttachmentHeaders()
	return fmt.Sprintf("\r\n%s\r\n\r\n%s\r\n--%s--", headers, encoded, marker)
}

func encrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	b := base64.URLEncoding.EncodeToString(text)

	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

func decrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.URLEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	return data, nil
}
