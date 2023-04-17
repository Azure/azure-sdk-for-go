//go:build go1.20
// +build go1.20

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aznotificationhubs

import (
	"encoding/json"
	"testing"
)

func TestInstallationSimple(t *testing.T) {
	jsonStr := `
        {
            "installationId": "123",
            "userId": "456",
            "lastActiveOn": "2022-01-01T00:00:00Z",
            "expirationTime": "2023-01-01T00:00:00Z",
            "lastUpdate": "2022-04-01T00:00:00Z",
            "platform": "apple",
            "expiredPushChannel": true,
            "tags": ["tag1", "tag2"],
            "templates": {
                "template1": {
                    "body": "Hello, {{name}}!"
                },
                "template2": {
                    "body": "Goodbye, {{name}}!"
                }
            },
            "pushChannel": "abc123"
        }
    `
	var installation Installation
	err := json.Unmarshal([]byte(jsonStr), &installation)
	if err != nil {
		t.Errorf("error unmarshalling JSON: %v", err)
	}

	if installation.InstallationID != "123" {
		t.Errorf("expected InstallationId to be 123, got %s", installation.InstallationID)
	}

	if installation.UserId != "456" {
		t.Errorf("expected UserId to be 456, got %s", installation.UserId)
	}

	if installation.Platform != "apple" {
		t.Errorf("expected Platform to be apple, got %s", installation.Platform)
	}

	if installation.Templates["template1"].Body != "Hello, {{name}}!" {
		t.Errorf("expected template1.Body to be 'Hello, {{name}}!', got '%s'", installation.Templates["template1"].Body)
	}

	if installation.Templates["template2"].Body != "Goodbye, {{name}}!" {
		t.Errorf("expected template2.Body to be 'Goodbye, {{name}}!', got '%s'", installation.Templates["template2"].Body)
	}

	pushChannelType := installation.PushChannel
	switch pushChannelType.(type) {
	case string:
		if installation.PushChannel.(string) != "abc123" {
			t.Errorf("expected PushChannel to be 'abc123', got %v", installation.PushChannel)
		}
	default:
		t.Errorf("expected PushChannel to be 'abc123', got %v", installation.PushChannel)
	}
}

func TestInstallationWebPush(t *testing.T) {
	jsonStr := `
        {
            "installationId": "123",
            "userId": "456",
            "lastActiveOn": "2022-01-01T00:00:00Z",
            "expirationTime": "2023-01-01T00:00:00Z",
            "lastUpdate": "2022-04-01T00:00:00Z",
            "platform": "apple",
            "expiredPushChannel": true,
            "tags": ["tag1", "tag2"],
            "templates": {
                "template1": {
                    "body": "Hello, {{name}}!"
                },
                "template2": {
                    "body": "Goodbye, {{name}}!"
                }
            },
            "pushChannel": {
				"endpoint": "https://example.com",
				"p256dh": "abc123",
				"auth": "def456"
			}
        }
    `
	var installation Installation
	err := json.Unmarshal([]byte(jsonStr), &installation)
	if err != nil {
		t.Errorf("error unmarshalling JSON: %v", err)
	}

	if installation.InstallationID != "123" {
		t.Errorf("expected InstallationId to be 123, got %s", installation.InstallationID)
	}

	if installation.UserId != "456" {
		t.Errorf("expected UserId to be 456, got %s", installation.UserId)
	}

	if installation.Platform != "apple" {
		t.Errorf("expected Platform to be apple, got %s", installation.Platform)
	}

	if installation.Templates["template1"].Body != "Hello, {{name}}!" {
		t.Errorf("expected template1.Body to be 'Hello, {{name}}!', got '%s'", installation.Templates["template1"].Body)
	}

	if installation.Templates["template2"].Body != "Goodbye, {{name}}!" {
		t.Errorf("expected template2.Body to be 'Goodbye, {{name}}!', got '%s'", installation.Templates["template2"].Body)
	}

	switch pushChannel := installation.PushChannel.(type) {
	case map[string]interface{}:
		if pushChannel["endpoint"].(string) != "https://example.com" {
			t.Errorf("expected PushChannel.Endpoint to be 'https://example.com', got %v", pushChannel["endpoint"])
		}
		if pushChannel["p256dh"].(string) != "abc123" {
			t.Errorf("expected PushChannel.P256DH to be 'abc123', got %v", pushChannel["p256dh"])
		}
		if pushChannel["auth"].(string) != "def456" {
			t.Errorf("expected PushChannel.Auth to be 'abc123', got %v", pushChannel["auth"])
		}
	default:
		t.Errorf("expected PushChannel to be map[string]interface{}, got %v", installation.PushChannel)
	}
}
