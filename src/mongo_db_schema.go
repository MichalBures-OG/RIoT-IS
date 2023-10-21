package main

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoDBRootKPIDefinition struct {
	ID                       primitive.ObjectID `bson:"_id,omitempty"`
	DeviceTypeSpecification  string             `bson:"device_type_specification"`
	HumanReadableDescription string             `bson:"human_readable_description"`
	DefinitionRoot           bson.Raw           `bson:"definition_root"`
}

type MongoDBSubKPIDefinitionBaseNode struct {
	DeviceParameterSpecification string `bson:"device_parameter_specification"`
	SubKPIDefinitionType         string `bson:"sub_kpi_definition_type"`
}

type MongoDBStringEqualitySubKPIDefinitionNode struct {
	MongoDBSubKPIDefinitionBaseNode `bson:",inline"`
	ReferenceValue                  string `bson:"reference_value"`
}

type MongoDBNumericLessThanSubKPIDefinitionNode struct {
	MongoDBSubKPIDefinitionBaseNode `bson:",inline"`
	ReferenceValue                  float64 `bson:"reference_value"`
}

type MongoDBNumericGreaterThanSubKPIDefinitionNode struct {
	MongoDBSubKPIDefinitionBaseNode `bson:",inline"`
	ReferenceValue                  float64 `bson:"reference_value"`
}

type MongoDBNumericEqualitySubKPIDefinitionNode struct {
	MongoDBSubKPIDefinitionBaseNode `bson:",inline"`
	ReferenceValue                  float64 `bson:"reference_value"`
}

type MongoDBNumericInRangeSubKPIDefinitionNode struct {
	MongoDBSubKPIDefinitionBaseNode `bson:",inline"`
	LowerBoundaryValue              float64 `bson:"lower_boundary_value"`
	UpperBoundaryValue              float64 `bson:"upper_boundary_value"`
}

type MongoDBBooleanEqualitySubKPIDefinitionNode struct {
	MongoDBSubKPIDefinitionBaseNode `bson:",inline"`
	ReferenceValue                  bool `bson:"reference_value"`
}

type MongoDBLogicalOperatorNode struct {
	Type       LogicalOperatorNodeType `bson:"type"`
	ChildNodes []bson.Raw              `bson:"child_nodes"`
}
