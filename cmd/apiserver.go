package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"goansible/work"
	"log"
	"strings"
)

var apiserverCmd = &cobra.Command{
	Use:   "apiserver",
	Short: "Apiserver the server",
	Long:  `Apiserver the server with a ansilbe route,default port 9090.`,
	Run: func(cmd *cobra.Command, args []string) {
		port := ":9090"
		if len(args) != 0 {
			port = args[0]
		}

		if !strings.HasPrefix(port, ":") {
			port = ":" + port
		}

		r := gin.Default()
		work.SetupRouter(r)
		err := r.Run(port)
		if err != nil {
			log.Println("apiserver err:", err)
			return
		}
	},
}

func init() {
	goansibleCmd.AddCommand(apiserverCmd)
}
