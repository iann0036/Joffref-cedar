package cedar

type ValidationResult struct {
	Passed bool              `json:"passed"`
	Errors []ValidationError `json:"errors"`
}

type ValidationError struct {
	PolicyID string `json:"policyId"`
	Note     string `json:"note"`
}

const (
	ValidationModeStrict     = "Strict"
	ValidationModePermissive = "Permissive"
)
