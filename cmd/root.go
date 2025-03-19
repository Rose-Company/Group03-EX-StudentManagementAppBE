// cmd/root.go
package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "student-management",
	Short: "Student Management Application Backend",
	Long:  "A backend service for managing student data, faculty, and authentication",
}

func Execute() {
	rootCmd.AddCommand(serverCmd)
	InitFlags()
	rootCmd.Execute()
}

func InitFlags() {
	serverCmd.PersistentFlags().Bool("start", false, "Start the server with default port 8080")
	serverCmd.PersistentFlags().String("port", "8080", "Port to run the server on")
	serverCmd.PersistentFlags().String("mode", "release", "Run mode (debug/release)")
}
