package auth

import (
	"testing"
	"time"
)

func TestSignVerifyRoundTrip(t *testing.T) {
	t.Parallel()
	key := []byte("test-signing-key")
	want := Identity{UserID: 522224, Login: "theding0x", IsOwner: true, Exp: time.Now().Add(time.Hour).Unix()}

	token, err := SignIdentity(want, key)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	got, err := VerifyIdentity(token, key, time.Now())
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if got != want {
		t.Fatalf("round-trip mismatch: got %+v want %+v", got, want)
	}
}

func TestVerifyRejectsTamperedSignature(t *testing.T) {
	t.Parallel()
	key := []byte("test-signing-key")
	token, _ := SignIdentity(Identity{UserID: 1, Exp: time.Now().Add(time.Hour).Unix()}, key)
	if _, err := VerifyIdentity(token+"x", key, time.Now()); err == nil {
		t.Fatal("expected error for tampered token")
	}
}

func TestVerifyRejectsWrongKey(t *testing.T) {
	t.Parallel()
	token, _ := SignIdentity(Identity{UserID: 1, Exp: time.Now().Add(time.Hour).Unix()}, []byte("key-a"))
	if _, err := VerifyIdentity(token, []byte("key-b"), time.Now()); err != ErrBadSignature {
		t.Fatalf("want ErrBadSignature, got %v", err)
	}
}

func TestVerifyRejectsExpired(t *testing.T) {
	t.Parallel()
	key := []byte("k")
	token, _ := SignIdentity(Identity{UserID: 1, Exp: time.Now().Add(-time.Minute).Unix()}, key)
	if _, err := VerifyIdentity(token, key, time.Now()); err != ErrExpired {
		t.Fatalf("want ErrExpired, got %v", err)
	}
}
