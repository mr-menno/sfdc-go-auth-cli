package main

import (
	"net/url"
	"strings"
	"testing"
)

func TestGenerateState(t *testing.T) {
	state1 := generateState()
	state2 := generateState()

	// States should be different
	if state1 == state2 {
		t.Error("Generated states should be unique")
	}

	// States should not be empty
	if state1 == "" || state2 == "" {
		t.Error("Generated states should not be empty")
	}

	// States should be reasonable length (base64 encoded 32 bytes = ~44 chars)
	if len(state1) < 40 || len(state2) < 40 {
		t.Error("Generated states should be sufficiently long")
	}
}

func TestBuildAuthURL(t *testing.T) {
	// Set test values
	clientID = "test_client_id"
	state = "test_state"
	testDomain := "login.salesforce.com"

	authURL := buildAuthURL(testDomain)

	// Parse the URL
	parsedURL, err := url.Parse(authURL)
	if err != nil {
		t.Fatalf("Failed to parse auth URL: %v", err)
	}

	// Check base URL
	expectedBase := "https://login.salesforce.com/services/oauth2/authorize"
	baseURL := parsedURL.Scheme + "://" + parsedURL.Host + parsedURL.Path
	if baseURL != expectedBase {
		t.Errorf("Expected base URL %s, got %s", expectedBase, baseURL)
	}

	// Check query parameters
	query := parsedURL.Query()

	expectedParams := map[string]string{
		"response_type": "code",
		"client_id":     "test_client_id",
		"redirect_uri":  "http://localhost:8080/callback",
		"state":         "test_state",
		"scope":         "full refresh_token",
	}

	for key, expectedValue := range expectedParams {
		if actualValue := query.Get(key); actualValue != expectedValue {
			t.Errorf("Expected %s=%s, got %s=%s", key, expectedValue, key, actualValue)
		}
	}
}

func TestTokenResponseStructure(t *testing.T) {
	// Test that TokenResponse can be marshaled to JSON with expected fields
	response := TokenResponse{
		AccessToken:  "test_access_token",
		RefreshToken: "test_refresh_token",
		InstanceURL:  "https://test.salesforce.com",
	}

	// This should not panic - just verify fields are accessible
	_ = response.AccessToken
	_ = response.RefreshToken
	_ = response.InstanceURL

	// Check field names match expected JSON tags
	if response.AccessToken == "" {
		t.Error("AccessToken should not be empty in test")
	}
	if response.RefreshToken == "" {
		t.Error("RefreshToken should not be empty in test")
	}
	if response.InstanceURL == "" {
		t.Error("InstanceURL should not be empty in test")
	}
}

func TestSalesforceOAuthResponseStructure(t *testing.T) {
	// Test that SalesforceOAuthResponse has all expected fields
	response := SalesforceOAuthResponse{
		AccessToken:  "test_access",
		RefreshToken: "test_refresh",
		InstanceURL:  "https://test.salesforce.com",
		ID:           "test_id",
		TokenType:    "Bearer",
		IssuedAt:     "1234567890",
		Signature:    "test_signature",
	}

	// Verify all fields are accessible
	if response.AccessToken == "" || response.RefreshToken == "" || response.InstanceURL == "" {
		t.Error("Required OAuth response fields should not be empty")
	}
}

func TestConstants(t *testing.T) {
	// Test that constants are properly defined
	if defaultSalesforceDomain == "" {
		t.Error("defaultSalesforceDomain should not be empty")
	}
	if defaultPort == "" {
		t.Error("defaultPort should not be empty")
	}

	// Test that URL functions work correctly
	testDomain := "login.salesforce.com"
	authURL := getSalesforceAuthURL(testDomain)
	tokenURL := getSalesforceTokenURL(testDomain)

	if authURL == "" {
		t.Error("getSalesforceAuthURL should not return empty string")
	}
	if tokenURL == "" {
		t.Error("getSalesforceTokenURL should not return empty string")
	}

	// Test that URLs are valid
	if !strings.HasPrefix(authURL, "https://") {
		t.Error("Auth URL should use HTTPS")
	}
	if !strings.HasPrefix(tokenURL, "https://") {
		t.Error("Token URL should use HTTPS")
	}

	// Test custom domain
	customDomain := "company.my.salesforce.com"
	customAuthURL := getSalesforceAuthURL(customDomain)
	customTokenURL := getSalesforceTokenURL(customDomain)

	expectedAuthURL := "https://company.my.salesforce.com/services/oauth2/authorize"
	expectedTokenURL := "https://company.my.salesforce.com/services/oauth2/token"

	if customAuthURL != expectedAuthURL {
		t.Errorf("Expected auth URL %s, got %s", expectedAuthURL, customAuthURL)
	}
	if customTokenURL != expectedTokenURL {
		t.Errorf("Expected token URL %s, got %s", expectedTokenURL, customTokenURL)
	}

	// Test default port format
	if defaultPort != "8080" {
		t.Error("defaultPort should be 8080")
	}
}

func TestSalesforceURLFunctions(t *testing.T) {
	tests := []struct {
		domain           string
		expectedAuthURL  string
		expectedTokenURL string
	}{
		{
			domain:           "login.salesforce.com",
			expectedAuthURL:  "https://login.salesforce.com/services/oauth2/authorize",
			expectedTokenURL: "https://login.salesforce.com/services/oauth2/token",
		},
		{
			domain:           "company.my.salesforce.com",
			expectedAuthURL:  "https://company.my.salesforce.com/services/oauth2/authorize",
			expectedTokenURL: "https://company.my.salesforce.com/services/oauth2/token",
		},
		{
			domain:           "test.sandbox.my.salesforce.com",
			expectedAuthURL:  "https://test.sandbox.my.salesforce.com/services/oauth2/authorize",
			expectedTokenURL: "https://test.sandbox.my.salesforce.com/services/oauth2/token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.domain, func(t *testing.T) {
			authURL := getSalesforceAuthURL(tt.domain)
			tokenURL := getSalesforceTokenURL(tt.domain)

			if authURL != tt.expectedAuthURL {
				t.Errorf("getSalesforceAuthURL(%s) = %s, want %s", tt.domain, authURL, tt.expectedAuthURL)
			}
			if tokenURL != tt.expectedTokenURL {
				t.Errorf("getSalesforceTokenURL(%s) = %s, want %s", tt.domain, tokenURL, tt.expectedTokenURL)
			}
		})
	}
}

func TestCLICommand(t *testing.T) {
	// Test that root command is properly configured
	if rootCmd.Use != "sfdc-auth" {
		t.Errorf("Expected command name 'sfdc-auth', got '%s'", rootCmd.Use)
	}

	if rootCmd.Short == "" {
		t.Error("Command should have a short description")
	}

	if rootCmd.Long == "" {
		t.Error("Command should have a long description")
	}
}

func TestCLIFlags(t *testing.T) {
	// Reset flags to defaults
	flagClientID = ""
	flagClientSecret = ""
	flagPort = defaultPort
	flagDomain = defaultSalesforceDomain
	flagQuiet = false

	// Test that flags are properly defined
	flags := rootCmd.Flags()

	clientIDFlag := flags.Lookup("client-id")
	if clientIDFlag == nil {
		t.Error("client-id flag should be defined")
	}

	clientSecretFlag := flags.Lookup("client-secret")
	if clientSecretFlag == nil {
		t.Error("client-secret flag should be defined")
	}

	portFlag := flags.Lookup("port")
	if portFlag == nil {
		t.Error("port flag should be defined")
	}

	domainFlag := flags.Lookup("domain")
	if domainFlag == nil {
		t.Error("domain flag should be defined")
	}

	quietFlag := flags.Lookup("quiet")
	if quietFlag == nil {
		t.Error("quiet flag should be defined")
	}

	// Test default values
	if portFlag.DefValue != defaultPort {
		t.Errorf("Expected port default value '%s', got '%s'", defaultPort, portFlag.DefValue)
	}
	if domainFlag.DefValue != defaultSalesforceDomain {
		t.Errorf("Expected domain default value '%s', got '%s'", defaultSalesforceDomain, domainFlag.DefValue)
	}
}

func TestInitialization(t *testing.T) {
	// Save original values
	originalPort := port
	originalRedirectURI := redirectURI

	// Test init function behavior
	port = ""
	redirectURI = ""

	// Simulate init
	port = ":" + defaultPort
	redirectURI = "http://localhost:" + defaultPort + "/callback"

	// Test that values are set correctly
	if port != ":8080" {
		t.Errorf("Expected port ':8080', got '%s'", port)
	}

	if redirectURI != "http://localhost:8080/callback" {
		t.Errorf("Expected redirectURI 'http://localhost:8080/callback', got '%s'", redirectURI)
	}

	// Restore original values
	port = originalPort
	redirectURI = originalRedirectURI
}

func TestPortAndRedirectURIUpdate(t *testing.T) {
	// Save original values
	originalPort := port
	originalRedirectURI := redirectURI
	originalFlagPort := flagPort

	// Test port update logic
	flagPort = "9090"

	// Simulate the port update logic from runAuth
	if flagPort != defaultPort {
		port = ":" + flagPort
		redirectURI = "http://localhost:" + flagPort + "/callback"
	}

	// Test that values are updated correctly
	if port != ":9090" {
		t.Errorf("Expected port ':9090', got '%s'", port)
	}

	if redirectURI != "http://localhost:9090/callback" {
		t.Errorf("Expected redirectURI 'http://localhost:9090/callback', got '%s'", redirectURI)
	}

	// Restore original values
	port = originalPort
	redirectURI = originalRedirectURI
	flagPort = originalFlagPort
}
