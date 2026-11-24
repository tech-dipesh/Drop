package ui

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/tech-dipesh/drop/internal/peer"
)

func Pick(ctx context.Context, found <-chan peer.Peer) (peer.Peer, error) {
	fmt.Println("searching nearby devices, type a number and press enter")
	var peers []peer.Peer
	answer := readLine()
	hint := time.After(3 * time.Second)
	for {
		select {
		case p := <-found:
			peers = append(peers, p)
			fmt.Printf("  %d  %s\n", len(peers), clean(p.Name))
		case text := <-answer:
			n, err := strconv.Atoi(text)
			if err == nil && n >= 1 && n <= len(peers) {
				return peers[n-1], nil
			}
			fmt.Println("type a number from the list")
			answer = readLine()
		case <-hint:
			if len(peers) == 0 {
				fmt.Println("nothing yet, run drop on the other device and keep both on the same network")
			}
		case <-ctx.Done():
			return peer.Peer{}, ctx.Err()
		}
	}
}

func readLine() <-chan string {
	line := make(chan string, 1)
	go func() {
		text, err := stdin.ReadString('\n')
		if err != nil && text == "" {
			return
		}
		line <- strings.ToLower(strings.TrimSpace(text))
	}()
	return line
}

func clean(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, s)
}
