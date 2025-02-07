package main

import (
	"codet-cli/funcs"
	"fmt"

	"github.com/spf13/cobra"
)

func CloneFunction(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: codat clone <codatlink>")
		return
	}
	coddatlink := args[0]
	funcs.Clone(coddatlink)
}

func InitFunction(cmd *cobra.Command, args []string) {
	funcs.Init()
}

func SetOrigin(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: codat set-origin <cod	atlink>")
		return
	}
	coddatlink := args[0]
	funcs.SetOrigin(coddatlink)
}

func SetToken(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: codat set-origin <codatlink>")
		return
	}

	token := args[0]
	funcs.SetToken(token)
}

func Push(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: codat psuh <filename>")
		return
	}

	file := args[0]
	funcs.Push(file)
}

func main() {

	var rootCmd = &cobra.Command{
		Use:   "codet",
		Short: "A cli tool for codet",
	}

	var cloneCmd = &cobra.Command{
		Use:   "clone <codatlink>",
		Short: "clone a codat",
		Run:   CloneFunction,
	}
	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "initialize a codat repo",
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
	rootCmd.Execute()
}
