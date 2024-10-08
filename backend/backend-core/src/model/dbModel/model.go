package dbModel

type KPIDefinitionEntity struct {
	ID                                         uint32                                      `gorm:"column:id;primaryKey"`
	UserIdentifier                             string                                      `gorm:"column:user_identifier;not null"`
	RootNodeID                                 *uint32                                     `gorm:"column:root_node_id;not null"`
	RootNode                                   *KPINodeEntity                              `gorm:"foreignKey:RootNodeID"`
	SDTypeID                                   uint32                                      `gorm:"column:sd_type_id;not null"`
	SDType                                     SDTypeEntity                                `gorm:"foreignKey:SDTypeID;constraint:OnDelete:CASCADE"`
	SDInstanceMode                             string                                      `gorm:"column:sd_instance_mode;not null"`
	KPIFulfillmentCheckResults                 []KPIFulfillmentCheckResultEntity           `gorm:"foreignKey:KPIDefinitionID;constraint:OnDelete:CASCADE"`
	SDInstanceKPIDefinitionRelationshipRecords []SDInstanceKPIDefinitionRelationshipEntity `gorm:"foreignKey:KPIDefinitionID;constraint:OnDelete:CASCADE"`
}

func (KPIDefinitionEntity) TableName() string {
	return "kpi_definitions"
}

type KPINodeEntity struct {
	ID           uint32         `gorm:"column:id;primaryKey;not null"`
	ParentNodeID *uint32        `gorm:"column:parent_node_id"`
	ParentNode   *KPINodeEntity `gorm:"foreignKey:ParentNodeID;constraint:OnDelete:CASCADE"`
}

func (KPINodeEntity) TableName() string {
	return "kpi_nodes"
}

type LogicalOperationKPINodeEntity struct {
	NodeID *uint32        `gorm:"column:node_id;primaryKey;not null"`
	Node   *KPINodeEntity `gorm:"foreignKey:NodeID;constraint:OnDelete:CASCADE"`
	Type   string         `gorm:"column:type;not null"`
}

func (LogicalOperationKPINodeEntity) TableName() string {
	return "logical_operation_kpi_nodes"
}

type AtomKPINodeEntity struct {
	NodeID                *uint32           `gorm:"column:node_id;primaryKey;not null"`
	Node                  *KPINodeEntity    `gorm:"foreignKey:NodeID;constraint:OnDelete:CASCADE"`
	SDParameterID         uint32            `gorm:"column:sd_parameter_id;not null"`
	SDParameter           SDParameterEntity `gorm:"foreignKey:SDParameterID;constraint:OnDelete:CASCADE"`
	Type                  string            `gorm:"column:type;not null"`
	StringReferenceValue  *string           `gorm:"column:string_reference_value"`
	BooleanReferenceValue *bool             `gorm:"column:boolean_reference_value"`
	NumericReferenceValue *float64          `gorm:"column:numeric_reference_value"`
}

func (AtomKPINodeEntity) TableName() string {
	return "atom_kpi_nodes"
}

type SDTypeEntity struct {
	ID         uint32              `gorm:"column:id;primaryKey;not null"`
	Denotation string              `gorm:"column:denotation;not null;index"` // Denotation is a separately indexed field
	Parameters []SDParameterEntity `gorm:"foreignKey:SDTypeID;constraint:OnDelete:CASCADE"`
}

func (SDTypeEntity) TableName() string {
	return "sd_types"
}

type SDParameterEntity struct {
	ID         uint32 `gorm:"column:id;primaryKey;not null"`
	SDTypeID   uint32 `gorm:"column:sd_type_id;not null"`
	Denotation string `gorm:"column:denotation;not null"`
	Type       string `gorm:"column:type;not null"`
}

func (SDParameterEntity) TableName() string {
	return "sd_parameters"
}

type SDInstanceEntity struct {
	ID                                         uint32                                      `gorm:"column:id;primaryKey;not null"`
	UID                                        string                                      `gorm:"column:uid;not null;index"` // UID is a separately indexed field
	ConfirmedByUser                            bool                                        `gorm:"column:confirmed_by_user;not null"`
	UserIdentifier                             string                                      `gorm:"column:user_identifier;not null"`
	SDTypeID                                   uint32                                      `gorm:"column:sd_type_id"`
	SDType                                     SDTypeEntity                                `gorm:"foreignKey:SDTypeID;constraint:OnDelete:CASCADE"`
	GroupMembershipRecords                     []SDInstanceGroupMembershipEntity           `gorm:"foreignKey:SDInstanceID;constraint:OnDelete:CASCADE"`
	KPIFulfillmentCheckResults                 []KPIFulfillmentCheckResultEntity           `gorm:"foreignKey:SDInstanceID;constraint:OnDelete:CASCADE"`
	SDInstanceKPIDefinitionRelationshipRecords []SDInstanceKPIDefinitionRelationshipEntity `gorm:"foreignKey:SDInstanceID;constraint:OnDelete:CASCADE"`
}

func (SDInstanceEntity) TableName() string {
	return "sd_instances"
}

type KPIFulfillmentCheckResultEntity struct {
	KPIDefinitionID uint32 `gorm:"column:kpi_definition_id;primaryKey;not null"`
	SDInstanceID    uint32 `gorm:"column:sd_instance_id;primaryKey;not null"`
	Fulfilled       bool   `gorm:"column:fulfilled;not null"`
}

func (KPIFulfillmentCheckResultEntity) TableName() string {
	return "kpi_fulfillment_check_results"
}

type SDInstanceGroupEntity struct {
	ID                     uint32                            `gorm:"column:id;primaryKey;not null"`
	UserIdentifier         string                            `gorm:"column:user_identifier;not null"`
	GroupMembershipRecords []SDInstanceGroupMembershipEntity `gorm:"foreignKey:SDInstanceGroupID;constraint:OnDelete:CASCADE"`
}

func (SDInstanceGroupEntity) TableName() string {
	return "sd_instance_groups"
}

type SDInstanceGroupMembershipEntity struct {
	SDInstanceGroupID uint32 `gorm:"column:sd_instance_group_id;primaryKey;not null;index"` // SDInstanceGroupID is a separately indexed field
	SDInstanceID      uint32 `gorm:"column:sd_instance_id;primaryKey;not null"`
}

func (SDInstanceGroupMembershipEntity) TableName() string {
	return "sd_instance_group_membership"
}

type SDInstanceKPIDefinitionRelationshipEntity struct {
	KPIDefinitionID uint32 `gorm:"column:kpi_definition_id;primaryKey;not null"`
	SDInstanceID    uint32 `gorm:"column:sd_instance_id;primaryKey;not null"`
	SDInstanceUID   string `gorm:"column:sd_instance_uid;not null"`
}
