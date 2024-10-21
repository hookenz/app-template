package main

import (
	"embed"
	"flag"
	"fmt"
	"os"
	"syscall"

	"github.com/hookenz/app-template/api/db"
	"github.com/hookenz/app-template/api/server"

	"golang.org/x/term"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//go:embed assets
var assets embed.FS

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	var addUserName string
	flag.StringVar(&addUserName, "add-user", "", "add a new user")
	flag.Parse()

	store := db.NewSqliteStore("database.db")

	err := store.Open()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	if addUserName != "" {
		err = PromptAddUser(store, addUserName)
		if err != nil {
			fmt.Printf("Error adding user: %v", err)
		}

		return
	}

	s := server.New(":9000", store, assets)
	s.Start()
}

func PromptAddUser(store db.Database, addUserName string) error {
	fmt.Println("Adding a new user to the database")
	fmt.Printf("Password: ")

	password, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("password cannot be empty")
	}

	return store.InsertUser(addUserName, string(password))
}
