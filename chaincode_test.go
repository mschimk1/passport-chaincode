package main

import "testing"

func TestGenerateHash(t *testing.T) {
	msgHash1 := generateHash("{\"id\":\"1\", \"from_account\":\"2\", \"to_account\":\"2\", \"amount\":1000, \"fee\":0, \"currency\":\"AUD\"\"}")
	msgHash2 := generateHash("{\"id\":\"2\", \"from_account\":\"2\", \"to_account\":\"1\", \"amount\":1000, \"fee\":0, \"currency\":\"AUD\"\"}")
	msgHash3 := generateHash("{\"id\":\"1\", \"from_account\":\"2\", \"to_account\":\"2\", \"amount\":1000, \"fee\":0, \"currency\":\"AUD\"\"}")
	if msgHash1 == msgHash2 {
		t.Error("Expect different message hashes, but got the same")
	}
	if msgHash1 != msgHash3 {
		t.Error("Expected same messages hashes, but got different ones")
	}
}

func TestGenerateID(t *testing.T) {
	id1 := generateID(8)
	id2 := generateID(8)
	if id1 == id2 {
		t.Error("Expect different message hashes, but got the same")
	}
}

func TestMin(t *testing.T) {
	if min(2, 3) != 2 {
		t.Error("Expected 2 to be the minimum")
	}
}
