package graphql

import (
	"github.com/graphql-go/graphql"
	"log"
)

func SetupGraphQLSchema() (graphql.Schema, error) {

	deviceParameterTypeEnum := graphql.NewEnum(graphql.EnumConfig{
		Name: "DeviceParameterType",
		Values: graphql.EnumValueConfigMap{
			"NUMBER": &graphql.EnumValueConfig{
				Value: "number",
			},
			"STRING": &graphql.EnumValueConfig{
				Value: "string",
			},
			"BOOLEAN": &graphql.EnumValueConfig{
				Value: "boolean",
			},
		},
	})

	deviceParameterDefinitionType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "DeviceParameterDefinition",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.NewNonNull(graphql.ID),
				},
				"name": &graphql.Field{
					Type: graphql.NewNonNull(graphql.String),
				},
				"type": &graphql.Field{
					Type: graphql.NewNonNull(deviceParameterTypeEnum),
				},
			},
		},
	)

	userDefinedDeviceType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "UserDefinedDevice",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.NewNonNull(graphql.ID),
				},
				"denotation": &graphql.Field{
					Type: graphql.NewNonNull(graphql.String),
				},
				"parameters": &graphql.Field{
					Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(deviceParameterDefinitionType))),
				},
			},
		},
	)

	schemaConfig := graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"singleUserDefinedDeviceType": &graphql.Field{
					Type: userDefinedDeviceType,
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.ID),
						},
					},
					Resolve: singleUserDefinedDeviceTypeQueryResolverFunction,
				},
				"userDefinedDeviceTypes": &graphql.Field{
					Type:    graphql.NewList(userDefinedDeviceType),
					Resolve: userDefinedDeviceTypesQueryResolverFunction,
				},
			},
		}),
	}

	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Printf("Failed to generate GraphQL schema, error: %v", err)
		return schema, err
	}
	return schema, nil
}
