package message

import (
	messaging2 "abbassmortazavi/go-microservice/services/auth-service/messaging"
)

func GetPublisher() *messaging2.Publisher {
	return Publisher
}
