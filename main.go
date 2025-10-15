package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// TokenResponse represents the JSON response structure
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	InstanceURL  string `json:"instance_url"`
}

// SalesforceOAuthResponse represents the OAuth response from Salesforce
type SalesforceOAuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	InstanceURL  string `json:"instance_url"`
	ID           string `json:"id"`
	TokenType    string `json:"token_type"`
	IssuedAt     string `json:"issued_at"`
	Signature    string `json:"signature"`
}

const (
	salesforceAuthURL  = "https://login.salesforce.com/services/oauth2/authorize"
	salesforceTokenURL = "https://login.salesforce.com/services/oauth2/token"
	defaultPort        = "8080"
)

var (
	authCode     string
	authError    string
	serverDone   = make(chan bool)
	clientID     string
	clientSecret string
	state        string
	redirectURI  string
	port         string

	// CLI flags
	flagClientID     string
	flagClientSecret string
	flagPort         string
	flagQuiet        bool
)

var rootCmd = &cobra.Command{
	Use:   "sfdc-auth",
	Short: "Salesforce OAuth2 Authentication CLI",
	Long: `A command-line tool that authenticates with Salesforce using OAuth2
and returns access tokens, refresh tokens, and instance URLs in JSON format.`,
	Run: runAuth,
}

func init() {
	// Initialize default values
	port = ":" + defaultPort
	redirectURI = "http://localhost:" + defaultPort + "/callback"

	rootCmd.Flags().StringVarP(&flagClientID, "client-id", "c", "", "Salesforce Client ID (Consumer Key)")
	rootCmd.Flags().StringVarP(&flagClientSecret, "client-secret", "s", "", "Salesforce Client Secret (Consumer Secret)")
	rootCmd.Flags().StringVarP(&flagPort, "port", "p", defaultPort, "Port for OAuth callback server")
	rootCmd.Flags().BoolVarP(&flagQuiet, "quiet", "q", false, "Suppress informational output")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func runAuth(cmd *cobra.Command, args []string) {
	if !flagQuiet {
		fmt.Println("Salesforce OAuth2 Authentication CLI")
		fmt.Println("====================================")
	}

	// Use flag values if provided, otherwise prompt
	clientID = flagClientID
	clientSecret = flagClientSecret

	if clientID == "" || clientSecret == "" {
		if err := getClientCredentials(); err != nil {
			log.Fatalf("Error getting client credentials: %v", err)
		}
	}

	// Update port if specified
	if flagPort != defaultPort {
		port = ":" + flagPort
		redirectURI = "http://localhost:" + flagPort + "/callback"
	}

	// Generate state parameter for security
	state = generateState()

	// Start local server for OAuth callback
	server := &http.Server{Addr: port}
	http.HandleFunc("/callback", handleCallback)

	go func() {
		if !flagQuiet {
			fmt.Printf("Starting local server on %s for OAuth callback...\n", port)
		}
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Build authorization URL
	authURL := buildAuthURL()
	if !flagQuiet {
		fmt.Printf("\nPlease open the following URL in your browser to authenticate:\n%s\n", authURL)
		fmt.Println("\nWaiting for OAuth callback...")
	}

	// Wait for callback
	<-serverDone

	// Shutdown server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	if authError != "" {
		log.Fatalf("OAuth error: %s", authError)
	}

	if authCode == "" {
		log.Fatal("No authorization code received")
	}

	// Exchange authorization code for tokens
	tokenResponse, err := exchangeCodeForTokens(authCode)
	if err != nil {
		log.Fatalf("Error exchanging code for tokens: %v", err)
	}

	// Output the result as JSON
	result := TokenResponse{
		AccessToken:  tokenResponse.AccessToken,
		RefreshToken: tokenResponse.RefreshToken,
		InstanceURL:  tokenResponse.InstanceURL,
	}

	jsonOutput, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	if !flagQuiet {
		fmt.Println("\nAuthentication successful!")
	}
	fmt.Println(string(jsonOutput))
}

func getClientCredentials() error {
	reader := bufio.NewReader(os.Stdin)

	// Get Client ID
	fmt.Print("Enter Salesforce Client ID: ")
	clientIDInput, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("error reading client ID: %v", err)
	}
	clientID = strings.TrimSpace(clientIDInput)

	if clientID == "" {
		return fmt.Errorf("client ID cannot be empty")
	}

	// Get Client Secret (hidden input)
	fmt.Print("Enter Salesforce Client Secret: ")
	clientSecretBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("error reading client secret: %v", err)
	}
	clientSecret = strings.TrimSpace(string(clientSecretBytes))
	fmt.Println() // New line after hidden input

	if clientSecret == "" {
		return fmt.Errorf("client secret cannot be empty")
	}

	return nil
}

func generateState() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		// Fallback to timestamp-based state if random generation fails
		fallback := make([]byte, 0, 32)
		fallback = fmt.Appendf(fallback, "state_%d", time.Now().UnixNano())
		return base64.URLEncoding.EncodeToString(fallback)
	}
	return base64.URLEncoding.EncodeToString(b)
}

func buildAuthURL() string {
	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", clientID)
	params.Add("redirect_uri", redirectURI)
	params.Add("state", state)
	params.Add("scope", "full refresh_token")

	return salesforceAuthURL + "?" + params.Encode()
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	defer func() {
		serverDone <- true
	}()

	// Check for error parameter
	if errorParam := r.URL.Query().Get("error"); errorParam != "" {
		authError = fmt.Sprintf("%s: %s", errorParam, r.URL.Query().Get("error_description"))
		http.Error(w, "OAuth error occurred. Check your terminal.", http.StatusBadRequest)
		return
	}

	// Verify state parameter
	receivedState := r.URL.Query().Get("state")
	if receivedState != state {
		authError = "Invalid state parameter"
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	// Get authorization code
	authCode = r.URL.Query().Get("code")
	if authCode == "" {
		authError = "No authorization code received"
		http.Error(w, "No authorization code received", http.StatusBadRequest)
		return
	}

	// Send success response
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(`
		<html>
			<body>
				<h2>Authentication Successful!</h2>
				<p>You can close this window and return to your terminal.</p>
			</body>
		</html>
	`)); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func exchangeCodeForTokens(code string) (*SalesforceOAuthResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("redirect_uri", redirectURI)
	data.Set("code", code)

	resp, err := http.PostForm(salesforceTokenURL, data)
	if err != nil {
		return nil, fmt.Errorf("error making token request: %v", err)
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token request failed with status: %d", resp.StatusCode)
	}

	var tokenResp SalesforceOAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("error decoding token response: %v", err)
	}

	return &tokenResp, nil
}
