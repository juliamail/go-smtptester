package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/smtp"
	"os"
	"os/signal"
	"time"

	"github.com/uponusolutions/go-smtptester"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	s := smtptester.StandardWithAddress(":2525")

	go func() {
		if err := s.ListenAndServe(ctx); err != nil {
			slog.Error("smtp server response %s", slog.Any("error", err))
		}
	}()

	defer func() {
		if err := s.Close(); err != nil {
			slog.Error("error closing server", slog.Any("error", err))
		}
	}()

	// Wait a second to let the server come up.
	time.Sleep(time.Second)

	// Send email.
	from := "alice@i.com"
	to := []string{"bob@e.com", "mal@b.com"}
	msg := []byte("Test\r\n")
	if err := smtp.SendMail("[::1]:2525", nil, from, to, msg); err != nil {
		panic(err)
	}

	// Lookup email.
	m, found := smtptester.GetBackend(s).Load(from, to)
	fmt.Printf("Found %t, mail %+v\n", found, m)
}
