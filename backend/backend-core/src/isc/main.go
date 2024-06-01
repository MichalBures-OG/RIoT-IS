package isc

import (
	"errors"
	"fmt"
	"github.com/MichalBures-OG/bp-bures-SfPDfSD-backend-core/src/db/dbClient"
	"github.com/MichalBures-OG/bp-bures-SfPDfSD-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-SfPDfSD-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-SfPDfSD-backend-core/src/modelMapping/dll2gql"
	"github.com/MichalBures-OG/bp-bures-SfPDfSD-commons/src/rabbitmq"
	"github.com/MichalBures-OG/bp-bures-SfPDfSD-commons/src/sharedConstants"
	"github.com/MichalBures-OG/bp-bures-SfPDfSD-commons/src/sharedModel"
	"github.com/MichalBures-OG/bp-bures-SfPDfSD-commons/src/sharedUtils"
	"sync"
)

var (
	rabbitMQClient rabbitmq.Client
	once           sync.Once
)

func getRabbitMQClient() rabbitmq.Client {
	once.Do(func() {
		rabbitMQClient = rabbitmq.NewClient()
	})
	return rabbitMQClient
}

func ProcessIncomingSDInstanceRegistrationRequests(sdInstanceGraphQLSubscriptionChannel *chan graphQLModel.SDInstance) {
	consumeSDInstanceRegistrationRequestJSONMessages(func(sdInstanceRegistrationRequestISCMessage sharedModel.SDInstanceRegistrationRequestISCMessage) error {
		newSDInstanceUID := sdInstanceRegistrationRequestISCMessage.SDInstanceUID
		newSDInstanceSDTypeSpecification := sdInstanceRegistrationRequestISCMessage.SDTypeSpecification
		newSDInstancePersistResult := dbClient.GetRelationalDatabaseClientInstance().PersistNewSDInstance(newSDInstanceUID, newSDInstanceSDTypeSpecification)
		if newSDInstancePersistResult.IsSuccess() {
			select {
			case *sdInstanceGraphQLSubscriptionChannel <- dll2gql.ToGraphQLModelSDInstance(newSDInstancePersistResult.GetPayload()):
			default:
			}
		} else if newSDInstancePersistError := newSDInstancePersistResult.GetError(); !errors.Is(newSDInstancePersistError, dbClient.ErrOperationWouldLeadToForeignKeyIntegrityBreach) {
			return fmt.Errorf("failed to persist a new SD instance with UID = %s and SD type specification = %s: %w", newSDInstanceUID, newSDInstanceSDTypeSpecification, newSDInstancePersistError)
		}
		return nil
	})
}

func ProcessIncomingKPIFulfillmentCheckResults(kpiFulfillmentCheckResultGraphQLSubscriptionChannel *chan graphQLModel.KPIFulfillmentCheckResult) {
	consumeKPIFulfillmentCheckResultJSONMessages(func(kpiFulfillmentCheckResultISCMessage sharedModel.KPIFulfillmentCheckResultISCMessage) error {
		kpiDefinitionID := kpiFulfillmentCheckResultISCMessage.KPIDefinitionID
		sdInstanceUID := kpiFulfillmentCheckResultISCMessage.SDInstanceUID
		fulfilled := kpiFulfillmentCheckResultISCMessage.Fulfilled
		kpiFulfillmentCheckResultPersistResult := dbClient.GetRelationalDatabaseClientInstance().PersistKPIFulFulfillmentCheckResult(kpiDefinitionID, sdInstanceUID, fulfilled)
		if kpiFulfillmentCheckResultPersistResult.IsSuccess() {
			select {
			case *kpiFulfillmentCheckResultGraphQLSubscriptionChannel <- dll2gql.ToGraphQLModelKPIFulfillmentCheckResult(kpiFulfillmentCheckResultPersistResult.GetPayload()):
			default:
			}
		} else if kpiFulfillmentCheckResultPersistError := kpiFulfillmentCheckResultPersistResult.GetError(); !errors.Is(kpiFulfillmentCheckResultPersistError, dbClient.ErrOperationWouldLeadToForeignKeyIntegrityBreach) {
			return fmt.Errorf("failed to persist a KPI fulfillment check result with KPI definition ID = %d, SD instance UID = %s and fulfillment status = %t: %w", kpiDefinitionID, sdInstanceUID, fulfilled, kpiFulfillmentCheckResultPersistError)
		}
		return nil
	})
}

func EnqueueMessageRepresentingCurrentSDTypeConfiguration() {
	sharedUtils.TerminateOnError(func() error {
		sdTypesLoadResult := dbClient.GetRelationalDatabaseClientInstance().LoadSDTypes()
		if sdTypesLoadResult.IsFailure() {
			return sdTypesLoadResult.GetError()
		}
		sdTypeDenotations := sharedUtils.Map[dllModel.SDType, string](sdTypesLoadResult.GetPayload(), func(sdType dllModel.SDType) string {
			return sdType.Denotation
		})
		sdTypeDenotationsJSONSerializationResult := sharedUtils.SerializeToJSON(sdTypeDenotations)
		if sdTypeDenotationsJSONSerializationResult.IsFailure() {
			return sdTypeDenotationsJSONSerializationResult.GetError()
		}
		return getRabbitMQClient().PublishJSONMessage(sharedUtils.NewEmptyOptional[string](), sharedUtils.NewOptionalOf(sharedConstants.SetOfSDTypesUpdatesQueueName), sdTypeDenotationsJSONSerializationResult.GetPayload())
	}(), "[ISC] Failed to enqueue RabbitMQ messages representing current SD type configuration")
}

func EnqueueMessageRepresentingCurrentSDInstanceConfiguration() {
	sharedUtils.TerminateOnError(func() error {
		sdInstancesLoadResult := dbClient.GetRelationalDatabaseClientInstance().LoadSDInstances()
		if sdInstancesLoadResult.IsFailure() {
			return sdInstancesLoadResult.GetError()
		}
		sdInstancesInfo := sharedUtils.Map[dllModel.SDInstance, sharedModel.SDInstanceInfo](sdInstancesLoadResult.GetPayload(), func(sdInstance dllModel.SDInstance) sharedModel.SDInstanceInfo {
			return sharedModel.SDInstanceInfo{
				SDInstanceUID:   sdInstance.UID,
				ConfirmedByUser: sdInstance.ConfirmedByUser,
			}
		})
		sdInstancesInfoJSONSerializationResult := sharedUtils.SerializeToJSON(sdInstancesInfo)
		if sdInstancesInfoJSONSerializationResult.IsFailure() {
			return sdInstancesInfoJSONSerializationResult.GetError()
		}
		return getRabbitMQClient().PublishJSONMessage(sharedUtils.NewEmptyOptional[string](), sharedUtils.NewOptionalOf(sharedConstants.SetOfSDInstancesUpdatesQueueName), sdInstancesInfoJSONSerializationResult.GetPayload())
	}(), "[ISC] Failed to enqueue RabbitMQ messages representing current SD instance configuration")
}

func EnqueueMessageRepresentingCurrentKPIDefinitionConfiguration() {
	sharedUtils.TerminateOnError(func() error {
		kpiDefinitionsLoadResult := dbClient.GetRelationalDatabaseClientInstance().LoadKPIDefinitions()
		if kpiDefinitionsLoadResult.IsFailure() {
			return kpiDefinitionsLoadResult.GetError()
		}
		kpiDefinitions := kpiDefinitionsLoadResult.GetPayload()
		kpiDefinitionsBySDTypeDenotationMap := make(sharedModel.KPIConfigurationUpdateISCMessage)
		for _, kpiDefinition := range kpiDefinitions {
			sdTypeSpecification := kpiDefinition.SDTypeSpecification
			if _, exists := kpiDefinitionsBySDTypeDenotationMap[sdTypeSpecification]; !exists {
				kpiDefinitionsBySDTypeDenotationMap[sdTypeSpecification] = make([]sharedModel.KPIDefinition, 0)
			}
			kpiDefinitionsBySDTypeDenotationMap[sdTypeSpecification] = append(kpiDefinitionsBySDTypeDenotationMap[sdTypeSpecification], kpiDefinition)
		}
		jsonSerializationResult := sharedUtils.SerializeToJSON(kpiDefinitionsBySDTypeDenotationMap)
		if jsonSerializationResult.IsFailure() {
			return jsonSerializationResult.GetError()
		}
		return getRabbitMQClient().PublishJSONMessage(sharedUtils.NewOptionalOf(sharedConstants.MainFanoutExchangeName), sharedUtils.NewEmptyOptional[string](), jsonSerializationResult.GetPayload())
	}(), "[ISC] Failed to enqueue RabbitMQ messages representing current KPI definition configuration")
}

func EnqueueMessagesRepresentingCurrentSystemConfiguration() {
	EnqueueMessageRepresentingCurrentSDTypeConfiguration()
	EnqueueMessageRepresentingCurrentSDInstanceConfiguration()
	EnqueueMessageRepresentingCurrentKPIDefinitionConfiguration()
}
