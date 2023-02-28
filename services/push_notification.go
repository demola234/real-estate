package services

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/demola234/real-estate/interfaces"
	"google.golang.org/api/option"
)

func SendPushNotificationToAll(p *interfaces.PushNotificationToAll) error {

	app, _, _ := SetupFirebase()

	ctx := context.Background()
	client, error := app.Messaging(ctx)
	if error != nil {
		log.Fatalln(error)
	}

	response, error := client.SendMulticast(ctx, &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: p.Title,
			Body:  p.Body,
		},
		Tokens: p.Tokens,
	})

	if error != nil {
		log.Fatalln(error)
	}

	fmt.Println("Successfully sent message:", response)
	return error
}

func SendPushNotificationToUser(p *interfaces.PushNotificationToUser) error {

	app, _, _ := SetupFirebase()

	ctx := context.Background()
	client, error := app.Messaging(ctx)
	if error != nil {
		log.Fatalln(error)
	}

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: p.Title,
			Body:  p.Body,
		},
		Token: p.To,
	}

	response, error := client.Send(ctx, message)

	fmt.Println("Successfully sent message:", response)
	return error
}

func SetupFirebase() (*firebase.App, context.Context, *messaging.Client) {

	ctx := context.Background()

	//Firebase service account key file path
	serviceAccountKeyFilePath, err := filepath.Abs("homely-push-notification-firebase.json")
	if err != nil {
		panic("Unable to load homely-push-notification-firebase.json file")
	}

	opt := option.WithCredentialsFile(serviceAccountKeyFilePath)

	//Firebase admin SDK initialization
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic("Firebase load error")
	}

	//Messaging client
	client, _ := app.Messaging(ctx)

	return app, ctx, client
}
