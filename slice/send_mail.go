package main

import (
    "log"
    "net/smtp"
)

var (
    // variable
    from ="xxx.coco@gmail.com"
    pass ="xxxx"
    to =[]string{"zhang.xxx@gmail.com"}
    msg =[]byte("To: zhang.xxx@gmail.com\r\n" +
                 "Subject: discount Gophers!\r\n" +
                 "\r\n" +
                 "This is the email body.\r\n")
)

func main() {
    // Set up authentication info
    smtp_srv := "smtp.gmail.com"
    auth :=smtp.PlainAuth("",from,pass,smtp_srv)
    err :=smtp.SendMail(smtp_srv+":587", auth, from,to,msg)
    if err !=nil {
           log.Fatal(err)
    }

}
