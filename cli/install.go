package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"web/models"
)

//Install 初始化
func Install() {
	if len(os.Args) > 2  && os.Args[1] == "install" {
		user := flag.String("user","admin","Administrator user.")
		pwd  := flag.String("password","","Administrator password")
		email := flag.String("email","","Administrator email.")

		flag.CommandLine.Parse(os.Args[2:])
		pasword := strings.TrimSpace(*pwd)

		if pasword == "" {
			fmt.Println("Administrator password  is required.")
			os.Exit(0)
		}

		if *email == "" {
			fmt.Println("Administrator email is required")
			os.Exit(0)
		}

		member := models.NewUser()
		member.Account = *user
		member.Password = pasword
		member.Email = *email

		if err := member.Add(); err != nil  {
			fmt.Println(err.Error())
			os.Exit(0)
		}
		fmt.Println("ok")
		os.Exit(0)
	}
}