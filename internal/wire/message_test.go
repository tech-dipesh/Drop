package wire

import (
	"bytes"
	"net"
	"reflect"
	"testing"
)

func TestSendReceiveOverPipe(t *testing.T) {
	a, b := net.Pipe()
	defer a.Close()
	defer b.Close()
	offer := Offer{Version: Version, From: "Ana-PC", Files: []File{{"photo.jpg", 2048}, {"notes.pdf", 512}}}
	errs := make(chan error, 1)
	go func() { errs <- Send(a, offer) }()
	var got Offer
	if err := Receive(b, &got); err != nil {
		t.Fatal(err)
	}
	if err := <-errs; err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, offer) {
		t.Fatalf("got %+v, want %+v", got, offer)
	}
}

func TestOfferTotal(t *testing.T) {
	offer := Offer{Files: []File{{"a", 100}, {"b", 50}, {"empty", 0}}}
	if total := offer.Total(); total != 150 {
		t.Fatalf("got %d, want 150", total)
	}
}

func TestReceiveBadJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := writeFrame(&buf, []byte("{")); err != nil {
		t.Fatal(err)
	}
	var answer Answer
	if err := Receive(&buf, &answer); err == nil {
		t.Fatal("expected a decode error")
	}
}
