package main

const (
	RootKPIDefinitionTableName    = "root_kpi_definitions"
	GenericKPINodeTableName       = "generic_kpi_nodes"
	LogicalOperatorNodeTableName  = "logical_operator_nodes"
	SubKPIDefinitionNodeTableName = "sub_kpi_definition_nodes"

	UserDefinedDeviceTypesTableName = "user_defined_device_types"
	DeviceTypeParametersTableName   = "device_type_parameters"
)

type RootKPIDefinitionEntity struct {
	ID                       uint32                `gorm:"column:id;primarykey;not null"`
	DeviceTypeSpecification  string                `gorm:"column:device_type_specification;not null"`
	HumanReadableDescription string                `gorm:"column:human_readable_description;not null"`
	DefinitionRootNodeID     *uint32               `gorm:"column:definition_root_node_id;not null"`
	DefinitionRootNode       *GenericKPINodeEntity `gorm:"foreignKey:DefinitionRootNodeID"`
}

func (RootKPIDefinitionEntity) TableName() string {

	return RootKPIDefinitionTableName
}

type GenericKPINodeEntity struct {
	ID           uint32                `gorm:"column:id;primarykey;not null"`
	ParentNodeID *uint32               `gorm:"column:parent_node_id"`
	ParentNode   *GenericKPINodeEntity `gorm:"foreignKey:ParentNodeID"`
}

func (GenericKPINodeEntity) TableName() string {

	return GenericKPINodeTableName
}

type LogicalOperatorNodeEntity struct {
	NodeID *uint32                 `gorm:"column:node_id;primarykey;not null"`
	Node   *GenericKPINodeEntity   `gorm:"foreignKey:NodeID"`
	Type   LogicalOperatorNodeType `gorm:"column:type;not null"`
}

func (LogicalOperatorNodeEntity) TableName() string {

	return LogicalOperatorNodeTableName
}

type SubKPIDefinitionNodeEntity struct {
	NodeID                       *uint32               `gorm:"column:node_id;primarykey;not null"`
	Node                         *GenericKPINodeEntity `gorm:"foreignKey:NodeID"`
	DeviceParameterSpecification string                `gorm:"column:device_parameter_specification;not null"`
	Type                         string                `gorm:"column:type;not null"`
	FirstNumericReferenceValue   *float64              `gorm:"column:first_numeric_reference_value"`
	SecondNumericReferenceValue  *float64              `gorm:"column:second_numeric_reference_value"`
	StringReferenceValue         *string               `gorm:"column:string_reference_value"`
	BooleanReferenceValue        *bool                 `gorm:"column:boolean_reference_value"`
}

func (SubKPIDefinitionNodeEntity) TableName() string {

	return SubKPIDefinitionNodeTableName
}

type UserDefinedDeviceTypesEntity struct {
	ID         uint32                       `gorm:"column:id;primarykey;not null"`
	Denotation string                       `gorm:"column:denotation;not null"`
	Parameters []DeviceTypeParametersEntity `gorm:"foreignKey:UserDefinedDeviceTypeID"`
}

func (UserDefinedDeviceTypesEntity) TableName() string {

	return UserDefinedDeviceTypesTableName
}

type DeviceTypeParametersEntity struct {
	ID                      uint32 `gorm:"column:id;primarykey;not null"`
	UserDefinedDeviceTypeID uint32 `gorm:"column:user_defined_device_type_id;not null"`
	Name                    string `gorm:"column:name;not null"`
	Type                    string `gorm:"column:type;not null"`
}

func (DeviceTypeParametersEntity) TableName() string {

	return DeviceTypeParametersTableName
}
