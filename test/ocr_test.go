package tests

import (
	"testing"

	"github.com/otiai10/gosseract/v2"
	"github.com/stretchr/testify/assert"
)

func TestOCR(t *testing.T) {
	client := gosseract.NewClient()
	defer client.Close()

	// Set the image to be processed
	client.SetImage("test-images/1.png")

	// Perform OCR on the image
	text, err := client.Text()
	if err != nil {
		t.Fatalf("Failed to perform OCR: %v", err)
	}

	// Assert that the recognized text is as expected
	expectedText := "Hello, World!"
	assert.Equal(t, expectedText, text, "The recognized text does not match the expected text")
}
