package config

import "strings"

type ValidationError struct {
	Field   string
	Message string
}

type ValidationWarning struct {
	Field   string
	Message string
}

type ValidationResult struct {
	Errors   []ValidationError
	Warnings []ValidationWarning
}

func (r *ValidationResult) IsValid() bool {
	return len(r.Errors) == 0
}

var placeholders = []string{"change_me", "MAINTAINER_EMAIL@example.com", "example.com"}

func isPlaceholder(v string) bool {
	for _, p := range placeholders {
		if strings.EqualFold(v, p) {
			return true
		}
	}
	return false
}

func Validate(cfg *Config, sec *Secrets) ValidationResult {
	var result ValidationResult

	if cfg.ServerName == "" || cfg.ServerName == "example.com" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "server_name",
			Message: "must be set to your actual domain name",
		})
	}

	if cfg.MaintainerEmail == "" || cfg.MaintainerEmail == "MAINTAINER_EMAIL@example.com" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "maintainer_email",
			Message: "must be set to a real email address",
		})
	}

	if cfg.ConfluentCloud && sec.ConfluentCloud.BootstrapServer == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "confluent_cloud.bootstrapServerurl",
			Message: "required when confluent_cloud is enabled",
		})
	}

	if cfg.EnableFitbit && sec.FitbitClientID == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "fitbit_client_id",
			Message: "required when Fitbit integration is enabled",
		})
	}

	secretFields := map[string]string{
		"fitbit_client_id":     sec.FitbitClientID,
		"fitbit_client_secret": sec.FitbitClientSecret,
		"garmin_consumer_key":  sec.GarminConsumerKey,
		"redcap_token":         sec.RedcapToken,
	}
	for field, val := range secretFields {
		if val != "" && isPlaceholder(val) {
			result.Warnings = append(result.Warnings, ValidationWarning{
				Field:   field,
				Message: "still contains a placeholder value",
			})
		}
	}

	return result
}
