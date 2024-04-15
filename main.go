package main

import (
	"assignment/pkg/reddit"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	clientID := "MnqvvrkYBC8TYjVMxbb37A"
	clientSecret := "ko05fJ_Mpz0zkqf2ZoFf3LBrP5r4Bw"
	redirectURI := "https://localhost/token"

	redditClient := reddit.NewClient("https://oauth.reddit.com", clientID, clientSecret, redirectURI)

	authURL := redditClient.AuthURL("code", "permanent", "read report history")
	fmt.Println("Authorize this app by visiting the URL:", authURL)

	var code string
	fmt.Print("Enter the authorization code: ")
	fmt.Scan(&code)

	// code := "er1wdtFqXfu_SHcEknqNOgB2xB11zA"
	code = code[:len(code)-2]

	ctx := context.Background()
	token, err := redditClient.ExchangeCode(ctx, code)
	if err != nil {
		log.Fatalf("Failed to exchange authorization code: %v", err)
	}

	redditClient.SetAccessToken(token)

	// res, _ := redditClient.GetPosts(ctx, "funny")
	// x, _ := json.Marshal(res)
	// fmt.Println(string(x))

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		redditClient.ProcessPosts(ctx, "funny")
	}()

	// Wait for termination signal (Ctrl+C)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	<-sigCh

}
