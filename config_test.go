package main

import "testing"

func TestConfig(t *testing.T) {
	configText := `
    access "allow" {       
        path = "/api"
    }

    access "allow" {
        user = "simas"
    }

    access "deny" {
        user = "serviceaccount:*"
    }

    access "allow" {
        path = "/apis"
    }
    `
	LoadConfigFromByteArray([]byte(configText))

	if len(config.Rules) != 4 {
		t.Errorf("Config: error on rules length - %d", len(config.Rules))
	}

	r := config.Rules[0]
	if r.Mode != "allow" || r.Path != "/api" {
		t.Errorf("Config: error on rules [0]")
	}
}
