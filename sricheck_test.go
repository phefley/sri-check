package sricheck

import (
	"testing"
)

func TestKnownGood(t *testing.T) {
	url := "https://code.jquery.com/jquery-3.3.1.slim.min.js"
	integrity := "sha384-q8i/X+965DzO0rT7abK41JStQIAqVgRVzpbzo5smXKp4YfRvH+8abtTE1Pi6jizo"
	check, err := CheckIntegrity(url, integrity)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if !check {
		t.Fatalf("A known, valid integrity check failed.")
	}
}

func TestKnownBad(t *testing.T) {
	url := "https://code.jquery.com/jquery-3.3.1.slim.min.js"
	integrity := "sha384-q8i/X+965DzO0rT7abK41JStQIAqVgRVzpbzo5smXKp4YfRvH+8abtTE1Pi6jUWU"
	check, err := CheckIntegrity(url, integrity)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if check {
		t.Fatalf("A known, invalid integrity check passed.")
	}
}
