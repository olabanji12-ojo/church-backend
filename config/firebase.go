package config

import (
	"context"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

var FirebaseApp *firebase.App

func InitFirebase() {
	// The path to our service account key we moved earlier
	// We'll use os.Getwd() to ensure the path is absolute from the backend root
	cwd, _ := os.Getwd()
	opt := option.WithCredentialsFile(cwd + "/firebase-admin.json")

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		fmt.Printf("❌ error initializing Firebase app: %v\n", err)
		return
	}

	FirebaseApp = app
	fmt.Println("✅ Successfully initialized Firebase Admin SDK")
}
