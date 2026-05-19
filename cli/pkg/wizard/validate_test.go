package wizard

import "testing"

func TestValidateHostname(t *testing.T) {
	ok := []string{
		"radar.example.org",
		"my-radar.institution.ac.uk",
		"sub.domain.example.com",
		"xn--nxasmq6b.example.com",
	}
	bad := []string{
		"",
		"nodot",
		"http://radar.example.org",
		"https://radar.example.org",
		"192.168.1.1",
		"bad label!.example.org",
		"-starts-with-hyphen.example.org",
		"ends-with-hyphen-.example.org",
		"has space.example.org",
		"has/slash.example.org",
		"double..dot.example.org",
	}
	for _, s := range ok {
		if err := validateHostname(s); err != nil {
			t.Errorf("validateHostname(%q) unexpected error: %v", s, err)
		}
	}
	for _, s := range bad {
		if err := validateHostname(s); err == nil {
			t.Errorf("validateHostname(%q) expected error, got nil", s)
		}
	}
}

func TestValidateEmail(t *testing.T) {
	ok := []string{
		"ops@example.com",
		"user.name+tag@institution.ac.uk",
		"a@b.io",
	}
	bad := []string{
		"",
		"notanemail",
		"@nodomain",
		"noatsign",
		"missing@tld",
		"has space@example.com",
		"user@",
		"user@double..dot.com",
	}
	for _, s := range ok {
		if err := validateEmail(s); err != nil {
			t.Errorf("validateEmail(%q) unexpected error: %v", s, err)
		}
	}
	for _, s := range bad {
		if err := validateEmail(s); err == nil {
			t.Errorf("validateEmail(%q) expected error, got nil", s)
		}
	}
}

func TestValidateConfluentURL(t *testing.T) {
	ok := []string{
		"pkc-xxxxx.us-east-1.aws.confluent.cloud:9092",
		"localhost:9092",
		"broker.example.com:29092",
	}
	bad := []string{
		"",
		"noport",
		"https://broker.example.com:9092",
		"host:notaport",
		"host:0",
		"host:99999",
		":9092",
	}
	for _, s := range ok {
		if err := validateConfluentURL(s); err != nil {
			t.Errorf("validateConfluentURL(%q) unexpected error: %v", s, err)
		}
	}
	for _, s := range bad {
		if err := validateConfluentURL(s); err == nil {
			t.Errorf("validateConfluentURL(%q) expected error, got nil", s)
		}
	}
}

func TestValidateS3Bucket(t *testing.T) {
	ok := []string{
		"my-bucket",
		"radar-data-2024",
		"abc",
		"bucket.name.with.dots",
	}
	bad := []string{
		"",
		"ab",
		"UPPERCASE",
		"-starts-with-hyphen",
		"ends-with-hyphen-",
		".starts-with-dot",
		"has..double.dots",
		"has space",
		"192.168.1.1",
		"this-bucket-name-is-way-too-long-it-exceeds-the-sixty-three-char-limit",
	}
	for _, s := range ok {
		if err := validateS3Bucket(s); err != nil {
			t.Errorf("validateS3Bucket(%q) unexpected error: %v", s, err)
		}
	}
	for _, s := range bad {
		if err := validateS3Bucket(s); err == nil {
			t.Errorf("validateS3Bucket(%q) expected error, got nil", s)
		}
	}
}

func TestValidateS3Region(t *testing.T) {
	ok := []string{
		"us-east-1",
		"eu-west-2",
		"ap-southeast-1",
		"ca-central-1",
	}
	bad := []string{
		"",
		"us",
		"useast1",
		"has space",
	}
	for _, s := range ok {
		if err := validateS3Region(s); err != nil {
			t.Errorf("validateS3Region(%q) unexpected error: %v", s, err)
		}
	}
	for _, s := range bad {
		if err := validateS3Region(s); err == nil {
			t.Errorf("validateS3Region(%q) expected error, got nil", s)
		}
	}
}
