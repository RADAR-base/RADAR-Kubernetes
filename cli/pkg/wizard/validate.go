package wizard

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// validateHostname checks that s is a valid DNS hostname suitable for use as a
// server name (e.g. radar.myinstitution.org). It rejects bare TLDs, IP addresses,
// URLs with schemes, and labels with invalid characters.
func validateHostname(s string) error {
	if s == "" {
		return fmt.Errorf("server hostname is required")
	}
	if strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
		return fmt.Errorf("enter a hostname only, without http:// or https://")
	}
	if strings.ContainsAny(s, " \t/") {
		return fmt.Errorf("hostname must not contain spaces, tabs, or slashes")
	}
	if net.ParseIP(s) != nil {
		return fmt.Errorf("enter a domain name, not an IP address")
	}
	if len(s) > 253 {
		return fmt.Errorf("hostname too long (max 253 characters)")
	}
	labels := strings.Split(s, ".")
	if len(labels) < 2 {
		return fmt.Errorf("hostname must include a domain, e.g. radar.example.org")
	}
	for _, label := range labels {
		if label == "" {
			return fmt.Errorf("hostname contains an empty label (double dot or leading/trailing dot)")
		}
		if len(label) > 63 {
			return fmt.Errorf("hostname label %q is too long (max 63 characters)", label)
		}
		if strings.HasPrefix(label, "-") || strings.HasSuffix(label, "-") {
			return fmt.Errorf("hostname label %q must not start or end with a hyphen", label)
		}
		for _, ch := range label {
			if !isAlphaNumeric(ch) && ch != '-' {
				return fmt.Errorf("hostname label %q contains invalid character %q (only letters, digits, hyphens allowed)", label, ch)
			}
		}
	}
	return nil
}

// validateEmail checks that s looks like a valid email address.
func validateEmail(s string) error {
	if s == "" {
		return fmt.Errorf("email is required")
	}
	if strings.ContainsAny(s, " \t") {
		return fmt.Errorf("email must not contain spaces")
	}
	at := strings.LastIndex(s, "@")
	if at < 1 {
		return fmt.Errorf("email must contain an @ sign with a local part before it")
	}
	local := s[:at]
	domain := s[at+1:]
	if local == "" {
		return fmt.Errorf("email is missing the local part before @")
	}
	if domain == "" {
		return fmt.Errorf("email is missing the domain part after @")
	}
	if !strings.Contains(domain, ".") {
		return fmt.Errorf("email domain %q must contain at least one dot", domain)
	}
	parts := strings.Split(domain, ".")
	for _, part := range parts {
		if part == "" {
			return fmt.Errorf("email domain contains an empty label (double dot or trailing dot)")
		}
	}
	return nil
}

// validateConfluentURL checks that s is a valid bootstrap server in host:port format.
func validateConfluentURL(s string) error {
	if s == "" {
		return fmt.Errorf("Confluent bootstrap server URL is required")
	}
	if strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
		return fmt.Errorf("enter a host:port pair, not a URL (e.g. pkc-xxx.us-east-1.aws.confluent.cloud:9092)")
	}
	host, portStr, err := net.SplitHostPort(s)
	if err != nil {
		return fmt.Errorf("expected host:port format, e.g. pkc-xxx.us-east-1.aws.confluent.cloud:9092")
	}
	if host == "" {
		return fmt.Errorf("hostname part of the bootstrap server URL is empty")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("port %q is not a valid port number (1–65535)", portStr)
	}
	return nil
}

// validateS3Bucket checks AWS S3 bucket naming rules.
func validateS3Bucket(s string) error {
	if s == "" {
		return fmt.Errorf("S3 bucket name is required")
	}
	if len(s) < 3 {
		return fmt.Errorf("S3 bucket name must be at least 3 characters")
	}
	if len(s) > 63 {
		return fmt.Errorf("S3 bucket name must be 63 characters or fewer")
	}
	if strings.HasPrefix(s, "-") || strings.HasSuffix(s, "-") ||
		strings.HasPrefix(s, ".") || strings.HasSuffix(s, ".") {
		return fmt.Errorf("S3 bucket name must not start or end with a hyphen or dot")
	}
	if strings.Contains(s, "..") {
		return fmt.Errorf("S3 bucket name must not contain consecutive dots")
	}
	if net.ParseIP(s) != nil {
		return fmt.Errorf("S3 bucket name must not be formatted as an IP address")
	}
	for _, ch := range s {
		if !isLowerAlphaNumeric(ch) && ch != '-' && ch != '.' {
			if isUpper(ch) {
				return fmt.Errorf("S3 bucket name must be lowercase (found %q)", ch)
			}
			return fmt.Errorf("S3 bucket name contains invalid character %q (only lowercase letters, digits, hyphens, dots allowed)", ch)
		}
	}
	return nil
}

// validateS3Region checks that s looks like a valid AWS region identifier.
func validateS3Region(s string) error {
	if s == "" {
		return fmt.Errorf("S3 region is required")
	}
	if strings.ContainsAny(s, " \t") {
		return fmt.Errorf("S3 region must not contain spaces")
	}
	// Regions follow the pattern: area-direction-number, e.g. us-east-1, eu-west-2, ap-southeast-1
	parts := strings.Split(s, "-")
	if len(parts) < 3 {
		return fmt.Errorf("S3 region %q doesn't look valid — expected format like us-east-1 or eu-west-2", s)
	}
	for _, p := range parts {
		if p == "" {
			return fmt.Errorf("S3 region %q contains an empty segment", s)
		}
	}
	return nil
}

// validateS3Key checks that an S3 access key ID is non-empty.
func validateS3Key(s string) error {
	if s == "" {
		return fmt.Errorf("S3 access key ID is required")
	}
	if strings.ContainsAny(s, " \t") {
		return fmt.Errorf("S3 access key ID must not contain spaces")
	}
	return nil
}

func isAlphaNumeric(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9')
}

func isLowerAlphaNumeric(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9')
}

func isUpper(ch rune) bool {
	return ch >= 'A' && ch <= 'Z'
}
