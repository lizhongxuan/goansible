package main

import (
	"fmt"
	"os"
)

import (
	"github.com/spf13/cobra"
)

var goansibleCmd = &cobra.Command{
	Use:   "goansible",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Run: func(cmd *cobra.Command, args []string) { },
}

func main() {
	if err := goansibleCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
