package service

import (
	"testing"
)

func TestHashPasswordProducesDifferentValuesForSameInput(t *testing.T) {
	first, err := hashPassword("super-secret-1")
	if err != nil {
		t.Fatalf("hashPassword error = %v", err)
	}
	second, err := hashPassword("super-secret-1")
	if err != nil {
		t.Fatalf("hashPassword error = %v", err)
	}
	if first == second {
		t.Fatal("hashes should not be equal for same password due to salt")
	}
	if first == "super-secret-1" {
		t.Fatal("password should not be stored in plain text")
	}
	if len(first) < 40 {
		t.Fatalf("hash length = %d, want a bcrypt hash", len(first))
	}
}

func TestHashPasswordRejectsEmpty(t *testing.T) {
	_, err := hashPassword("")
	if err == nil {
		t.Fatal("expected error for empty password")
	}
}

func TestComparePasswordMatchesHashedPassword(t *testing.T) {
	hashed, err := hashPassword("correct-horse-battery")
	if err != nil {
		t.Fatalf("hashPassword error = %v", err)
	}

	if !comparePassword(hashed, "correct-horse-battery") {
		t.Fatal("comparePassword should return true for matching password")
	}
	if comparePassword(hashed, "wrong-password") {
		t.Fatal("comparePassword should return false for mismatched password")
	}
	if comparePassword(hashed, "") {
		t.Fatal("comparePassword should return false for empty password")
	}
	if comparePassword("", "correct-horse-battery") {
		t.Fatal("comparePassword should return false for empty hash")
	}
}
