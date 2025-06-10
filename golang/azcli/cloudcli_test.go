package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestConfigParsing(t *testing.T) {
	// Create a temporary test file
	testFile := filepath.Join("tests", "tenants.json")
	file, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Error reading test file: %v", err)
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		t.Fatalf("Error parsing config: %v", err)
	}

	// Test Azure tenants
	if len(config.Azure.Tenants) != 2 {
		t.Errorf("Expected 2 Azure tenants, got %d", len(config.Azure.Tenants))
	}
	expectedTenants := []string{
		"98765432-1234-5678-9012-345678901234",
		"87654321-4321-4321-4321-210987654321",
	}
	for i, tenant := range config.Azure.Tenants {
		if tenant != expectedTenants[i] {
			t.Errorf("Expected tenant %s, got %s", expectedTenants[i], tenant)
		}
	}

	// Test GCP projects
	if len(config.GCP.Projects) != 2 {
		t.Errorf("Expected 2 GCP projects, got %d", len(config.GCP.Projects))
	}
	expectedProjects := []string{
		"project-12345",
		"my-project-2",
	}
	for i, project := range config.GCP.Projects {
		if project != expectedProjects[i] {
			t.Errorf("Expected project %s, got %s", expectedProjects[i], project)
		}
	}
}

func TestClusterParsing(t *testing.T) {
	// Read test AKS cluster data
	testFile := filepath.Join("tests", "aks.json")
	file, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Error reading test file: %v", err)
	}

	var cluster Cluster
	if err := json.Unmarshal(file, &cluster); err != nil {
		t.Fatalf("Error parsing cluster: %v", err)
	}

	// Test cluster fields
	expectedName := "testestsst"
	if cluster.Name != expectedName {
		t.Errorf("Expected cluster name %s, got %s", expectedName, cluster.Name)
	}

	expectedLocation := "northeurope"
	if cluster.Location != expectedLocation {
		t.Errorf("Expected location %s, got %s", expectedLocation, cluster.Location)
	}

	expectedResourceGroup := "anoop-test"
	if cluster.ResourceGroup != expectedResourceGroup {
		t.Errorf("Expected resource group %s, got %s", expectedResourceGroup, cluster.ResourceGroup)
	}

	expectedPowerState := "Running"
	if cluster.PowerState.Code != expectedPowerState {
		t.Errorf("Expected power state %s, got %s", expectedPowerState, cluster.PowerState.Code)
	}
}

func TestSubscriptionParsing(t *testing.T) {
	// Test subscription JSON parsing
	subscriptionJSON := `{
		"id": "test-subscription-id",
		"name": "Test Subscription",
		"state": "Enabled",
		"tenantId": "test-tenant-id"
	}`

	var subscription Subscription
	if err := json.Unmarshal([]byte(subscriptionJSON), &subscription); err != nil {
		t.Fatalf("Error parsing subscription: %v", err)
	}

	// Test subscription fields
	expectedID := "test-subscription-id"
	if subscription.ID != expectedID {
		t.Errorf("Expected subscription ID %s, got %s", expectedID, subscription.ID)
	}

	expectedName := "Test Subscription"
	if subscription.Name != expectedName {
		t.Errorf("Expected subscription name %s, got %s", expectedName, subscription.Name)
	}

	expectedState := "Enabled"
	if subscription.State != expectedState {
		t.Errorf("Expected subscription state %s, got %s", expectedState, subscription.State)
	}

	expectedTenant := "test-tenant-id"
	if subscription.Tenant != expectedTenant {
		t.Errorf("Expected tenant ID %s, got %s", expectedTenant, subscription.Tenant)
	}
}
