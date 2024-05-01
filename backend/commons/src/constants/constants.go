package constants

const (
	MIMETypeOfJSONData                         = "application/json"
	URLOfRabbitMQ                              = "amqp://guest:guest@rabbitmq:5672/"
	ContextTimeoutInSeconds                    = 5
	KPIFulfillmentCheckRequestsQueueName       = "kpi-fulfillment-check-requests"
	KPIFulfillmentCheckResultsQueueName        = "kpi-fulfillment-check-results"
	KPIDefinitionsBySDTypeDenotationMapUpdates = "kpi-definitions-by-sd-type-denotation-map-updates"
	SDInstanceRegistrationRequestsQueueName    = "sd-instance-registration-requests"
	SetOfSDInstancesUpdatesQueueName           = "set-of-sd-instances-updates"
	SetOfSDTypesUpdatesQueueName               = "set-of-sd-types-updates"
)
