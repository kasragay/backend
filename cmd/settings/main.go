package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kasragay/backend/internal/ports"
	"github.com/kasragay/backend/internal/repository"
	"github.com/kasragay/backend/internal/utils"

	"syscall"

	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/term"
)

func CreateSuperUser() {
	cmd := flag.NewFlagSet("createsuperuser", flag.ExitOnError)
	name := cmd.String("name", "", "name of the superuser")
	username := cmd.String("username", "", "username of the superuser")
	password := cmd.String("password", "", "password of the superuser")
	if err := cmd.Parse(os.Args[2:]); err != nil {
		log.Fatalf("error parsing command line arguments: %v", err)
	}
	logger := utils.NewLogger()
	relRepo := repository.NewRelationalRepo(logger)

	reader := bufio.NewReader(os.Stdin)
	if *username == "" {
		*username = InputString(reader, "username")
	}
	if *password == "" {
		*password = InputPassword("password")
		if *password != "" {
			rPassword := InputPassword("repeated-password")
			if *password != rPassword {
				log.Fatal("passwords do not match.")
			}
		}

	}
	if *name == "" {
		*name = InputString(reader, "name")
	}

	req := &ports.AuthSignupPostRequest{
		Name:     *name,
		Username: *username,
		UserType: ports.AdminUserType,
		Password: *password,
		Token:    "12345",
		Key:      "XXXXXXXX",
	}
	if err := ports.Validate(context.Background(), logger, req); err != nil {
		logger.Fatal(context.Background(), err.Error())
	}

	resp, err := relRepo.CreateUser(context.Background(), req)
	if err != nil {
		log.Fatalf("error creating superuser: %v", err)
	}
	fmt.Println("superuser created successfully.")
	fmt.Println("Id:", resp.Id)
}

func DeleteSuperUser() {
	cmd := flag.NewFlagSet("deletesuperuser", flag.ExitOnError)
	username := cmd.String("username", "", "Username of the superuser")
	if err := cmd.Parse(os.Args[2:]); err != nil {
		log.Fatalf("error parsing command line arguments: %v", err)
	}

	logger := utils.NewLogger()
	relRepo := repository.NewRelationalRepo(logger)

	reader := bufio.NewReader(os.Stdin)
	if *username == "" {
		*username = InputString(reader, "username")
	}
	areYouSure := InputString(reader, "are you sure? [y/n]")
	if areYouSure == "n" {
		log.Fatal("aborted.")
	} else if areYouSure != "y" {
		log.Fatal("invalid input.")
	}
	user, _, err := relRepo.GetUserByUsername(context.Background(), *username, ports.AdminUserType)
	if err != nil {
		log.Fatalf("error deleting superuser: %v", err)
	}
	if err := relRepo.DeleteUserById(context.Background(), user.GetId(), ports.AdminUserType); err != nil {
		log.Fatalf("error deleting superuser: %v", err)
	}
	fmt.Println("superuser deleted successfully.")
	fmt.Println("Id:", user.GetId())
}

func main() {
	subcommands := map[string]func(){
		"createsuperuser": CreateSuperUser,
		"deletesuperuser": DeleteSuperUser,
	}

	if len(os.Args) < 2 {
		errString := "expected "
		for k := range subcommands {
			errString += fmt.Sprintf("'%s' or ", k)
		}
		errString = errString[:len(errString)-4] + " subcommands."
		fmt.Println(errString)
		os.Exit(1)
	}
	if _, ok := subcommands[os.Args[1]]; !ok {
		fmt.Printf("unknown subcommand '%s'\n", os.Args[1])
		os.Exit(1)
	}
	subcommands[os.Args[1]]()
}

func InputString(reader *bufio.Reader, prompt string) string {
	fmt.Print("Enter " + prompt + ": ")
	text, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("error reading input: %v", err)
	}
	if len(text) == 1 {
		log.Fatalf("empty input for %s.", prompt)
	}
	return text[:len(text)-1]
}

func InputPassword(prompt string) string {
	fmt.Print("Enter " + prompt + ": ")
	password, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatalf("\nError reading %s: %v", prompt, err)
	}
	fmt.Println()
	return string(password)
}
