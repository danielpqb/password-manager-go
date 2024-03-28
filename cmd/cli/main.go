package main

import (
	"flag"
	"os"
	database_manager "poc-go/internal/database"
	password_manager "poc-go/internal/password"
)

func main() {
	cmd := os.Args[len(os.Args)-1]
	var keyflag string
	flag.StringVar(&keyflag, "k", "", "Chave")
	var valueflag string
	flag.StringVar(&valueflag, "v", "", "Valor")
	var userflag string
	flag.StringVar(&userflag, "u", "", "Usuário")
	var passwordflag string
	flag.StringVar(&passwordflag, "p", "", "Senha")

	flag.Parse()

	db := database_manager.NewDatabase()
	db.ConnectDB()
	defer db.Db.Close()

	pw := password_manager.NewPassword(db.Db)

	// fmt.Println("CmdArg:", cmd)
	// fmt.Println("Flags:", keyflag, valueflag, userflag, passwordflag)

	switch cmd {
	case "add":
		{
			pw.ValidatePassword(userflag, passwordflag)
			pw.Add(keyflag, valueflag, passwordflag)
		}
	case "get":
		{
			pw.ValidatePassword(userflag, passwordflag)
			pw.Get(keyflag, passwordflag)
		}
	case "register":
		{
			pw.Register(userflag, passwordflag)
		}
	}

}
