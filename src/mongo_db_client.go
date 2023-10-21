package main

import (
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	connectionString = "mongodb://admin:password@localhost:27017"
	databaseName     = "main"
	collectionName   = "kpi_definitions"
)

type MongoDBClient interface {
	ConnectToDatabase() error
	InitializeDatabase() error
	ObtainRootKPIDefinitionsForTheGivenDeviceType(deviceType string) ([]RootKPIDefinition, error)
}

type mongoDBClientImpl struct {
	client *mongo.Client
}

// NewMongoDBClient is a constructor-like function that returns an instance of MongoDBClient
func NewMongoDBClient() MongoDBClient {
	return &mongoDBClientImpl{client: nil}
}

func (m *mongoDBClientImpl) ConnectToDatabase() error {

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(connectionString).SetServerAPIOptions(serverAPIOptions)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting to database:", err)
		return err
	}

	m.client = client
	return nil
}

func transformNodeIntoMongoDBNode(node Node) (bson.M, error) {
	switch typedNode := node.(type) {
	case *StringEqualitySubKPIDefinitionNode:
		return bson.M{
			"device_parameter_specification": typedNode.DeviceParameterSpecification,
			"sub_kpi_definition_type":        "string_equality",
			"reference_value":                typedNode.ReferenceValue,
		}, nil
	case *NumericLessThanSubKPIDefinitionNode:
		return bson.M{
			"device_parameter_specification": typedNode.DeviceParameterSpecification,
			"sub_kpi_definition_type":        "numeric_less_than",
			"reference_value":                typedNode.ReferenceValue,
		}, nil
	case *NumericGreaterThanSubKPIDefinitionNode:
		return bson.M{
			"device_parameter_specification": typedNode.DeviceParameterSpecification,
			"sub_kpi_definition_type":        "numeric_greater_than",
			"reference_value":                typedNode.ReferenceValue,
		}, nil
	case *NumericEqualitySubKPIDefinitionNode:
		return bson.M{
			"device_parameter_specification": typedNode.DeviceParameterSpecification,
			"sub_kpi_definition_type":        "numeric_equality",
			"reference_value":                typedNode.ReferenceValue,
		}, nil
	case *NumericInRangeSubKPIDefinitionNode:
		return bson.M{
			"device_parameter_specification": typedNode.DeviceParameterSpecification,
			"sub_kpi_definition_type":        "numeric_in_range",
			"lower_boundary_value":           typedNode.LowerBoundaryValue,
			"upper_boundary_value":           typedNode.UpperBoundaryValue,
		}, nil
	case *BooleanEqualitySubKPIDefinitionNode:
		return bson.M{
			"device_parameter_specification": typedNode.DeviceParameterSpecification,
			"sub_kpi_definition_type":        "boolean_equality",
			"reference_value":                typedNode.ReferenceValue,
		}, nil
	case *LogicalOperatorNode:
		mongoDBChildNodes := make([]bson.M, len(typedNode.ChildNodes))
		for childNodeIndex, childNode := range typedNode.ChildNodes {
			mongoDBChildNode, err := transformNodeIntoMongoDBNode(childNode)
			if err != nil {
				return nil, err
			}
			mongoDBChildNodes[childNodeIndex] = mongoDBChildNode
		}
		return bson.M{
			"type":        typedNode.Type,
			"child_nodes": mongoDBChildNodes,
		}, nil
	default:
		return nil, errors.New("unknown node type")
	}
}

func transformRootKPIDefinitionIntoMongoDBRootKPIDefinition(rootKPIDefinition RootKPIDefinition) (MongoDBRootKPIDefinition, error) {

	mongoDBNodeMap, err := transformNodeIntoMongoDBNode(rootKPIDefinition.DefinitionRoot)
	if err != nil {
		return MongoDBRootKPIDefinition{}, err
	}

	mongoDBNodeRaw, err := bson.Marshal(mongoDBNodeMap)
	if err != nil {
		return MongoDBRootKPIDefinition{}, err
	}

	return MongoDBRootKPIDefinition{
		DeviceTypeSpecification:  rootKPIDefinition.DeviceTypeSpecification,
		HumanReadableDescription: rootKPIDefinition.HumanReadableDescription,
		DefinitionRoot:           mongoDBNodeRaw,
	}, nil
}

func (m *mongoDBClientImpl) InitializeDatabase() error {

	var err error

	databaseExists := false
	var databaseNames []string
	databaseNames, err = m.client.ListDatabaseNames(context.TODO(), bson.D{})
	if err != nil {
		log.Printf("Failed to list databases: %v\n", err)
		return err
	}
	for _, dbName := range databaseNames {
		if dbName == databaseName {
			databaseExists = true
			break
		}
	}

	if !databaseExists {
		db := m.client.Database(databaseName)
		err = db.CreateCollection(context.TODO(), collectionName)
		if err != nil {
			log.Printf("Failed to create collection: %v\n", err)
			return err
		}
		log.Println("Both database and collection created successfully, inserting data...")

		collection := db.Collection(collectionName)

		for _, rootKPIDefinitionObject := range GetRootKPIDefinitions() {
			var mongoDBRootKPIDefinitionObject MongoDBRootKPIDefinition
			mongoDBRootKPIDefinitionObject, err = transformRootKPIDefinitionIntoMongoDBRootKPIDefinition(rootKPIDefinitionObject)
			if err != nil {
				log.Printf("Data insertion failed: %v\n", err)
				return err
			}
			_, err := collection.InsertOne(context.TODO(), mongoDBRootKPIDefinitionObject)
			if err != nil {
				log.Printf("Data insertion failed: %v\n", err)
				return err
			}
		}
		log.Println("Data inserted successfully")

	} else {
		log.Println("Database already exists, skipping data insertion...")
	}

	return nil
}

func transformMongoDBNodeToNode(mongoDBNodeRaw bson.Raw) (Node, error) {

	var doc bson.M
	if err := bson.Unmarshal(mongoDBNodeRaw, &doc); err != nil {
		return nil, err
	}

	subKPIDefinitionType, ok := doc["sub_kpi_definition_type"].(string)
	if ok { // sub-KPI definition node

		deviceParameterSpecification, _ := doc["device_parameter_specification"].(string)

		switch subKPIDefinitionType {
		case "string_equality":
			referenceValue, _ := doc["reference_value"].(string)
			return &StringEqualitySubKPIDefinitionNode{
				SubKPIDefinitionBaseNode: SubKPIDefinitionBaseNode{
					DeviceParameterSpecification: deviceParameterSpecification,
				},
				ReferenceValue: referenceValue,
			}, nil
		case "numeric_less_than":
			referenceValue, _ := doc["reference_value"].(float64)
			return &NumericLessThanSubKPIDefinitionNode{
				SubKPIDefinitionBaseNode: SubKPIDefinitionBaseNode{
					DeviceParameterSpecification: deviceParameterSpecification,
				},
				ReferenceValue: referenceValue,
			}, nil
		case "numeric_greater_than":
			referenceValue, _ := doc["reference_value"].(float64)
			return &NumericGreaterThanSubKPIDefinitionNode{
				SubKPIDefinitionBaseNode: SubKPIDefinitionBaseNode{
					DeviceParameterSpecification: deviceParameterSpecification,
				},
				ReferenceValue: referenceValue,
			}, nil
		case "numeric_equality":
			referenceValue, _ := doc["reference_value"].(float64)
			return &NumericEqualitySubKPIDefinitionNode{
				SubKPIDefinitionBaseNode: SubKPIDefinitionBaseNode{
					DeviceParameterSpecification: deviceParameterSpecification,
				},
				ReferenceValue: referenceValue,
			}, nil
		case "numeric_in_range":
			lowerBoundaryValue, _ := doc["lower_boundary_value"].(float64)
			upperBoundaryValue, _ := doc["upper_boundary_value"].(float64)
			return &NumericInRangeSubKPIDefinitionNode{
				SubKPIDefinitionBaseNode: SubKPIDefinitionBaseNode{
					DeviceParameterSpecification: deviceParameterSpecification,
				},
				LowerBoundaryValue: lowerBoundaryValue,
				UpperBoundaryValue: upperBoundaryValue,
			}, nil
		case "boolean_equality":
			referenceValue, _ := doc["reference_value"].(bool)
			return &BooleanEqualitySubKPIDefinitionNode{
				SubKPIDefinitionBaseNode: SubKPIDefinitionBaseNode{
					DeviceParameterSpecification: deviceParameterSpecification,
				},
				ReferenceValue: referenceValue,
			}, nil
		}

	} else { // logical operator node

		typeValue, _ := doc["type"].(string)
		var logicalType LogicalOperatorNodeType
		switch typeValue {
		case "AND":
			logicalType = AND
		case "OR":
			logicalType = OR
		case "NOR":
			logicalType = NOR
		}

		/** FIXME: Completely broken...
		childNodesRaw, ok := doc["child_nodes"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("failed to cast child_nodes to []interface{}")
		}

		childNodes := make([]Node, len(childNodesRaw))
		for index, childNodeRaw := range childNodesRaw {
			// Attempt to cast childNodeRaw to bson.M
			childNodeMap, ok := childNodeRaw.(bson.M)
			if !ok {
				return nil, fmt.Errorf("failed to cast childNodeRaw to bson.M at index %d", index)
			}

			// Convert bson.M to bson.Raw
			childNodeBytes, err := bson.Marshal(childNodeMap)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal childNodeMap to bytes at index %d: %w", index, err)
			}
			var childNodeBsonRaw bson.Raw
			err = bson.Unmarshal(childNodeBytes, &childNodeBsonRaw)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal bytes to bson.Raw at index %d: %w", index, err)
			}

			// Now you have childNodeBsonRaw as bson.Raw, proceed with your transformation
			childNode, err := transformMongoDBNodeToNode(childNodeBsonRaw)
			if err != nil {
				return nil, fmt.Errorf("failed to transform MongoDBNode to Node at index %d: %w", index, err)
			}

			childNodes[index] = childNode
		}
		*/

		return &LogicalOperatorNode{
			Type:       logicalType,
			ChildNodes: []Node{}, // FIXME: Just a stupid placeholder here
		}, nil
	}

	return nil, errors.New("unknown node type")
}

func transformMongoDBRootKPIDefinitionToRootKPIDefinition(mongoDBRootKPIDefinition MongoDBRootKPIDefinition) (RootKPIDefinition, error) {

	node, err := transformMongoDBNodeToNode(mongoDBRootKPIDefinition.DefinitionRoot)
	return RootKPIDefinition{
		DeviceTypeSpecification:  mongoDBRootKPIDefinition.DeviceTypeSpecification,
		HumanReadableDescription: mongoDBRootKPIDefinition.HumanReadableDescription,
		DefinitionRoot:           node,
	}, err
}

func (m *mongoDBClientImpl) ObtainRootKPIDefinitionsForTheGivenDeviceType(deviceType string) ([]RootKPIDefinition, error) {

	collection := m.client.Database(databaseName).Collection(collectionName)

	var mongoDBRootKPIDefinitions []MongoDBRootKPIDefinition
	filter := bson.D{{"device_type_specification", deviceType}}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.TODO()) {

		var rawDoc bson.Raw
		if err := cursor.Decode(&rawDoc); err != nil {
			return nil, err
		}

		var mongoDBRootKPIDefinition MongoDBRootKPIDefinition
		if err := bson.Unmarshal(rawDoc, &mongoDBRootKPIDefinition); err != nil {
			return nil, err
		}

		mongoDBRootKPIDefinitions = append(mongoDBRootKPIDefinitions, mongoDBRootKPIDefinition)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	if err := cursor.Close(context.TODO()); err != nil {
		return nil, err
	}

	var rootKPIDefinitions []RootKPIDefinition

	for _, mongoDBRootKPIDefinition := range mongoDBRootKPIDefinitions {
		var rootKPIDefinitionObject RootKPIDefinition
		rootKPIDefinitionObject, err = transformMongoDBRootKPIDefinitionToRootKPIDefinition(mongoDBRootKPIDefinition)
		if err != nil {
			log.Printf("Data loading failed: %v\n", err)
			return nil, err
		}
		rootKPIDefinitions = append(rootKPIDefinitions, rootKPIDefinitionObject)
	}

	return rootKPIDefinitions, nil
}
