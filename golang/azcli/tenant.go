package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		os.Exit(1)
	}

	// Path to the tenant IDs file
	tenantFile := filepath.Join(homeDir, "tenants")

	// Open the file
	file, err := os.Open(tenantFile)
	if err != nil {
		fmt.Printf("Error: Tenant file not found at %s\n", tenantFile)
		fmt.Println("\nTo fix this:")
		fmt.Println("1. Create a file named 'tenants' in your home directory")
		fmt.Println("2. Add your tenant IDs, one per line")
		fmt.Println("3. Example format:")
		fmt.Println("   12345678-1234-1234-1234-123456789012")
		fmt.Println("   87654321-4321-4321-4321-210987654321")
		fmt.Println("\nYou can create the file using:")
		fmt.Printf("   echo 'your-tenant-id' > %s\n", tenantFile)
		os.Exit(1)
	}
	defer file.Close()

	// Read all tenant IDs into a slice
	var tenants []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tenantID := scanner.Text()
		if tenantID != "" {
			tenants = append(tenants, tenantID)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading tenant file: %v\n", err)
		os.Exit(1)
	}

	if len(tenants) == 0 {
		fmt.Println("No tenants found in the file")
		os.Exit(1)
	}

	// Display available tenants
	fmt.Println("\nAvailable tenants:")
	for i, tenant := range tenants {
		fmt.Printf("%d. %s\n", i+1, tenant)
	}

	// Get user selection
	fmt.Print("\nEnter tenant number to login (1-", len(tenants), "): ")
	var input string
	fmt.Scanln(&input)

	var selectedTenants []string
	idx := 0
	_, err = fmt.Sscanf(input, "%d", &idx)
	if err != nil || idx < 1 || idx > len(tenants) {
		fmt.Printf("Invalid selection. Please enter a number between 1 and %d\n", len(tenants))
		os.Exit(1)
	}
	selectedTenants = append(selectedTenants, tenants[idx-1])

	if len(selectedTenants) == 0 {
		fmt.Println("No valid tenants selected")
		os.Exit(1)
	}

	// Process selected tenants
	for _, tenantID := range selectedTenants {
		fmt.Printf("\nLogging in to tenant: %s\n", tenantID)

		// Execute az login command
		cmd := exec.Command("az", "login", "-t", tenantID)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error logging in to tenant %s: %v\n", tenantID, err)
			fmt.Printf("Command output: %s\n", output)
			continue
		}

		fmt.Printf("Successfully logged in to tenant: %s\n", tenantID)
	}
}
