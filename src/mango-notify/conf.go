package main

import (
	"flag"
	"log"
)

//Conf is flag information configuration nessesary to run the application.
type Conf struct {
	File       *string
	From       *string
	Pwd        *string
	To         *string
	ServerName *string
}

//GetConf reads nessesary flags and validates them.
func GetConf() *Conf {
	file := flag.String("file", "/var/log/auth.log", "Define auth file to watch.")
	from := flag.String("mailfrom", "", "Define user who send the mail.")
	pwd := flag.String("pwd", "", "Define user pwd for sending mail account.")
	to := flag.String("mailto", "", "Define user who receive the mail.")
	servername := flag.String("server", "", "Define mail server name and port. ex: smtp.gmail.com:465 (OBS only tsl supported.)")
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

	return &Conf{File: file, From: from, Pwd: pwd, To: to, ServerName: servername}
}
