package main

import (
	"fmt"
	"mango-notify/io"
	"mango-notify/mail"
	"mango-notify/models"
	"sort"
	"time"
)

func main() {
	conf := GetConf()

	for {
		fileChangedChan := make(chan bool)
		io.OnFileChange(conf.File, fileChangedChan, conf.Time)
		<-fileChangedChan
		err := models.UpdateChangedLines(conf.File)

		if err != nil {
			fmt.Println("Error while parsing file:", err)
		} else {
			lines, err := models.GetUnsentLines()
			if err != nil {
				fmt.Println("Error getting unsent:", err)
			}
			sort.Sort(lines)
			t := time.Now()
			m := mail.Mail{
				Subject:    "Mango notify",
				Body:       mail.Body{Content: "Notify sent:" + t.Format("Mon Jan _2 15:04:05 2006")},
				Attachment: mail.Attachment{Lines: lines, EncryptionKey: conf.EncryptionKey},
				Conf: mail.Conf{From: conf.From,
					To:         conf.To,
					Password:   conf.Pwd,
					ServerName: conf.ServerName},
			}

			m.Send()

			err = models.SetLinesAsSent(&lines)
			if err != nil {
				fmt.Println("Error while saving sent lines:", err)
			}
		}

		//grep sshd.\*Accepted publickey for /var/log/auth.log | less ?
	}
}
