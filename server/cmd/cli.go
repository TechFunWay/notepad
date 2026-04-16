package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"notepad/config"
	"notepad/database"
	"notepad/model"
)

func ExecuteCLI(args []string) {
	if len(args) == 0 {
		printUsage()
		return
	}

	cfg := config.Load(0, "")
	if err := database.Init(cfg.DBPath(), appVersion); err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		return
	}

	switch args[0] {
	case "recover-admin":
		recoverAdmin()
	case "find-admin":
		findAdmin()
	case "list-users":
		listUsers()
	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Notepad CLI Tool")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  notepad                  Start the web server")
	fmt.Println("  notepad recover-admin    Reset admin password")
	fmt.Println("  notepad find-admin       Show admin username")
	fmt.Println("  notepad list-users       List all users")
}

func recoverAdmin() {
	admin, err := model.GetAdminUser()
	if err != nil {
		fmt.Println("Error: No admin user found")
		return
	}

	fmt.Printf("Admin username: %s\n", admin.Username)
	fmt.Print("Enter new password: ")

	reader := bufio.NewReader(os.Stdin)
	password, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading password: %v\n", err)
		return
	}
	password = strings.TrimSpace(password)

	if len(password) < 6 {
		fmt.Println("Error: Password must be at least 6 characters")
		return
	}

	if err := model.UpdatePassword(admin.ID, password); err != nil {
		fmt.Printf("Error resetting password: %v\n", err)
		return
	}

	fmt.Println("Admin password reset successfully!")
}

func findAdmin() {
	admin, err := model.GetAdminUser()
	if err != nil {
		fmt.Println("Error: No admin user found")
		return
	}
	fmt.Printf("Admin username: %s\n", admin.Username)
	fmt.Printf("Created at: %s\n", admin.CreatedAt.Format("2006-01-02 15:04:05"))
}

func listUsers() {
	users, err := model.ListUsers()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if len(users) == 0 {
		fmt.Println("No users found")
		return
	}

	fmt.Printf("%-5s %-20s %-10s %-20s\n", "ID", "Username", "Role", "Created At")
	fmt.Println(strings.Repeat("-", 60))
	for _, u := range users {
		fmt.Printf("%-5d %-20s %-10s %-20s\n",
			u.ID, u.Username, u.Role, u.CreatedAt.Format("2006-01-02 15:04:05"))
	}
}
