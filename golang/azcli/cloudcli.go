package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Subscription struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	State  string `json:"state"`
	Tenant string `json:"tenantId"`
}

type Config struct {
	Azure struct {
		Tenants []string `json:"tenants"`
	} `json:"azure"`
	GCP struct {
		Projects []string `json:"projects"`
	} `json:"gcp"`
}

type Cluster struct {
	Name          string `json:"name"`
	Location      string `json:"location"`
	ResourceGroup string `json:"resourceGroup"`
	PowerState    struct {
		Code string `json:"code"`
	} `json:"powerState"`
}

type GCPCluster struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Status   string `json:"status"`
}

func main() {
	// Define subcommands
	loginCmd := flag.NewFlagSet("login", flag.ExitOnError)

	// Parse command line arguments
	if len(os.Args) < 2 {
		fmt.Println("cloudcli - A CLI tool for Azure authentication")
		fmt.Println("\nUsage:")
		fmt.Println("  cloudcli <command>")
		fmt.Println("\nCommands:")
		fmt.Println("  login         Login to Azure and select tenant/subscription")
		fmt.Println("\nExamples:")
		fmt.Println("  cloudcli login")
		os.Exit(1)
	}

	// Check which subcommand is used
	switch os.Args[1] {
	case "login":
		loginCmd.Parse(os.Args[2:])
		handleLogin()
	case "help", "--help", "-h":
		fmt.Println("cloudcli - A CLI tool for Azure authentication")
		fmt.Println("\nUsage:")
		fmt.Println("  cloudcli <command>")
		fmt.Println("\nCommands:")
		fmt.Println("  login         Login to Azure and select tenant/subscription")
		fmt.Println("\nExamples:")
		fmt.Println("  cloudcli login")
		os.Exit(0)
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		fmt.Println("\nAvailable commands:")
		fmt.Println("  login         Login to Azure and select tenant/subscription")
		fmt.Println("\nUse 'cloudcli help' for more information")
		os.Exit(1)
	}
}

func handleLogin() {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		os.Exit(1)
	}

	// Path to the tenants file
	tenantFile := filepath.Join(homeDir, "tenants.json")

	// Open and read the JSON file
	file, err := os.ReadFile(tenantFile)
	if err != nil {
		fmt.Printf("Error: Tenants file not found at %s\n", tenantFile)
		fmt.Println("\nTo fix this:")
		fmt.Println("1. Create a file named 'tenants.json' in your home directory")
		fmt.Println("2. Add your tenants in JSON format")
		fmt.Println("3. Example format:")
		fmt.Println(`   {
     "azure": {
       "tenants": [
         "12345678-1234-1234-1234-123456789012",
         "87654321-4321-4321-4321-210987654321"
       ]
     },
     "gcp": {
       "projects": [
         "project-1",
         "project-2"
       ]
     }
   }`)
		fmt.Println("\nYou can create the file using:")
		fmt.Printf("   echo '%s' > %s\n", `{"azure":{"tenants":["your-tenant-id"]},"gcp":{"projects":["your-project"]}}`, tenantFile)
		os.Exit(1)
	}

	// Parse the JSON file
	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		fmt.Printf("Error parsing tenants file: %v\n", err)
		fmt.Printf("Raw file content: %s\n", file)
		os.Exit(1)
	}

	// Display Azure tenants
	if len(config.Azure.Tenants) > 0 {
		fmt.Println("\nAvailable Azure tenants:")
		for i, tenantID := range config.Azure.Tenants {
			fmt.Printf("%d. %s\n", i+1, tenantID)
		}
	}

	// Display GCP projects
	if len(config.GCP.Projects) > 0 {
		fmt.Println("\nAvailable GCP projects:")
		for i, project := range config.GCP.Projects {
			fmt.Printf("%d. %s\n", len(config.Azure.Tenants)+i+1, project)
		}
	}

	if len(config.Azure.Tenants) == 0 && len(config.GCP.Projects) == 0 {
		fmt.Println("No tenants or projects found in the file")
		os.Exit(1)
	}

	// Get user selection
	fmt.Print("\nEnter number to select (1-", len(config.Azure.Tenants)+len(config.GCP.Projects), "): ")
	var input string
	fmt.Scanln(&input)

	idx := 0
	_, err = fmt.Sscanf(input, "%d", &idx)
	if err != nil || idx < 1 || idx > len(config.Azure.Tenants)+len(config.GCP.Projects) {
		fmt.Printf("Invalid selection. Please enter a number between 1 and %d\n", len(config.Azure.Tenants)+len(config.GCP.Projects))
		os.Exit(1)
	}

	// Check if selection is Azure tenant or GCP project
	if idx <= len(config.Azure.Tenants) {
		// Azure tenant selected
		selectedTenantID := config.Azure.Tenants[idx-1]
		fmt.Printf("\nLogging in to Azure tenant: %s\n", selectedTenantID)

		// Execute az login command
		cmd := exec.Command("az", "login", "-t", selectedTenantID)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error logging in to tenant: %v\n", err)
			fmt.Printf("Command output: %s\n", output)
			os.Exit(1)
		}

		fmt.Printf("Successfully logged in to Azure tenant: %s\n", selectedTenantID)

		// Get list of subscriptions
		cmd = exec.Command("az", "account", "list", "--output", "json")
		output, err = cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error getting subscription list: %v\n", err)
			fmt.Printf("Command output: %s\n", output)
			os.Exit(1)
		}

		// Parse the JSON output
		var subscriptions []Subscription
		if err := json.Unmarshal(output, &subscriptions); err != nil {
			fmt.Printf("Error parsing subscription list: %v\n", err)
			fmt.Printf("Raw output: %s\n", output)
			os.Exit(1)
		}

		if len(subscriptions) == 0 {
			fmt.Println("No subscriptions found for this tenant")
			os.Exit(0)
		}

		// Display available subscriptions
		fmt.Println("\nAvailable subscriptions:")
		for i, sub := range subscriptions {
			fmt.Printf("%d. %s (%s)\n", i+1, sub.Name, sub.ID)
		}

		// Get subscription selection
		fmt.Print("\nEnter subscription number to select (1-", len(subscriptions), "): ")
		fmt.Scanln(&input)

		subIdx := 0
		_, err = fmt.Sscanf(input, "%d", &subIdx)
		if err != nil || subIdx < 1 || subIdx > len(subscriptions) {
			fmt.Printf("Invalid selection. Please enter a number between 1 and %d\n", len(subscriptions))
			os.Exit(1)
		}

		// Set the selected subscription
		selectedSub := subscriptions[subIdx-1]
		cmd = exec.Command("az", "account", "set", "--subscription", selectedSub.ID)
		output, err = cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error setting subscription: %v\n", err)
			fmt.Printf("Command output: %s\n", output)
			os.Exit(1)
		}

		fmt.Printf("Successfully set subscription to: %s\n", selectedSub.Name)

		// List AKS clusters
		listAKSClusters()
	} else {
		// GCP project selected
		gcpIdx := idx - len(config.Azure.Tenants) - 1
		selectedProject := config.GCP.Projects[gcpIdx]
		fmt.Printf("\nSetting GCP project: %s\n", selectedProject)

		// Execute gcloud config set project command
		cmd := exec.Command("gcloud", "config", "set", "project", selectedProject)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error setting GCP project: %v\n", err)
			fmt.Printf("Command output: %s\n", output)
			os.Exit(1)
		}

		fmt.Printf("Successfully set GCP project to: %s\n", selectedProject)

		// List GKE clusters
		listGKEClusters()
	}
}

func getKubeconfig(clusterName, resourceGroup string) error {
	cmd := exec.Command("az", "aks", "get-credentials",
		"--resource-group", resourceGroup,
		"--name", clusterName,
		"--overwrite-existing")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error getting kubeconfig: %v\n", err)
		fmt.Printf("Command output: %s\n", output)
		return err
	}
	fmt.Printf("Successfully connected to cluster: %s\n", clusterName)
	return nil
}

func listAKSClusters() {
	// Get list of clusters
	cmd := exec.Command("az", "aks", "list", "--output", "json")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error getting cluster list: %v\n", err)
		fmt.Printf("Command output: %s\n", output)
		return
	}

	var clusters []Cluster
	if err := json.Unmarshal(output, &clusters); err != nil {
		fmt.Printf("Error parsing cluster list: %v\n", err)
		fmt.Printf("Raw output: %s\n", output)
		return
	}

	if len(clusters) == 0 {
		fmt.Println("No Kubernetes clusters found in the current subscription")
		return
	}

	// Display clusters
	fmt.Println("\nKubernetes clusters in current subscription:")
	fmt.Println("-------------------------------------------")
	for i, cluster := range clusters {
		fmt.Printf("%d. Name: %s\n", i+1, cluster.Name)
		fmt.Printf("   Location: %s\n", cluster.Location)
		fmt.Printf("   Resource Group: %s\n", cluster.ResourceGroup)
		fmt.Printf("   Power State: %s\n", cluster.PowerState.Code)
		fmt.Println("-------------------------------------------")
	}

	// Ask if user wants to connect to a cluster
	fmt.Print("\nDo you want to connect to a cluster? (y/n): ")
	var input string
	fmt.Scanln(&input)

	if input == "y" || input == "Y" {
		fmt.Print("Enter cluster number to connect (1-", len(clusters), "): ")
		fmt.Scanln(&input)

		idx := 0
		_, err = fmt.Sscanf(input, "%d", &idx)
		if err != nil || idx < 1 || idx > len(clusters) {
			fmt.Printf("Invalid selection. Please enter a number between 1 and %d\n", len(clusters))
			return
		}

		selectedCluster := clusters[idx-1]
		if err := getKubeconfig(selectedCluster.Name, selectedCluster.ResourceGroup); err != nil {
			fmt.Printf("Failed to connect to cluster: %v\n", err)
		}
	}
}

func getGCPKubeconfig(clusterName, location string) error {
	cmd := exec.Command("gcloud", "container", "clusters", "get-credentials",
		clusterName,
		"--region", location)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error getting kubeconfig: %v\n", err)
		fmt.Printf("Command output: %s\n", output)
		return err
	}
	fmt.Printf("Successfully connected to cluster: %s\n", clusterName)
	return nil
}

func listGKEClusters() {
	// Get list of clusters
	cmd := exec.Command("gcloud", "container", "clusters", "list", "--format", "json")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error getting cluster list: %v\n", err)
		fmt.Printf("Command output: %s\n", output)
		return
	}

	var clusters []GCPCluster
	if err := json.Unmarshal(output, &clusters); err != nil {
		fmt.Printf("Error parsing cluster list: %v\n", err)
		fmt.Printf("Raw output: %s\n", output)
		return
	}

	if len(clusters) == 0 {
		fmt.Println("No Kubernetes clusters found in the current project")
		return
	}

	// Display clusters
	fmt.Println("\nKubernetes clusters in current project:")
	fmt.Println("-------------------------------------------")
	for i, cluster := range clusters {
		fmt.Printf("%d. Name: %s\n", i+1, cluster.Name)
		fmt.Printf("   Location: %s\n", cluster.Location)
		fmt.Printf("   Status: %s\n", cluster.Status)
		fmt.Println("-------------------------------------------")
	}

	// Ask if user wants to connect to a cluster
	fmt.Print("\nDo you want to connect to a cluster? (y/n): ")
	var input string
	fmt.Scanln(&input)

	if input == "y" || input == "Y" {
		fmt.Print("Enter cluster number to connect (1-", len(clusters), "): ")
		fmt.Scanln(&input)

		idx := 0
		_, err = fmt.Sscanf(input, "%d", &idx)
		if err != nil || idx < 1 || idx > len(clusters) {
			fmt.Printf("Invalid selection. Please enter a number between 1 and %d\n", len(clusters))
			return
		}

		selectedCluster := clusters[idx-1]
		if err := getGCPKubeconfig(selectedCluster.Name, selectedCluster.Location); err != nil {
			fmt.Printf("Failed to connect to cluster: %v\n", err)
		}
	}
}
