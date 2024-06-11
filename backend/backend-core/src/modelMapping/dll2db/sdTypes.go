package dll2db

import (
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dbModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

func ToDBModelEntitySDType(sdType dllModel.SDType) dbModel.SDTypeEntity {
	return dbModel.SDTypeEntity{
		ID:         sdType.ID.GetPayloadOrDefault(0),
		Denotation: sdType.Denotation,
		Parameters: sharedUtils.Map(sdType.Parameters, func(sdParameter dllModel.SDParameter) dbModel.SDParameterEntity {
			return dbModel.SDParameterEntity{
				ID:         sdParameter.ID.GetPayloadOrDefault(0),
				Denotation: sdParameter.Denotation,
				Type:       string(sdParameter.Type),
			}
		}),
	}
}
