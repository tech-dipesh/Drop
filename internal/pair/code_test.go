package pair

import (
	"regexp"
	"testing"
)

func TestBothSidesGetTheSameCode(t *testing.T) {
	client, server := pipe(t)
	codes := make(chan string, 1)
	go func() {
		code, err := AsReceiver(server)
		if err != nil {
			t.Error(err)
		}
		codes <- code
	}()
	mine, err := AsSender(client)
	if err != nil {
		t.Fatal(err)
	}
	if theirs := <-codes; theirs != mine {
		t.Fatalf("sender sees %s, receiver sees %s", mine, theirs)
	}
	if !regexp.MustCompile(`^\d{3} \d{3}$`).MatchString(mine) {
		t.Fatalf("code %q isn't six digits", mine)
	}
}

func TestSeparateSessionsGetDifferentCodes(t *testing.T) {
	first, _ := pipe(t)
	second, _ := pipe(t)
	sender, receiver := []byte("same sender nonce"), []byte("same receiver nonce")
	a, err := code(first, sender, receiver)
	if err != nil {
		t.Fatal(err)
	}
	b, err := code(second, sender, receiver)
	if err != nil {
		t.Fatal(err)
	}
	if a == b {
		t.Fatalf("two tls sessions produced the same code %s", a)
	}
}
