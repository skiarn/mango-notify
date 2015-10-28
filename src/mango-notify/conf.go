package main

import (
	"encoding/base64"
	"flag"
	"log"
)

//Conf is flag information configuration nessesary to run the application.
type Conf struct {
	File          string
	From          string
	Pwd           string
	To            string
	ServerName    string
	EncryptionKey string
}

//GetConf reads nessesary flags and validates them.
func GetConf() *Conf {
	file := flag.String("file", "/var/log/auth.log", "Define auth file to watch.")
	from := flag.String("mailfrom", "", "Define user who send the mail.")
	pwd := flag.String("pwd", "", "Define user pwd for sending mail account.")
	to := flag.String("mailto", "", "Define user who receive the mail.")
	servername := flag.String("server", "", "Define mail server name and port. ex: smtp.gmail.com:465 (OBS only tsl supported.)")
	encryptonKey := flag.String("encKey", "", "Base64 encoded key, needs to be 32 bytes decoded. Used encryption is AES iv CFB.")
	flag.Parse()

	if *from == "" {
		log.Fatal("mailfrom is required.")
	}
	if *pwd == "" {
		log.Fatal("pwd is required.")
	}
	if *to == "" {
		log.Fatal("mailto is required.")
	}
	if *servername == "" {
		log.Fatal("server is required. ex: smtp.gmail.com:465")
	}

	if *encryptonKey == "" {
		log.Fatal("encKey is required.")
	}
	if *encryptonKey != "" {
		key, err := base64.URLEncoding.DecodeString(*encryptonKey)
		if err != nil {
			log.Fatal(err)
		}
		if len(key) != 32 {
			log.Fatal("encKey has to be 32 bytes decoded.")
		}
	}

	return &Conf{File: *file, From: *from, Pwd: *pwd, To: *to, ServerName: *servername, EncryptionKey: *encryptonKey}
}
