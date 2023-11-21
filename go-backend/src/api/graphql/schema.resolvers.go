package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.40

import (
	"bp-bures-SfPDfSD/src/api/graphql/generated"
	"bp-bures-SfPDfSD/src/api/graphql/model"
	"bp-bures-SfPDfSD/src/mapping/api2dto"
	"bp-bures-SfPDfSD/src/mapping/dto2api"
	"context"
	"github.com/thoas/go-funk"
	"strconv"
)

// CreateNewUserDefinedDeviceType is the resolver for the createNewUserDefinedDeviceType field.
func (r *mutationResolver) CreateNewUserDefinedDeviceType(ctx context.Context, input model.NewUserDefinedDeviceTypeInput) (*model.UserDefinedDeviceType, error) {
	dto := api2dto.MapNewUserDefinedDeviceTypeInputToUserDefinedDeviceTypeDTO(input)
	id, err := r.rdbClient.InsertUserDefinedDeviceType(dto)
	if err != nil {
		return nil, err
	}

	userDefinedDeviceTypeDTO, err := r.rdbClient.ObtainUserDefinedDeviceTypeByID(id)
	if err != nil {
		return nil, err
	}

	return dto2api.MapUserDefinedDeviceTypeDTOToUserDefinedDeviceType(userDefinedDeviceTypeDTO)
}

// DeleteUserDefinedDeviceType is the resolver for the deleteUserDefinedDeviceType field.
func (r *mutationResolver) DeleteUserDefinedDeviceType(ctx context.Context, input string) (*bool, error) {
	id, err := strconv.ParseInt(input, 10, 32)
	if err != nil {
		return nil, err
	}

	if err := r.rdbClient.DeleteUserDefinedDeviceType(uint32(id)); err != nil {
		return nil, err
	}

	result := true
	return &result, nil
}

// SingleUserDefinedDeviceType is the resolver for the singleUserDefinedDeviceType field.
func (r *queryResolver) SingleUserDefinedDeviceType(ctx context.Context, input string) (*model.UserDefinedDeviceType, error) {
	id, err := strconv.ParseInt(input, 10, 32)
	if err != nil {
		return nil, err
	}

	userDefinedDeviceTypeDTO, err := r.rdbClient.ObtainUserDefinedDeviceTypeByID(uint32(id))
	if err != nil {
		return nil, err
	}

	return dto2api.MapUserDefinedDeviceTypeDTOToUserDefinedDeviceType(userDefinedDeviceTypeDTO)
}

// UserDefinedDeviceTypes is the resolver for the userDefinedDeviceTypes field.
func (r *queryResolver) UserDefinedDeviceTypes(ctx context.Context) ([]*model.UserDefinedDeviceType, error) {
	userDefinedDeviceTypeDTOs, err := r.rdbClient.ObtainAllUserDefinedDeviceTypes()
	if err != nil {
		return nil, err
	}

	return dto2api.MapUserDefinedDeviceTypeDTOsToUserDefinedDeviceTypes(userDefinedDeviceTypeDTOs)
}

// Devices is the resolver for the devices field.
func (r *queryResolver) Devices(ctx context.Context) ([]*model.Device, error) {

	deviceDTOs, err := r.rdbClient.ObtainAllDevices()
	if err != nil {
		return nil, err
	}

	return funk.Map(deviceDTOs, dto2api.MapDeviceDTOToDevice).([]*model.Device), nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
