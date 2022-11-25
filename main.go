package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"log"
	"net/smtp"
	"strings"
	"time"
)

const (
	SubjectHeader   = "Subject"
	FromHeader      = "From"
	ToHeader        = "To"
	DateHeader      = "Date"
	MessageIDHeader = "Message-Id"
)

const (
	NewLine    = "\r\n"
	Domain     = "ztmail.net"
	MailDomain = "mail.ztmail.net"
)

func main() {
	var rcpt string
	var from string
	var mx string
	var subject string
	var body string

	flag.StringVar(&rcpt, "rcpt", "", "recipient")
	flag.StringVar(&from, "from", "", "from")
	flag.StringVar(&mx, "mx", "", "mx:port")
	flag.StringVar(&subject, "subject", "", "subject")
	flag.StringVar(&body, "body", "", "body")

	flag.Parse()

	client, err := smtp.Dial(mx)
	if err != nil {
		log.Fatal("Dial", err)
	}
	defer client.Close()

	err = client.Hello(MailDomain)
	if err != nil {
		log.Fatal("Hello", err)
	}

	err = client.Mail(from)
	if err != nil {
		log.Fatal("Mail", err)
	}

	err = client.Rcpt(rcpt)
	if err != nil {
		log.Fatal("Rcpt", err)
	}

	wc, err := client.Data()
	if err != nil {
		log.Fatal("Data", err)
	}

	wc.Write(header(SubjectHeader, subject))
	wc.Write(header(FromHeader, from))
	wc.Write(header(ToHeader, rcpt))
	wc.Write(header(DateHeader, time.Now().Format(time.RFC1123Z)))
	wc.Write(header(MessageIDHeader, messageID()+"@"+Domain))
	wc.Write([]byte(NewLine))
	wc.Write([]byte(body))
	wc.Write([]byte(NewLine))
	wc.Close()
}

func header(header string, value string) []byte {
	return []byte(header + ": " + value + NewLine)
}

func messageID() string {
	hasher := sha256.New()
	hasher.Write([]byte(Domain))
	hasher.Write([]byte(time.Now().Format(time.RFC1123Z)))
	binaryHash := hasher.Sum(nil)

	hash := hex.EncodeToString(binaryHash)

	return strings.ReplaceAll(hash, "=", "")
}
