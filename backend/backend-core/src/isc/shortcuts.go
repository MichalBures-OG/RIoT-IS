package isc

import (
	"fmt"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/rabbitmq"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedConstants"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

func consumeMessageProcessingUnitConnectionNotificationJSONMessages(messageProcessingUnitConnectionNotificationConsumerFunction func(messageProcessingUnitConnectionNotification sharedModel.MessageProcessingUnitConnectionNotification) error) {
	err := rabbitmq.ConsumeJSONMessages[sharedModel.MessageProcessingUnitConnectionNotification](getRabbitMQClient(), sharedConstants.MessageProcessingUnitConnectionNotificationsQueue, messageProcessingUnitConnectionNotificationConsumerFunction)
	errorMessage := fmt.Sprintf("[ISC] Consumption of messages from the '%s' queue has failed", sharedConstants.MessageProcessingUnitConnectionNotificationsQueue)
	sharedUtils.TerminateOnError(err, errorMessage)
}

func consumeSDInstanceRegistrationRequestJSONMessages(sdInstanceRegistrationRequestConsumerFunction func(sdInstanceRegistrationRequest sharedModel.SDInstanceRegistrationRequestISCMessage) error) {
	err := rabbitmq.ConsumeJSONMessages[sharedModel.SDInstanceRegistrationRequestISCMessage](getRabbitMQClient(), sharedConstants.SDInstanceRegistrationRequestsQueueName, sdInstanceRegistrationRequestConsumerFunction)
	errorMessage := fmt.Sprintf("[ISC] Consumption of messages from the '%s' queue has failed", sharedConstants.SDInstanceRegistrationRequestsQueueName)
	sharedUtils.TerminateOnError(err, errorMessage)
}

func consumeKPIFulfillmentCheckResultJSONMessages(kpiFulfillmentCheckResultConsumerFunction func(kpiFulfillmentCheckResult sharedModel.KPIFulfillmentCheckResultISCMessage) error) {
	err := rabbitmq.ConsumeJSONMessages[sharedModel.KPIFulfillmentCheckResultISCMessage](getRabbitMQClient(), sharedConstants.KPIFulfillmentCheckResultsQueueName, kpiFulfillmentCheckResultConsumerFunction)
	errorMessage := fmt.Sprintf("[ISC] Consumption of messages from the '%s' queue has failed", sharedConstants.KPIFulfillmentCheckResultsQueueName)
	sharedUtils.TerminateOnError(err, errorMessage)
}
