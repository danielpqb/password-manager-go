package main

import (
	"flag"
	password_manager "poc-go/internal/password"
)

func main() {
	cmd := flag.Arg(0)
	keyflag := flag.String("k", "", "Chave")
	valueflag := flag.String("v", "", "Valor")
	userflag := flag.String("u", "", "Usuário")
	passwordflag := flag.String("p", "", "Senha")

	flag.Parse()

	switch cmd {
	case "add":
		{
			password_manager.ValidatePassword(*userflag, *passwordflag)
			password_manager.Add(*keyflag, *valueflag)
		}
	case "get":
		{
			password_manager.ValidatePassword(*userflag, *passwordflag)
			password_manager.Get(*keyflag)
		}
	case "register":
		{
			password_manager.Register(*userflag, *passwordflag)
		}
	}
}
