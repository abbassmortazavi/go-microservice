package message

import (
	notification "abbassmortazavi/go-microservice/services/notification-service/messaging"
)

func GetPublisher() *notification.Publisher {
	return Publisher
}
