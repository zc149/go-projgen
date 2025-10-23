package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "projgen",
	Short: "Project generator CLI with CI/CD templates",
	Long:  "projgen scaffolds new projects (Spring Boot, later Node/Go/React) with CI/CD boilerplate.",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
