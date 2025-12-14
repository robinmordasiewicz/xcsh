package cmd

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/robinmordasiewicz/f5xcctl/pkg/output"
)

var secretsCmd = &cobra.Command{
	Use:   "secrets",
	Short: "Execute commands for secret_management",
	Long: `Manage secrets in F5 Distributed Cloud using blindfold encryption.

The secrets command provides tools for encrypting, decrypting, and managing
secrets that can only be accessed by the F5 XC platform. This uses the
blindfold encryption mechanism for secure secret storage.`,
	Example: `  # Get the public key for encryption
  f5xcctl request secrets get-public-key

  # Encrypt a secret
  f5xcctl request secrets encrypt --policy-doc policy.json --public-key key.pem secret.txt

  # Build a Kubernetes secret bundle
  f5xcctl request secrets build-blindfold-bundle --name example-secret --data secret.txt`,
}

// Encrypt command
var encryptFlags struct {
	policyDocument string
	publicKey      string
	outfile        string
}

var encryptCmd = &cobra.Command{
	Use:   "encrypt [secret-file]",
	Short: "Encrypt secret",
	Long: `Encrypt a secret file using F5 XC blindfold encryption.

The encryption uses the public key obtained from the F5 XC API
and a policy document that defines the decryption policy.`,
	Example: `  # Encrypt a secret with policy and public key
  f5xcctl request secrets encrypt --policy-doc policy.json --public-key key.pem secret.txt

  # Encrypt and save to file
  f5xcctl request secrets encrypt --policy-doc policy.json --public-key key.pem --outfile encrypted.txt secret.txt`,
	Args: cobra.ExactArgs(1),
	RunE: runEncrypt,
}

// GetPublicKey command
var getPublicKeyFlags struct {
	keyVersion uint32
}

var getPublicKeyCmd = &cobra.Command{
	Use:   "get-public-key",
	Short: "Get Public Key",
	Long: `Retrieve the public key from F5 Distributed Cloud for encrypting secrets.

The public key is used with the blindfold encryption mechanism to
encrypt secrets that can only be decrypted by the F5 XC platform.`,
	Example: `  # Get the current public key
  f5xcctl request secrets get-public-key

  # Get a specific key version
  f5xcctl request secrets get-public-key --key-version 1`,
	RunE: runGetPublicKey,
}

// GetPolicyDocument command
var getPolicyDocumentCmd = &cobra.Command{
	Use:   "get-policy-document",
	Short: "Get Policy Document",
	Long: `Retrieve the policy document for secret encryption from F5 Distributed Cloud.

The policy document defines which services and conditions can decrypt
the encrypted secrets.`,
	Example: `  # Get the policy document
  f5xcctl request secrets get-policy-document`,
	RunE: runGetPolicyDocument,
}

// SecretInfo command
var secretInfoCmd = &cobra.Command{
	Use:   "secret-info [encrypted-secret-file]",
	Short: "Secret Info",
	Long: `Parse and display information about an encrypted secret.

This command reads an encrypted secret file and displays metadata
about the encryption, including the policy and key version used.`,
	Example: `  # Show info about an encrypted secret
  f5xcctl request secrets secret-info encrypted-secret.txt`,
	Args: cobra.ExactArgs(1),
	RunE: runSecretInfo,
}

// BuildBlindfoldBundle command
var buildBundleFlags struct {
	name      string
	namespace string
	dataFile  string
	outfile   string
}

var buildBlindfoldBundleCmd = &cobra.Command{
	Use:   "build-blindfold-bundle",
	Short: "Build blindfold secret bundle for k8s secret",
	Long: `Build a Kubernetes secret manifest with blindfold-encrypted data.

This command creates a Kubernetes Secret resource with the encrypted
secret data, ready to be applied to a cluster managed by F5 XC.`,
	Example: `  # Build a secret bundle
  f5xcctl request secrets build-blindfold-bundle --name example-secret --data secret.txt

  # Build with custom namespace
  f5xcctl request secrets build-blindfold-bundle --name example-secret --namespace production --data secret.txt

  # Build and save to file
  f5xcctl request secrets build-blindfold-bundle --name example-secret --data secret.txt --outfile k8s-secret.yaml`,
	RunE: runBuildBlindfoldBundle,
}

func init() {
	requestCmd.AddCommand(secretsCmd)

	// Enable AI-agent-friendly error handling for invalid subcommands
	secretsCmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return fmt.Errorf("unknown command %q for %q\n\nUsage: f5xcctl request secrets <action> [flags]\n\nAvailable actions:\n  encrypt, get-public-key, get-policy-document, secret-info, build-blindfold-bundle\n\nRun 'f5xcctl request secrets --help' for usage", args[0], cmd.CommandPath())
		}
		return cmd.Help()
	}
	secretsCmd.SuggestionsMinimumDistance = 2

	// Encrypt command
	encryptCmd.Flags().StringVar(&encryptFlags.policyDocument, "policy-document", "", "File containing policy document")
	encryptCmd.Flags().StringVar(&encryptFlags.publicKey, "public-key", "", "File containing public key")
	encryptCmd.Flags().StringVar(&encryptFlags.outfile, "outfile", "", "Output file for encrypted secret")
	secretsCmd.AddCommand(encryptCmd)

	// GetPublicKey command
	getPublicKeyCmd.Flags().Uint32Var(&getPublicKeyFlags.keyVersion, "key-version", 0, "Key version to fetch (0 for latest)")
	secretsCmd.AddCommand(getPublicKeyCmd)

	// GetPolicyDocument command
	secretsCmd.AddCommand(getPolicyDocumentCmd)

	// SecretInfo command
	secretsCmd.AddCommand(secretInfoCmd)

	// BuildBlindfoldBundle command
	buildBlindfoldBundleCmd.Flags().StringVar(&buildBundleFlags.name, "name", "", "Name of the Kubernetes secret")
	buildBlindfoldBundleCmd.Flags().StringVar(&buildBundleFlags.namespace, "namespace", "default", "Kubernetes namespace")
	buildBlindfoldBundleCmd.Flags().StringVar(&buildBundleFlags.dataFile, "data", "", "File containing secret data to encrypt")
	buildBlindfoldBundleCmd.Flags().StringVar(&buildBundleFlags.outfile, "outfile", "", "Output file for the bundle")
	_ = buildBlindfoldBundleCmd.MarkFlagRequired("name")
	_ = buildBlindfoldBundleCmd.MarkFlagRequired("data")
	secretsCmd.AddCommand(buildBlindfoldBundleCmd)
}

func runEncrypt(cmd *cobra.Command, args []string) error {
	// Read the secret file
	secretData, err := os.ReadFile(args[0])
	if err != nil {
		return fmt.Errorf("failed to read secret file: %w", err)
	}

	// Read public key
	if encryptFlags.publicKey == "" {
		return fmt.Errorf("--public-key is required")
	}
	pubKeyData, err := os.ReadFile(encryptFlags.publicKey)
	if err != nil {
		return fmt.Errorf("failed to read public key file: %w", err)
	}

	// Parse the public key
	block, _ := pem.Decode(pubKeyData)
	if block == nil {
		return fmt.Errorf("failed to parse PEM block containing the public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %w", err)
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("not an RSA public key")
	}

	// Read policy document if provided
	var policyDoc []byte
	if encryptFlags.policyDocument != "" {
		policyDoc, err = os.ReadFile(encryptFlags.policyDocument)
		if err != nil {
			return fmt.Errorf("failed to read policy document: %w", err)
		}
	}

	// Encrypt the secret using RSA-OAEP
	encrypted, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPub, secretData, nil)
	if err != nil {
		return fmt.Errorf("failed to encrypt secret: %w", err)
	}

	// Create the blindfold structure
	blindfold := map[string]interface{}{
		"type":     "blindfold",
		"location": "string:///",
		"secret_info": map[string]interface{}{
			"ciphertext":      base64.StdEncoding.EncodeToString(encrypted),
			"policy_document": string(policyDoc),
		},
	}

	// Output
	result, err := yaml.Marshal(blindfold)
	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}

	if encryptFlags.outfile != "" {
		if err := os.WriteFile(encryptFlags.outfile, result, 0600); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		output.PrintInfo(fmt.Sprintf("Encrypted secret written to %s", encryptFlags.outfile))
	} else {
		fmt.Println(string(result))
	}

	return nil
}

func runGetPublicKey(cmd *cobra.Command, args []string) error {
	client := GetClient()
	if client == nil {
		return fmt.Errorf("client not initialized - check configuration")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Build the API path
	path := "/api/web/secret_management/get_public_key"
	if getPublicKeyFlags.keyVersion > 0 {
		path = fmt.Sprintf("%s?key_version=%d", path, getPublicKeyFlags.keyVersion)
	}

	resp, err := client.Get(ctx, path, nil)
	if err != nil {
		return fmt.Errorf("failed to get public key: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API error: %s", string(resp.Body))
	}

	var result interface{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	return output.Print(result, GetOutputFormat())
}

func runGetPolicyDocument(cmd *cobra.Command, args []string) error {
	client := GetClient()
	if client == nil {
		return fmt.Errorf("client not initialized - check configuration")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	path := "/api/web/secret_management/get_policy_document"

	resp, err := client.Get(ctx, path, nil)
	if err != nil {
		return fmt.Errorf("failed to get policy document: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API error: %s", string(resp.Body))
	}

	var result interface{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	return output.Print(result, GetOutputFormat())
}

func runSecretInfo(cmd *cobra.Command, args []string) error {
	// Read the encrypted secret file
	data, err := os.ReadFile(args[0])
	if err != nil {
		return fmt.Errorf("failed to read encrypted secret file: %w", err)
	}

	// Parse as YAML first, then JSON
	var secretInfo map[string]interface{}
	if err := yaml.Unmarshal(data, &secretInfo); err != nil {
		if err := json.Unmarshal(data, &secretInfo); err != nil {
			return fmt.Errorf("failed to parse secret file (not valid YAML or JSON): %w", err)
		}
	}

	// Extract and display info
	info := map[string]interface{}{
		"type": secretInfo["type"],
	}

	if secretInfoData, ok := secretInfo["secret_info"].(map[string]interface{}); ok {
		if cipher, ok := secretInfoData["ciphertext"].(string); ok {
			info["ciphertext_length"] = len(cipher)
			info["ciphertext_preview"] = cipher[:min(50, len(cipher))] + "..."
		}
		if policy, ok := secretInfoData["policy_document"].(string); ok {
			info["has_policy_document"] = len(policy) > 0
		}
	}

	if location, ok := secretInfo["location"].(string); ok {
		info["location"] = location
	}

	return output.Print(info, GetOutputFormat())
}

func runBuildBlindfoldBundle(cmd *cobra.Command, args []string) error {
	// Read the data file
	data, err := os.ReadFile(buildBundleFlags.dataFile)
	if err != nil {
		return fmt.Errorf("failed to read data file: %w", err)
	}

	// Create Kubernetes Secret manifest
	secret := map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "Secret",
		"metadata": map[string]interface{}{
			"name":      buildBundleFlags.name,
			"namespace": buildBundleFlags.namespace,
			"annotations": map[string]string{
				"ves.io/blindfold": "true",
			},
		},
		"type": "Opaque",
		"data": map[string]string{
			"blindfold": base64.StdEncoding.EncodeToString(data),
		},
	}

	result, err := yaml.Marshal(secret)
	if err != nil {
		return fmt.Errorf("failed to marshal secret: %w", err)
	}

	if buildBundleFlags.outfile != "" {
		if err := os.WriteFile(buildBundleFlags.outfile, result, 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		output.PrintInfo(fmt.Sprintf("Kubernetes secret bundle written to %s", buildBundleFlags.outfile))
	} else {
		fmt.Println(string(result))
	}

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
