package cmd

// EnvVarSpec represents an environment variable specification for the CLI.
type EnvVarSpec struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	RelatedFlag string `json:"related_flag,omitempty" yaml:"related_flag,omitempty"`
	Required    bool   `json:"required,omitempty" yaml:"required,omitempty"`
	Sensitive   bool   `json:"sensitive,omitempty" yaml:"sensitive,omitempty"`
}

// EnvVarRegistry contains all environment variables supported by f5xcctl.
// This registry is the single source of truth for --help and --spec output.
var EnvVarRegistry = []EnvVarSpec{
	{
		Name:        "VES_API_TOKEN",
		Description: "API token for authentication",
		RelatedFlag: "--api-token",
		Required:    false,
		Sensitive:   true,
	},
	{
		Name:        "VES_API_URL",
		Description: "API endpoint URL",
		RelatedFlag: "--server-url",
		Required:    false,
		Sensitive:   false,
	},
	{
		Name:        "VES_CACERT",
		Description: "CA certificate for TLS verification",
		RelatedFlag: "--cacert",
		Required:    false,
		Sensitive:   false,
	},
	{
		Name:        "VES_CERT",
		Description: "Client certificate for mTLS",
		RelatedFlag: "--cert",
		Required:    false,
		Sensitive:   false,
	},
	{
		Name:        "VES_CONFIG",
		Description: "Configuration file path",
		RelatedFlag: "--config",
		Required:    false,
		Sensitive:   false,
	},
	{
		Name:        "VES_KEY",
		Description: "Client private key for mTLS",
		RelatedFlag: "--key",
		Required:    false,
		Sensitive:   false,
	},
	{
		Name:        "VES_OUTPUT",
		Description: "Output format: text, json, yaml, table",
		RelatedFlag: "--output-format",
		Required:    false,
		Sensitive:   false,
	},
	{
		Name:        "VES_P12_FILE",
		Description: "PKCS#12 bundle file path",
		RelatedFlag: "--p12-bundle",
		Required:    false,
		Sensitive:   false,
	},
	{
		Name:        "VES_P12_PASSWORD",
		Description: "PKCS#12 bundle password",
		RelatedFlag: "",
		Required:    false,
		Sensitive:   true,
	},
}

// GetEnvVarByName returns the EnvVarSpec for a given environment variable name.
func GetEnvVarByName(name string) *EnvVarSpec {
	for i := range EnvVarRegistry {
		if EnvVarRegistry[i].Name == name {
			return &EnvVarRegistry[i]
		}
	}
	return nil
}

// GetEnvVarByFlag returns the EnvVarSpec for a given CLI flag.
func GetEnvVarByFlag(flag string) *EnvVarSpec {
	for i := range EnvVarRegistry {
		if EnvVarRegistry[i].RelatedFlag == flag {
			return &EnvVarRegistry[i]
		}
	}
	return nil
}
