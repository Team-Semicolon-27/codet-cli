package main

import (
	"codet-cli/funcs"
	"fmt"

	"github.com/spf13/cobra"
)

func CloneFunction(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: codat clone <coddatlink>")
		return
	}
	coddatlink := args[0]
	funcs.Clone(coddatlink)
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

	rootCmd.AddCommand(cloneCmd)
	rootCmd.Execute()
}
