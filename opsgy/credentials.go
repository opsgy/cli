package opsgy

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/skratchdot/open-golang/open"
	"golang.org/x/oauth2"
)

var (
	conf        *oauth2.Config
	ctx         context.Context
	state       *string
	host        = "localhost"
	credentials *Credentials
	stop        chan bool
)

type Credentials struct {
	AccessToken  string    `yaml:"accessToken,omitempty"`
	RefreshToken string    `yaml:"refreshToken,omitempty"`
	TokenExpiry  time.Time `yaml:"tokenExpiry,omitempty"`
	TokenType    string    `yaml:"tokenType,omitempty"`
}

func generateRandomString(n int) (*string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	var hexString = hex.EncodeToString(b)
	return &hexString, nil
}

func handleCallback(w http.ResponseWriter, req *http.Request) {
	queryParts, _ := url.ParseQuery(req.URL.RawQuery)

	responseState := queryParts["state"][0]

	if responseState != *state {
		log.Println(responseState)
		log.Println(*state)
		log.Fatal("Incorrect state")
	}

	code := queryParts["code"][0]

	// Exchange will do the handshake to retrieve the initial access token.
	token, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}
	// The HTTP Client returned by conf.Client will refresh the token as necessary.
	// client := conf.Client(ctx, tok)
	fmt.Println()
	fmt.Println("You are authenticated!")

	credentials = &Credentials{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenExpiry:  token.Expiry,
		TokenType:    token.TokenType,
	}

	// show succes page
	msg := "<script>window.close();</script>"
	msg = msg + "<p>You are authenticated and can now return to the CLI.</p>"
	fmt.Fprintf(w, msg)

	// Stop process
	go func() {
		stop <- true
	}()
}

func Login(port int, scopes []string) (*Credentials, error) {
	ctx = context.Background()
	conf = &oauth2.Config{
		ClientID:     ClientID,
		ClientSecret: ClientSecret,
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  OpsgyAccountUrl + "/oauth/authorize",
			TokenURL: OpsgyAccountUrl + "/oauth/token",
		},
		// Own callback URL
		RedirectURL: fmt.Sprintf("http://%s:%d/oauth/callback", host, port),
	}

	var err error
	state, err = generateRandomString(16)
	if err != nil {
		return nil, err
	}

	url := conf.AuthCodeURL(*state, oauth2.AccessTypeOffline)

	fmt.Println("You will now be taken to your browser for authentication")
	open.Run(url)
	fmt.Println("Authentication URL:" + url)

	mux := http.NewServeMux()
	mux.HandleFunc("/oauth/callback", handleCallback)

	server := &http.Server{Addr: fmt.Sprintf("%s:%d", host, port), Handler: mux}

	go func() {
		err = server.ListenAndServe()
		if err != nil && fmt.Sprint(err) != "http: Server closed" {
			log.Fatal(err)
		}
	}()

	// Setting up signal capturing
	stop = make(chan bool, 1)

	// Waiting for SIGINT (pkill -2)
	<-stop

	if err := server.Shutdown(ctx); err != nil {
		return nil, err
	}

	return credentials, nil
}

func Logout() error {
	config := LoadConfig()
	config.AccessToken = ""
	config.RefreshToken = ""
	return SaveConfig(config)
}

func GetClient(config *Config) (*http.Client, error) {
	ctx = context.Background()
	conf = &oauth2.Config{
		ClientID:     ClientID,
		ClientSecret: ClientSecret,
		Endpoint: oauth2.Endpoint{
			TokenURL: OpsgyAccountUrl + "/oauth/token",
		},
	}

	token := &oauth2.Token{
		AccessToken:  config.AccessToken,
		RefreshToken: config.RefreshToken,
		Expiry:       config.TokenExpiry,
		TokenType:    config.TokenType,
	}

	if token.Expiry.Before(time.Now()) {
		src := conf.TokenSource(ctx, token)
		newToken, err := src.Token() // this actually goes and renews the tokens
		if err != nil {
			return nil, err
		}
		if newToken.AccessToken != token.AccessToken {
			config.AccessToken = newToken.AccessToken
			config.RefreshToken = newToken.RefreshToken
			config.TokenExpiry = newToken.Expiry
			config.TokenType = newToken.TokenType
			err = SaveConfig(config)
			if err != nil {
				return nil, err
			}

			token = newToken
		}
	}

	client := conf.Client(ctx, token)

	return client, nil
}
