package dbClient

import (
	"errors"
	"fmt"
	"github.com/MichalBures-OG/bp-bures-SfPDfSD-backend-core/src/db/dbUtil"
	"github.com/MichalBures-OG/bp-bures-SfPDfSD-backend-core/src/model/dbModel"
	"github.com/MichalBures-OG/bp-bures-SfPDfSD-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-SfPDfSD-backend-core/src/modelMapping/db2dll"
	"github.com/MichalBures-OG/bp-bures-SfPDfSD-backend-core/src/modelMapping/dll2db"
	"github.com/MichalBures-OG/bp-bures-SfPDfSD-commons/src/sharedModel"
	"github.com/MichalBures-OG/bp-bures-SfPDfSD-commons/src/sharedUtils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
)

const (
	dsn = "host=postgres user=admin password=password dbname=postgres-db port=5432"
)

var (
	rdbClientInstance RelationalDatabaseClient
	once              sync.Once
)

type RelationalDatabaseClient interface {
	setup()
	PersistKPIDefinition(kpiDefinition sharedModel.KPIDefinition) sharedUtils.Result[uint32] // TODO: Is ID enough?
	LoadKPIDefinition(id uint32) sharedUtils.Result[sharedModel.KPIDefinition]
	LoadKPIDefinitions() sharedUtils.Result[[]sharedModel.KPIDefinition]
	DeleteKPIDefinition(id uint32) error
	PersistSDType(sdType dllModel.SDType) sharedUtils.Result[dllModel.SDType]
	LoadSDType(id uint32) sharedUtils.Result[dllModel.SDType]
	LoadSDTypeBasedOnDenotation(denotation string) sharedUtils.Result[dllModel.SDType]
	LoadSDTypes() sharedUtils.Result[[]dllModel.SDType]
	DeleteSDType(id uint32) error
	PersistSDInstance(sdInstance dllModel.SDInstance) sharedUtils.Result[uint32] // TODO: Is ID enough?
	LoadSDInstance(id uint32) sharedUtils.Result[dllModel.SDInstance]
	LoadSDInstanceBasedOnUID(uid string) sharedUtils.Result[sharedUtils.Optional[dllModel.SDInstance]]
	DoesSDInstanceExist(uid string) sharedUtils.Result[bool]
	LoadSDInstances() sharedUtils.Result[[]dllModel.SDInstance]
	PersistKPIFulFulfillmentCheckResult(kpiFulfillmentCheckResult dllModel.KPIFulfillmentCheckResult) error
	LoadKPIFulFulfillmentCheckResult(kpiDefinitionID uint32, sdInstanceID uint32) sharedUtils.Result[sharedUtils.Optional[dllModel.KPIFulfillmentCheckResult]]
	LoadKPIFulFulfillmentCheckResults() sharedUtils.Result[[]dllModel.KPIFulfillmentCheckResult]
	LoadSDInstanceGroups() sharedUtils.Result[[]dllModel.SDInstanceGroup]
	LoadSDInstanceGroup(id uint32) sharedUtils.Result[dllModel.SDInstanceGroup]
	PersistSDInstanceGroup(sdInstanceGroup dllModel.SDInstanceGroup) sharedUtils.Result[uint32]
	DeleteSDInstanceGroup(id uint32) error
}

type relationalDatabaseClientImpl struct {
	db *gorm.DB
	mu sync.Mutex
}

func GetRelationalDatabaseClientInstance() RelationalDatabaseClient {
	once.Do(func() {
		rdbClientInstance = new(relationalDatabaseClientImpl)
		rdbClientInstance.setup()
	})
	return rdbClientInstance
}

func (r *relationalDatabaseClientImpl) setup() {
	db, err := gorm.Open(postgres.Open(dsn), new(gorm.Config))
	sharedUtils.TerminateOnError(err, "[RDB client (GORM)]: couldn't connect to the database")
	session := new(gorm.Session)
	session.Logger = logger.Default.LogMode(logger.Silent)
	r.db = db.Session(session)
	sharedUtils.TerminateOnError(r.db.AutoMigrate(
		new(dbModel.KPIDefinitionEntity),
		new(dbModel.KPINodeEntity),
		new(dbModel.LogicalOperationKPINodeEntity),
		new(dbModel.AtomKPINodeEntity),
		new(dbModel.SDTypeEntity),
		new(dbModel.SDParameterEntity),
		new(dbModel.SDInstanceEntity),
		new(dbModel.KPIFulfillmentCheckResultEntity),
		new(dbModel.SDInstanceGroupEntity),
		new(dbModel.SDInstanceGroupMembershipEntity),
	), "[RDB client (GORM)]: auto-migration failed")
}

func (r *relationalDatabaseClientImpl) PersistKPIDefinition(kpiDefinition sharedModel.KPIDefinition) sharedUtils.Result[uint32] {
	r.mu.Lock()
	defer r.mu.Unlock()
	kpiNodeEntity, kpiNodeEntities, logicalOperationNodeEntities, atomNodeEntities := dll2db.ToDBModelEntitiesKPIDefinition(kpiDefinition)
	kpiDefinitionEntity := dbModel.KPIDefinitionEntity{
		SDTypeID:       kpiDefinition.SDTypeID,
		UserIdentifier: kpiDefinition.UserIdentifier,
		RootNode:       kpiNodeEntity,
	}
	err := r.db.Transaction(func(tx *gorm.DB) error {
		for _, entity := range kpiNodeEntities {
			if err := dbUtil.PersistEntityIntoDB[dbModel.KPINodeEntity](tx, entity); err != nil {
				return err
			}
		}
		if err := dbUtil.PersistEntityIntoDB[dbModel.KPIDefinitionEntity](tx, &kpiDefinitionEntity); err != nil {
			return err
		}
		for _, entity := range logicalOperationNodeEntities {
			if err := dbUtil.PersistEntityIntoDB[dbModel.LogicalOperationKPINodeEntity](tx, &entity); err != nil {
				return err
			}
		}
		for _, entity := range atomNodeEntities {
			if err := dbUtil.PersistEntityIntoDB[dbModel.AtomKPINodeEntity](tx, &entity); err != nil {
				return err
			}
		}
		return nil
	})
	return sharedUtils.Ternary(err == nil, sharedUtils.NewSuccessResult[uint32](kpiDefinitionEntity.ID), sharedUtils.NewFailureResult[uint32](err))
}

func (r *relationalDatabaseClientImpl) LoadKPIDefinition(id uint32) sharedUtils.Result[sharedModel.KPIDefinition] {
	r.mu.Lock()
	defer r.mu.Unlock()
	// TODO: Implement
	return sharedUtils.NewFailureResult[sharedModel.KPIDefinition](errors.New("[RDB client (GORM)]: not implemented"))
}

func (r *relationalDatabaseClientImpl) LoadKPIDefinitions() sharedUtils.Result[[]sharedModel.KPIDefinition] {
	r.mu.Lock()
	defer r.mu.Unlock()
	kpiDefinitionEntitiesLoadResult := dbUtil.LoadEntitiesFromDB[dbModel.KPIDefinitionEntity](r.db, dbUtil.Preload("SDType"))
	if kpiDefinitionEntitiesLoadResult.IsFailure() {
		err := fmt.Errorf("failed to load KPI definition entities from the database: %w", kpiDefinitionEntitiesLoadResult.GetError())
		return sharedUtils.NewFailureResult[[]sharedModel.KPIDefinition](err)
	}
	kpiDefinitionEntities := kpiDefinitionEntitiesLoadResult.GetPayload()
	kpiNodeEntitiesLoadResult := dbUtil.LoadEntitiesFromDB[dbModel.KPINodeEntity](r.db)
	if kpiNodeEntitiesLoadResult.IsFailure() {
		err := fmt.Errorf("failed to load KPI node entities from the database: %w", kpiNodeEntitiesLoadResult.GetError())
		return sharedUtils.NewFailureResult[[]sharedModel.KPIDefinition](err)
	}
	kpiNodeEntities := kpiNodeEntitiesLoadResult.GetPayload()
	logicalOperationKPINodeEntitiesLoadResult := dbUtil.LoadEntitiesFromDB[dbModel.LogicalOperationKPINodeEntity](r.db)
	if logicalOperationKPINodeEntitiesLoadResult.IsFailure() {
		err := fmt.Errorf("failed to load logical operation KPI node entities from the database: %w", logicalOperationKPINodeEntitiesLoadResult.GetError())
		return sharedUtils.NewFailureResult[[]sharedModel.KPIDefinition](err)
	}
	logicalOperationKPINodeEntities := logicalOperationKPINodeEntitiesLoadResult.GetPayload()
	atomKPINodeEntitiesLoadResult := dbUtil.LoadEntitiesFromDB[dbModel.AtomKPINodeEntity](r.db, dbUtil.Preload("SDParameter"))
	if atomKPINodeEntitiesLoadResult.IsFailure() {
		err := fmt.Errorf("failed to load atom KPI node entities from the database: %w", atomKPINodeEntitiesLoadResult.GetError())
		return sharedUtils.NewFailureResult[[]sharedModel.KPIDefinition](err)
	}
	atomKPINodeEntities := atomKPINodeEntitiesLoadResult.GetPayload()
	kpiDefinitions := make([]sharedModel.KPIDefinition, 0, len(kpiDefinitionEntities))
	sharedUtils.ForEach(kpiDefinitionEntities, func(kpiDefinitionEntity dbModel.KPIDefinitionEntity) {
		kpiDefinition := db2dll.ToDLLModelKPIDefinition(kpiDefinitionEntity, kpiNodeEntities, logicalOperationKPINodeEntities, atomKPINodeEntities)
		kpiDefinitions = append(kpiDefinitions, kpiDefinition)
	})
	return sharedUtils.NewSuccessResult[[]sharedModel.KPIDefinition](kpiDefinitions)
}

func getIDsOfKPINodeEntitiesFormingTheKPIDefinition(g *gorm.DB, kpiDefinitionID uint32) sharedUtils.Result[[]uint32] {
	kpiDefinitionEntityLoadResult := dbUtil.LoadEntityFromDB[dbModel.KPIDefinitionEntity](g, dbUtil.Where("id = ?", kpiDefinitionID))
	if kpiDefinitionEntityLoadResult.IsFailure() {
		err := fmt.Errorf("failed to load KPI definition entity with ID = %d from the database: %w", kpiDefinitionID, kpiDefinitionEntityLoadResult.GetError())
		return sharedUtils.NewFailureResult[[]uint32](err)
	}
	kpiDefinitionEntity := kpiDefinitionEntityLoadResult.GetPayload()
	kpiNodeEntitiesLoadResult := dbUtil.LoadEntitiesFromDB[dbModel.KPINodeEntity](g)
	if kpiNodeEntitiesLoadResult.IsFailure() {
		err := fmt.Errorf("failed to load KPI node entities from the database: %w", kpiNodeEntitiesLoadResult.GetError())
		return sharedUtils.NewFailureResult[[]uint32](err)
	}
	setOfIDsOfKPINodeEntitiesFormingTheDefinition := sharedUtils.NewSetFromSlice(sharedUtils.SliceOf(*kpiDefinitionEntity.RootNodeID))
	setOfRemainingKPINodeEntities := sharedUtils.NewSetFromSlice(kpiNodeEntitiesLoadResult.GetPayload())
	for {
		nextLayerOfKPINodeEntities := sharedUtils.Filter(setOfRemainingKPINodeEntities.ToSlice(), func(kpiNodeEntity dbModel.KPINodeEntity) bool {
			parentNodeIDOptional := sharedUtils.NewOptionalFromPointer(kpiNodeEntity.ParentNodeID)
			if parentNodeIDOptional.IsEmpty() {
				return false
			}
			return setOfIDsOfKPINodeEntitiesFormingTheDefinition.Contains(parentNodeIDOptional.GetPayload())
		})
		if len(nextLayerOfKPINodeEntities) == 0 {
			break
		}
		sharedUtils.ForEach(nextLayerOfKPINodeEntities, func(kpiNodeEntity dbModel.KPINodeEntity) {
			setOfIDsOfKPINodeEntitiesFormingTheDefinition.Add(kpiNodeEntity.ID)
			setOfRemainingKPINodeEntities.Delete(kpiNodeEntity)
		})
	}
	return sharedUtils.NewSuccessResult[[]uint32](setOfIDsOfKPINodeEntitiesFormingTheDefinition.ToSlice())
}

func (r *relationalDatabaseClientImpl) DeleteKPIDefinition(id uint32) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	idsOfKPINodeEntitiesFormingTheDefinitionResult := getIDsOfKPINodeEntitiesFormingTheKPIDefinition(r.db, id)
	if idsOfKPINodeEntitiesFormingTheDefinitionResult.IsFailure() {
		return idsOfKPINodeEntitiesFormingTheDefinitionResult.GetError()
	}
	if err := dbUtil.DeleteEntitiesBasedOnSliceOfIds[dbModel.KPINodeEntity](r.db, idsOfKPINodeEntitiesFormingTheDefinitionResult.GetPayload()); err != nil {
		return err
	}
	return nil
}

func (r *relationalDatabaseClientImpl) PersistSDType(sdType dllModel.SDType) sharedUtils.Result[dllModel.SDType] {
	r.mu.Lock()
	defer r.mu.Unlock()
	sdTypeEntity := dll2db.ToDBModelEntitySDType(sdType)
	if err := dbUtil.PersistEntityIntoDB[dbModel.SDTypeEntity](r.db, &sdTypeEntity); err != nil {
		return sharedUtils.NewFailureResult[dllModel.SDType](err)
	}
	return sharedUtils.NewSuccessResult[dllModel.SDType](db2dll.ToDLLModelSDType(sdTypeEntity))
}

func loadSDType(g *gorm.DB, whereClause dbUtil.WhereClause) sharedUtils.Result[dllModel.SDType] {
	sdTypeEntityLoadResult := dbUtil.LoadEntityFromDB[dbModel.SDTypeEntity](g, dbUtil.Preload("Parameters"), whereClause)
	if sdTypeEntityLoadResult.IsFailure() {
		return sharedUtils.NewFailureResult[dllModel.SDType](sdTypeEntityLoadResult.GetError())
	}
	return sharedUtils.NewSuccessResult[dllModel.SDType](db2dll.ToDLLModelSDType(sdTypeEntityLoadResult.GetPayload()))
}

func (r *relationalDatabaseClientImpl) LoadSDType(id uint32) sharedUtils.Result[dllModel.SDType] {
	r.mu.Lock()
	defer r.mu.Unlock()
	return loadSDType(r.db, dbUtil.Where("id = ?", id))
}

func (r *relationalDatabaseClientImpl) LoadSDTypeBasedOnDenotation(denotation string) sharedUtils.Result[dllModel.SDType] {
	r.mu.Lock()
	defer r.mu.Unlock()
	return loadSDType(r.db, dbUtil.Where("denotation = ?", denotation))
}

func (r *relationalDatabaseClientImpl) LoadSDTypes() sharedUtils.Result[[]dllModel.SDType] {
	r.mu.Lock()
	defer r.mu.Unlock()
	sdTypeEntitiesLoadResult := dbUtil.LoadEntitiesFromDB[dbModel.SDTypeEntity](r.db, dbUtil.Preload("Parameters"))
	if sdTypeEntitiesLoadResult.IsFailure() {
		return sharedUtils.NewFailureResult[[]dllModel.SDType](sdTypeEntitiesLoadResult.GetError())
	}
	return sharedUtils.NewSuccessResult[[]dllModel.SDType](sharedUtils.Map(sdTypeEntitiesLoadResult.GetPayload(), db2dll.ToDLLModelSDType))
}

func (r *relationalDatabaseClientImpl) DeleteSDType(id uint32) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.db.Transaction(func(tx *gorm.DB) error {
		relatedKPIDefinitionEntitiesLoadResult := dbUtil.LoadEntitiesFromDB[dbModel.KPIDefinitionEntity](tx, dbUtil.Where("sd_type_id = ?", id))
		if relatedKPIDefinitionEntitiesLoadResult.IsFailure() {
			return fmt.Errorf("failed to load KPI definition entities related to the SD type with ID = %d from the database: %w", id, relatedKPIDefinitionEntitiesLoadResult.GetError())
		}
		if err := dbUtil.DeleteCertainEntityBasedOnId[dbModel.SDTypeEntity](tx, id); err != nil {
			return fmt.Errorf("failed to delete SD type entity with ID = %d from the database: %w", id, err)
		}
		relatedKPIDefinitionEntities := relatedKPIDefinitionEntitiesLoadResult.GetPayload()
		for _, relatedKPIDefinitionEntity := range relatedKPIDefinitionEntities {
			rootNodeID := sharedUtils.NewOptionalFromPointer(relatedKPIDefinitionEntity.RootNodeID).GetPayload()
			if err := dbUtil.DeleteCertainEntityBasedOnId[dbModel.KPINodeEntity](tx, rootNodeID); err != nil {
				return fmt.Errorf("failed to delete KPI node entity with ID = %d from the database: %w", rootNodeID, err)
			}
		}
		return nil
	})
}

func (r *relationalDatabaseClientImpl) PersistSDInstance(sdInstance dllModel.SDInstance) sharedUtils.Result[uint32] {
	r.mu.Lock()
	defer r.mu.Unlock()
	sdInstanceEntity := dll2db.ToDBModelEntitySDInstance(sdInstance)
	if err := dbUtil.PersistEntityIntoDB(r.db, &sdInstanceEntity); err != nil {
		return sharedUtils.NewFailureResult[uint32](err)
	}
	return sharedUtils.NewSuccessResult[uint32](sdInstanceEntity.ID)
}

func (r *relationalDatabaseClientImpl) LoadSDInstance(id uint32) sharedUtils.Result[dllModel.SDInstance] {
	r.mu.Lock()
	defer r.mu.Unlock()
	sdInstanceEntityLoadResult := dbUtil.LoadEntityFromDB[dbModel.SDInstanceEntity](r.db, dbUtil.Preload("SDType"), dbUtil.Where("id = ?", id))
	if sdInstanceEntityLoadResult.IsFailure() {
		return sharedUtils.NewFailureResult[dllModel.SDInstance](sdInstanceEntityLoadResult.GetError())
	}
	return sharedUtils.NewSuccessResult[dllModel.SDInstance](db2dll.ToDLLModelSDInstance(sdInstanceEntityLoadResult.GetPayload()))
}

func (r *relationalDatabaseClientImpl) LoadSDInstanceBasedOnUID(uid string) sharedUtils.Result[sharedUtils.Optional[dllModel.SDInstance]] {
	r.mu.Lock()
	defer r.mu.Unlock()
	sdInstanceEntityLoadResult := dbUtil.LoadEntityFromDB[dbModel.SDInstanceEntity](r.db, dbUtil.Preload("SDType"), dbUtil.Where("uid = ?", uid))
	if sdInstanceEntityLoadResult.IsFailure() {
		err := sdInstanceEntityLoadResult.GetError()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sharedUtils.NewSuccessResult[sharedUtils.Optional[dllModel.SDInstance]](sharedUtils.NewEmptyOptional[dllModel.SDInstance]())
		} else {
			return sharedUtils.NewFailureResult[sharedUtils.Optional[dllModel.SDInstance]](err)
		}
	}
	return sharedUtils.NewSuccessResult[sharedUtils.Optional[dllModel.SDInstance]](sharedUtils.NewOptionalOf(db2dll.ToDLLModelSDInstance(sdInstanceEntityLoadResult.GetPayload())))
}

func (r *relationalDatabaseClientImpl) DoesSDInstanceExist(uid string) sharedUtils.Result[bool] {
	r.mu.Lock()
	defer r.mu.Unlock()
	return dbUtil.DoesSuchEntityExist[dbModel.SDInstanceEntity](r.db, dbUtil.Where("uid = ?", uid))
}

func (r *relationalDatabaseClientImpl) LoadSDInstances() sharedUtils.Result[[]dllModel.SDInstance] {
	r.mu.Lock()
	defer r.mu.Unlock()
	sdInstanceEntitiesLoadResult := dbUtil.LoadEntitiesFromDB[dbModel.SDInstanceEntity](r.db, dbUtil.Preload("SDType"))
	if sdInstanceEntitiesLoadResult.IsFailure() {
		return sharedUtils.NewFailureResult[[]dllModel.SDInstance](sdInstanceEntitiesLoadResult.GetError())
	}
	return sharedUtils.NewSuccessResult[[]dllModel.SDInstance](sharedUtils.Map(sdInstanceEntitiesLoadResult.GetPayload(), db2dll.ToDLLModelSDInstance))
}

func (r *relationalDatabaseClientImpl) PersistKPIFulFulfillmentCheckResult(kpiFulfillmentCheckResult dllModel.KPIFulfillmentCheckResult) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	kpiFulfillmentCheckEntity := dll2db.ToDBModelEntityKPIFulfillmentCheckResult(kpiFulfillmentCheckResult)
	return dbUtil.PersistEntityIntoDB(r.db, &kpiFulfillmentCheckEntity)
}

func (r *relationalDatabaseClientImpl) LoadKPIFulFulfillmentCheckResult(kpiDefinitionID uint32, sdInstanceID uint32) sharedUtils.Result[sharedUtils.Optional[dllModel.KPIFulfillmentCheckResult]] {
	r.mu.Lock()
	defer r.mu.Unlock()
	kpiFulFulfillmentCheckResultEntityLoadResult := dbUtil.LoadEntityFromDB[dbModel.KPIFulfillmentCheckResultEntity](r.db, dbUtil.Where("kpi_definition_id = ?", kpiDefinitionID), dbUtil.Where("sd_instance_id = ?", sdInstanceID))
	if kpiFulFulfillmentCheckResultEntityLoadResult.IsFailure() {
		err := kpiFulFulfillmentCheckResultEntityLoadResult.GetError()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sharedUtils.NewSuccessResult[sharedUtils.Optional[dllModel.KPIFulfillmentCheckResult]](sharedUtils.NewEmptyOptional[dllModel.KPIFulfillmentCheckResult]())
		} else {
			return sharedUtils.NewFailureResult[sharedUtils.Optional[dllModel.KPIFulfillmentCheckResult]](err)
		}
	}
	return sharedUtils.NewSuccessResult[sharedUtils.Optional[dllModel.KPIFulfillmentCheckResult]](sharedUtils.NewOptionalOf(db2dll.ToDLLModelKPIFulfillmentCheckResult(kpiFulFulfillmentCheckResultEntityLoadResult.GetPayload())))
}

func (r *relationalDatabaseClientImpl) LoadKPIFulFulfillmentCheckResults() sharedUtils.Result[[]dllModel.KPIFulfillmentCheckResult] {
	r.mu.Lock()
	defer r.mu.Unlock()
	kpiFulFulfillmentCheckResultEntitiesLoadResult := dbUtil.LoadEntitiesFromDB[dbModel.KPIFulfillmentCheckResultEntity](r.db)
	if kpiFulFulfillmentCheckResultEntitiesLoadResult.IsFailure() {
		return sharedUtils.NewFailureResult[[]dllModel.KPIFulfillmentCheckResult](kpiFulFulfillmentCheckResultEntitiesLoadResult.GetError())
	}
	return sharedUtils.NewSuccessResult[[]dllModel.KPIFulfillmentCheckResult](sharedUtils.Map(kpiFulFulfillmentCheckResultEntitiesLoadResult.GetPayload(), db2dll.ToDLLModelKPIFulfillmentCheckResult))
}

func (r *relationalDatabaseClientImpl) LoadSDInstanceGroups() sharedUtils.Result[[]dllModel.SDInstanceGroup] {
	r.mu.Lock()
	defer r.mu.Unlock()
	sdInstanceGroupEntitiesLoadResult := dbUtil.LoadEntitiesFromDB[dbModel.SDInstanceGroupEntity](r.db, dbUtil.Preload("GroupMembershipRecords"))
	if sdInstanceGroupEntitiesLoadResult.IsFailure() {
		return sharedUtils.NewFailureResult[[]dllModel.SDInstanceGroup](sdInstanceGroupEntitiesLoadResult.GetError())
	}
	return sharedUtils.NewSuccessResult(sharedUtils.Map(sdInstanceGroupEntitiesLoadResult.GetPayload(), db2dll.ToDLLModelSDInstanceGroup))
}

func (r *relationalDatabaseClientImpl) LoadSDInstanceGroup(id uint32) sharedUtils.Result[dllModel.SDInstanceGroup] {
	r.mu.Lock()
	defer r.mu.Unlock()
	sdInstanceGroupEntityLoadResult := dbUtil.LoadEntityFromDB[dbModel.SDInstanceGroupEntity](r.db, dbUtil.Preload("GroupMembershipRecords"), dbUtil.Where("id = ?", id))
	if sdInstanceGroupEntityLoadResult.IsFailure() {
		return sharedUtils.NewFailureResult[dllModel.SDInstanceGroup](sdInstanceGroupEntityLoadResult.GetError())
	}
	return sharedUtils.NewSuccessResult(db2dll.ToDLLModelSDInstanceGroup(sdInstanceGroupEntityLoadResult.GetPayload()))
}

func (r *relationalDatabaseClientImpl) PersistSDInstanceGroup(sdInstanceGroup dllModel.SDInstanceGroup) sharedUtils.Result[uint32] {
	r.mu.Lock()
	defer r.mu.Unlock()
	sdInstanceGroupEntity := dll2db.ToDBModelEntitySDInstanceGroup(sdInstanceGroup)
	if err := dbUtil.PersistEntityIntoDB(r.db, &sdInstanceGroupEntity); err != nil {
		return sharedUtils.NewFailureResult[uint32](err)
	}
	return sharedUtils.NewSuccessResult[uint32](sdInstanceGroupEntity.ID)
}

func (r *relationalDatabaseClientImpl) DeleteSDInstanceGroup(id uint32) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return dbUtil.DeleteCertainEntityBasedOnId[dbModel.SDInstanceGroupEntity](r.db, id)
}
