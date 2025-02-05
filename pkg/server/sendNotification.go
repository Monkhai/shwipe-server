package server

import (
	"log"

	"github.com/Monkhai/shwipe-server.git/pkg/notifications"
	"github.com/wagon-official/expo-notifications-sdk-golang"
)

func (s *Server) sendNotification(userPushToken, sessionId string, notificationType notifications.NotificationType) {
	pushToken, err := expo.NewExpoPushToken(userPushToken)
	if err != nil {
		log.Printf("Error creating push token: %v", err)
		return
	}

	client := expo.NewPushClient(nil)

	switch notificationType {
	case notifications.NotificationTypeSessionInvitation:
		{
			msg := &expo.PushMessage{
				To:    []expo.ExpoPushToken{pushToken},
				Title: "Session",
				Body:  "You have been invited to a session",
				Data: map[string]interface{}{
					"type":      notificationType,
					"sessionId": sessionId,
				},
			}
			_, err := client.Publish(msg)
			if err != nil {
				log.Printf("Error sending notification: %v", err)
			}
		}

	}
}
