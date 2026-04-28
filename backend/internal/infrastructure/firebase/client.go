package firebase

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

type Client struct {
	App  *firebase.App
	Auth *auth.Client
}

func NewClient(ctx context.Context) *Client {
	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	saJSON := os.Getenv("FIREBASE_SERVICE_ACCOUNT_JSON")
	saPath := os.Getenv("FIREBASE_SERVICE_ACCOUNT_KEY")

	var app *firebase.App
	var err error

	if saJSON != "" {
		// JSON文字列から初期化
		opt := option.WithCredentialsJSON([]byte(saJSON))
		app, err = firebase.NewApp(ctx, &firebase.Config{ProjectID: projectID}, opt)
	} else if saPath != "" {
		// ファイルパスから初期化
		if _, err := os.Stat(saPath); os.IsNotExist(err) {
			log.Printf("WARNING: Firebase service account file not found at %s. Falling back to default credentials.", saPath)
			app, err = firebase.NewApp(ctx, &firebase.Config{ProjectID: projectID})
		} else {
			opt := option.WithServiceAccountFile(saPath)
			app, err = firebase.NewApp(ctx, &firebase.Config{ProjectID: projectID}, opt)
		}
	} else {
		// デフォルト（環境変数など）
		app, err = firebase.NewApp(ctx, &firebase.Config{ProjectID: projectID})
	}

	if err != nil {
		log.Fatalf("error initializing firebase app: %v", err)
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting firebase auth client: %v", err)
	}

	return &Client{
		App:  app,
		Auth: authClient,
	}
}
