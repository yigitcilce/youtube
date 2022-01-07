package main

import "testing"

func TestConvertVideoTitletoFileName(t *testing.T) {
	// Invalid characters
	fileName := "a<b>c:d\\e\"f/g\\h|i?j*k****"
	proper := ConvertVideoTitletoFileName(fileName)
	if proper != "abcdefghijk" {
		t.Error("Invalid characters must get stripped")
	}

	// Already proper
	fileName = "aB Cd"
	proper = ConvertVideoTitletoFileName(fileName)
	if proper != "aB Cd" {
		t.Error("Upper lower case and whitespaces must be preserved")
	}

	// Interestingly allowed characters
	fileName = "~!@#$%^&()[].,"
	proper = ConvertVideoTitletoFileName(fileName)
	if proper != "~!@#$%^&()[].," {
		t.Error("Allowed characters by OS must be preserved")
	}
}
