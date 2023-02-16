package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/logging"
	"github.com/erandelax/go-dossh/internal/commands"
	"github.com/erandelax/go-dossh/internal/configuration"
	gossh "golang.org/x/crypto/ssh"
)

//
// SSH server
//

func main() {
	cfg := configuration.Get()

	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)),
		wish.WithHostKeyPath("./.ssh/host.key"),
		wish.WithMiddleware(
			func(h ssh.Handler) ssh.Handler {
				return func(s ssh.Session) {
					userConfig, ok := cfg.Users[s.User()]
					if !ok {
						return
					}
					authorizedKey := strings.Trim(string(gossh.MarshalAuthorizedKey(s.PublicKey())), " \n\r\t")
					userKey := strings.Trim(userConfig.PublicKey[:len(authorizedKey)], " \n\r\t")
					if userKey != authorizedKey {
						s.Write([]byte("Authorization failed: SSH key mismatch. Please send your new public key to server administrator to register it."))
						return
					}

					rootCmd := commands.NewRoot(userConfig, s.User(), s)
					rootCmd.SetArgs(s.Command())
					rootCmd.SetIn(s)
					rootCmd.SetOut(s)
					rootCmd.SetErr(s.Stderr())
					rootCmd.CompletionOptions.DisableDefaultCmd = true
					if err := rootCmd.Execute(); err != nil {
						_ = s.Exit(1)
						return
					}
					s.Write([]byte("\n"))

					h(s)
				}
			},
			logging.Middleware(),
		),
		ssh.PublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
			user := ctx.User()
			_, ok := cfg.Users[user]
			return ok
		}),
	)
	if err != nil {
		log.Fatalln(err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("Starting SSH server on %s:%d", cfg.Host, cfg.Port)
	go func() {
		if err = s.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()

	<-done
	log.Println("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}
}
