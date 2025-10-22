package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "projgen",
	Short: "Project generator CLI with CI/CD templates",
	Long:  "projgen scaffolds new projects (Spring Boot, later Node/Go/React) with CI/CD boilerplate.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
