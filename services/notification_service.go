package services

import (
	"context"

	firebaseMessaging "firebase.google.com/go/v4/messaging"
	"github.com/olabanji12-ojo/church-backend/config"
	"github.com/sirupsen/logrus"
)

// NotificationService defines the contract for sending notifications.
// This interface allows us to easily swap out Firebase for another provider (like Kafka, AWS SNS, etc.)
// in the future without touching the core MessageService logic.
type NotificationService interface {
	SendPush(pushToken, title, body, url string) error
}

// FirebaseNotificationService is the concrete implementation of NotificationService using Firebase Cloud Messaging.
type FirebaseNotificationService struct{}

func NewFirebaseNotificationService() *FirebaseNotificationService {
	return &FirebaseNotificationService{}
}

// SendPush executes the actual push notification via the Firebase Admin SDK.
func (f *FirebaseNotificationService) SendPush(pushToken, title, body, url string) error {
	if config.FirebaseApp == nil {
		logrus.Warn("FirebaseApp is not initialized, skipping push notification.")
		return nil
	}

	client, err := config.FirebaseApp.Messaging(context.Background())
	if err != nil {
		logrus.Error("Failed getting Firebase Messaging client: ", err)
		return err
	}

	message := &firebaseMessaging.Message{
		Token: pushToken,
		Notification: &firebaseMessaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: map[string]string{
			"url": url,
		},
	}

	response, err := client.Send(context.Background(), message)
	if err != nil {
		logrus.Error("Failed to send Firebase Push Notification: ", err)
		return err
	}

	logrus.Info("Firebase Push sent successfully! ID: ", response)
	return nil
}
