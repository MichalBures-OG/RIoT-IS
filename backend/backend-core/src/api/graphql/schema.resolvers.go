package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"fmt"

	"github.com/MichalBures-OG/bp-bures-SfPDfSD-backend-core/src/api/graphql/generated"
	"github.com/MichalBures-OG/bp-bures-SfPDfSD-backend-core/src/api/graphql/model"
	"github.com/MichalBures-OG/bp-bures-SfPDfSD-backend-core/src/service"
)

// CreateSDType is the resolver for the createSDType field.
func (r *mutationResolver) CreateSDType(ctx context.Context, input model.SDTypeInput) (*model.SDType, error) {
	return service.CreateSDType(input).Unwrap()
}

// DeleteSDType is the resolver for the deleteSDType field.
func (r *mutationResolver) DeleteSDType(ctx context.Context, id string) (bool, error) {
	err := service.DeleteSDType(id)
	return err == nil, err
}

// UpdateSDInstance is the resolver for the updateSDInstance field.
func (r *mutationResolver) UpdateSDInstance(ctx context.Context, id string, input model.SDInstanceUpdateInput) (*model.SDInstance, error) {
	return service.UpdateSDInstance(id, input).Unwrap()
}

// CreateKPIDefinition is the resolver for the createKPIDefinition field.
func (r *mutationResolver) CreateKPIDefinition(ctx context.Context, input model.KPIDefinitionInput) (*model.KPIDefinition, error) {
	return service.CreateKPIDefinition(input).Unwrap()
}

// SdType is the resolver for the sdType field.
func (r *queryResolver) SdType(ctx context.Context, id string) (*model.SDType, error) {
	return service.GetSDType(id).Unwrap()
}

// SdTypes is the resolver for the sdTypes field.
func (r *queryResolver) SdTypes(ctx context.Context) ([]*model.SDType, error) {
	return service.GetSDTypes().Unwrap()
}

// SdInstances is the resolver for the sdInstances field.
func (r *queryResolver) SdInstances(ctx context.Context) ([]*model.SDInstance, error) {
	return service.GetSDInstances().Unwrap()
}

// KpiDefinitions is the resolver for the kpiDefinitions field.
func (r *queryResolver) KpiDefinitions(ctx context.Context) ([]*model.KPIDefinition, error) {
	panic(fmt.Errorf("not implemented: KpiDefinitions - kpiDefinitions"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
