package dto2db

import (
	"bp-bures-SfPDfSD/src/dto"
	"bp-bures-SfPDfSD/src/persistence/rdb/schema"
	"github.com/thoas/go-funk"
	"strings"
)

func MapDeviceDTOToDeviceEntity(d dto.DeviceDTO) schema.DeviceEntity {

	return schema.DeviceEntity{
		UID:  d.UID,
		Name: d.Name,
		DeviceType: schema.UserDefinedDeviceTypeEntity{
			Denotation: d.DeviceType.Denotation,
			Parameters: funk.Map(d.DeviceType.Parameters, func(d dto.UserDefinedDeviceTypeParameterDTO) schema.DeviceTypeParameterEntity {
				return schema.DeviceTypeParameterEntity{
					Name: d.Name,
					Type: strings.ToLower(string(d.Type)),
				}
			}).([]schema.DeviceTypeParameterEntity),
		},
	}
}