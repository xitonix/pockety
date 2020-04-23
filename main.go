package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"runtime"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/xitonix/pockety/auth"
)

func main() {
	consumerKey := kingpin.Arg("consumer-key", "Pocket consumer key").Required().String()
	kingpin.Parse()

	ch := make(chan struct{})
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == "/favicon.ico" {
				http.Error(w, "Not Found", 404)
				return
			}

			w.Header().Set("Content-Type", "text/html")
			_, _ = fmt.Fprintln(w, "<h3>Authorized. You can close the browser now</h3>")
			ch <- struct{}{}
		}))
	defer ts.Close()

	redirectURL := ts.URL

	requestToken, err := auth.ObtainRequestToken(*consumerKey, redirectURL)
	if err != nil {
		log.Fatalf("Failed to request a new authentication token: %s", err)
	}

	url := auth.GenerateAuthorizationURL(requestToken, redirectURL)
	if err := openBrowser(url); err != nil {
		fmt.Printf("Failed to open the default browser. Please open a browser and go to: %s\n", url)
	}

	<-ch
	auth, err := auth.ObtainAccessToken(*consumerKey, requestToken)
	if err != nil {
		log.Fatalf("Failed to obtain a new authentication token: %s", err)
	}

	fmt.Printf("Consumer Key: %s\n", *consumerKey)
	fmt.Printf("  Auth Token: %s\n", auth.AccessToken)
	fmt.Printf("    Username: %s\n", auth.Username)
}

func openBrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		return err
	}
	return nil
}
