package main

import (
	"fmt"
	"github.com/MichalBures-OG/bp-bures-SfPDfSD-MQTT-preprocessor/src/mqtt"
	"github.com/MichalBures-OG/bp-bures-SfPDfSD-MQTT-preprocessor/src/types"
	"github.com/MichalBures-OG/bp-bures-SfPDfSD-commons/src/constants"
	"github.com/MichalBures-OG/bp-bures-SfPDfSD-commons/src/rabbitmq"
	cTypes "github.com/MichalBures-OG/bp-bures-SfPDfSD-commons/src/types"
	cUtil "github.com/MichalBures-OG/bp-bures-SfPDfSD-commons/src/util"
	"log"
	"sync"
	"time"
)

const (
	mqttBrokerURI      = "mqtt://localhost:1883"
	mqttClientID       = "bp-bures-SfPDfSD-MQTT-preprocessor"
	mqttTopic          = "topic"
	mqttBrokerUsername = "admin"
	mqttBrokerPassword = "password"
)

var (
	rabbitMQClient       rabbitmq.Client
	sdTypes              = cUtil.NewSet[string]()
	sdTypesMutex         sync.Mutex
	sdInstances          = cUtil.NewSet[cTypes.SDInstanceInfo]()
	sdInstancesMutex     sync.Mutex
	mqttMessageFIFO      = make([][]byte, 0)
	mqttMessageFIFOMutex sync.Mutex
)

func checkForSetOfSDTypesUpdates() {
	err := rabbitmq.ConsumeJSONMessages[[]string](rabbitMQClient, constants.SetOfSDTypesUpdatesQueueName, func(messagePayload []string) error {
		updatedSDTypes := cUtil.SliceToSet(messagePayload)
		sdTypesMutex.Lock()
		sdTypes = updatedSDTypes
		sdTypesMutex.Unlock()
		return nil
	})
	if err != nil {
		log.Printf("Consumption of messages from the '%s' queue has failed", constants.SetOfSDTypesUpdatesQueueName)
	}
}

func checkForSetOfSDInstancesUpdates() {
	err := rabbitmq.ConsumeJSONMessages[[]cTypes.SDInstanceInfo](rabbitMQClient, constants.SetOfSDInstancesUpdatesQueueName, func(messagePayload []cTypes.SDInstanceInfo) error {
		updatedSDInstances := cUtil.SliceToSet(messagePayload)
		sdInstancesMutex.Lock()
		sdInstances = updatedSDInstances
		sdInstancesMutex.Unlock()
		return nil
	})
	if err != nil {
		log.Printf("Consumption of messages from the '%s' queue has failed", constants.SetOfSDInstancesUpdatesQueueName)
	}
}

func mqttMessageSDTypeCorrespondsToSDTypeDefinitions(mqttMessageSDType string) bool {
	sdTypesMutex.Lock()
	ok := sdTypes.Contains(mqttMessageSDType)
	sdTypesMutex.Unlock()
	return ok
}

type sdInstanceScenario string

const (
	unknownSDInstance               = "unknown"
	sdInstanceNotYetConfirmedByUser = "notYetConfirmed"
	confirmedSDInstance             = "confirmed"
)

func determineSDInstanceScenario(uid string) sdInstanceScenario {
	var scenario sdInstanceScenario
	sdInstancesMutex.Lock()
	if sdInstances.Contains(cTypes.SDInstanceInfo{UID: uid, ConfirmedByUser: true}) {
		scenario = confirmedSDInstance
	} else if sdInstances.Contains(cTypes.SDInstanceInfo{UID: uid, ConfirmedByUser: false}) {
		scenario = sdInstanceNotYetConfirmedByUser
	} else {
		scenario = unknownSDInstance
	}
	sdInstancesMutex.Unlock()
	return scenario
}

func generateKPIFulfillmentCheckRequest(uid string, sdType string, parameters any, timestamp float32) {
	requestForKPIFulfillmentCheck := cTypes.RequestForKPIFulfillmentCheck{
		Timestamp: timestamp,
		SD: cTypes.SDInfo{
			UID:  uid,
			Type: sdType,
		},
		Parameters: parameters,
	}
	jsonSerializationResult := cUtil.SerializeToJSON(requestForKPIFulfillmentCheck)
	if jsonSerializationResult.IsFailure() {
		log.Println("Failed to serialize the object representing a KPI fulfillment check request into JSON")
	}
	err := rabbitMQClient.EnqueueJSONMessage(constants.KPIFulfillmentCheckRequestsQueueName, jsonSerializationResult.GetPayload())
	if err != nil {
		log.Println("Failed to publish a KPI fulfillment check request message") // TODO: This is here for debug purposes. Get rid of this line once it becomes unnecessary.
		return
	}
	log.Println("Successfully published a KPI fulfillment check request message") // TODO: This is here for debug purposes. Get rid of this line once it becomes unnecessary.
}

func generateSDInstanceRegistrationRequest(uid string, sdType string, timestamp float32) {
	requestForSDInstanceRegistration := cTypes.RequestForSDInstanceRegistration{
		Timestamp: timestamp,
		SD: cTypes.SDInfo{
			UID:  uid,
			Type: sdType,
		},
	}
	jsonSerializationResult := cUtil.SerializeToJSON(requestForSDInstanceRegistration)
	if jsonSerializationResult.IsFailure() {
		log.Println("Failed to serialize the object representing a SD instance registration request into JSON")
	}
	err := rabbitMQClient.EnqueueJSONMessage(constants.SDInstanceRegistrationRequestsQueueName, jsonSerializationResult.GetPayload())
	if err != nil {
		log.Println("Failed to publish a SD instance registration request message") // TODO: This is here for debug purposes. Get rid of this line once it becomes unnecessary.
		return
	}
	log.Println("Successfully published a SD instance registration request message") // TODO: This is here for debug purposes. Get rid of this line once it becomes unnecessary.
	sdInstancesMutex.Lock()
	sdInstances.Add(cTypes.SDInstanceInfo{
		UID:             uid,
		ConfirmedByUser: false,
	})
	sdInstancesMutex.Unlock()
}

func processMQTTMessagePayload(mqttMessagePayload []byte) {
	jsonDeserializationResult := cUtil.DeserializeFromJSON[types.UpstreamMQTTMessageInJSONBasedProprietaryFormatOfLogimic](mqttMessagePayload)
	if jsonDeserializationResult.IsFailure() {
		log.Println("Failed to deserialize the JSON payload of MQTT message")
		return
	}
	messagePayloadObject := jsonDeserializationResult.GetPayload()
	sd := messagePayloadObject.Data.SDArray[0]
	if !mqttMessageSDTypeCorrespondsToSDTypeDefinitions(sd.Type) {
		log.Println("Discarding the MQTT message: the SD type does not correspond to the current SD type definitions") // TODO: This is here for debug purposes. Get rid of this line once it becomes unnecessary.
		return
	}
	switch determineSDInstanceScenario(sd.UID) {
	case unknownSDInstance:
		generateSDInstanceRegistrationRequest(sd.UID, sd.Type, messagePayloadObject.Notification.Timestamp)
	case confirmedSDInstance:
		generateKPIFulfillmentCheckRequest(sd.UID, sd.Type, sd.Parameters, messagePayloadObject.Notification.Timestamp)
	case sdInstanceNotYetConfirmedByUser: // TODO: This is here for debug purposes. Get rid of this line once it becomes unnecessary.
		log.Println("Discarding the MQTT message: the SD instance is in the system, but has not yet been confirmed by the user") // TODO: This is here for debug purposes. Get rid of this line once it becomes unnecessary.
	}
}

func checkMQTTMessageFIFO() {
	for {
		mqttMessageFIFOMutex.Lock()
		if len(mqttMessageFIFO) > 0 {
			mqttMessagePayload := mqttMessageFIFO[0]
			processMQTTMessagePayload(mqttMessagePayload)
			mqttMessageFIFO = mqttMessageFIFO[1:]
			mqttMessageFIFOMutex.Unlock()
		} else {
			mqttMessageFIFOMutex.Unlock()
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func addIncomingMQTTMessageToFIFO(incomingMQTTMessagePayload []byte) {
	mqttMessageFIFOMutex.Lock()
	mqttMessageFIFO = append(mqttMessageFIFO, incomingMQTTMessagePayload)
	mqttMessageFIFOMutex.Unlock()
}

func main() {
	rabbitMQClient = rabbitmq.NewClient()
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		checkForSetOfSDTypesUpdates()
	}()
	go func() {
		defer wg.Done()
		checkForSetOfSDInstancesUpdates()
	}()
	mqttClient := mqtt.NewEclipsePahoBasedMqttClient(mqttBrokerURI, mqttClientID, mqttBrokerUsername, mqttBrokerPassword)
	cUtil.TerminateOnError(mqttClient.Connect(), fmt.Sprintf("Failed to connect to the MQTT broker [%s]", mqttBrokerURI))
	cUtil.TerminateOnError(mqttClient.Subscribe(mqttTopic, mqtt.QosAtLeastOnce, addIncomingMQTTMessageToFIFO), fmt.Sprintf("Failed to subscribe to the MQTT topic [%s]", mqttTopic))
	go func() {
		defer wg.Done()
		checkMQTTMessageFIFO()
	}()
	wg.Wait()
	rabbitMQClient.Dispose()
}