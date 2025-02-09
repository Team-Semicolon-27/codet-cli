package main

import (
	"codet-cli/funcs"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func CloneFunction(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Error: Missing codat link.\nUsage: codat clone <codatlink>")
		return
	}
	coddatlink := args[0]
	fmt.Println("Cloning codat from:", coddatlink)
	funcs.Clone(coddatlink)
}

func InitFunction(cmd *cobra.Command, args []string) {
	fmt.Println("Initializing codat repository...")
	funcs.Init()
	fmt.Println("Initialization complete!")
}

func SetOrigin(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Error: Missing codat link.\nUsage: codat set-origin <codatlink>")
		return
	}
	coddatlink := args[0]
	fmt.Println("Setting codat origin to:", coddatlink)
	funcs.SetOrigin(coddatlink)
	fmt.Println("Origin set successfully!")
}

func SetToken(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Error: Missing token.\nUsage: codat set-token <token>")
		return
	}
	token := args[0]
	fmt.Println("Setting authentication token...")
	funcs.SetToken(token)
	fmt.Println("Token set successfully!")
}

func Push(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Error: Missing file name.\nUsage: codat push <filename>")
		return
	}

	file := args[0]
	if _, err := os.Stat(file); os.IsNotExist(err) {
		fmt.Println("Error: File does not exist:", file)
		return
	}

	funcs.Push(file)
	fmt.Println("File pushed successfully!")
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "codet",
		Short: "A CLI tool for codet",
	}

	var cloneCmd = &cobra.Command{
		Use:   "clone <codatlink>",
		Short: "Clone a codat",
		Run:   CloneFunction,
	}
	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize a codat repo",
		Run:   InitFunction,
	}
	var setOriginCmd = &cobra.Command{
		Use:   "set-origin <codatlink>",
		Short: "Set the remote codat origin",
		Run:   SetOrigin,
	}
	var setTokenCmd = &cobra.Command{
		Use:   "set-token <token>",
		Short: "Set the codat config token",
		Run:   SetToken,
	}
	var pushCmd = &cobra.Command{
		Use:   "push <fileName>",
		Short: "Push the codat",
		Run:   Push,
	}

	rootCmd.AddCommand(cloneCmd, initCmd, setOriginCmd, setTokenCmd, pushCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error executing command:", err)
		os.Exit(1)
	}
}
