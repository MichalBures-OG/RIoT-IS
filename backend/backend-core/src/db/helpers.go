package db

import (
	cUtil "github.com/MichalBures-OG/bp-bures-SfPDfSD-commons/src/util"
	"gorm.io/gorm"
	"log"
	"reflect"
)

type WhereClausePair = cUtil.Pair[string, any]
type PreloadPathString = string

func WhereClause(s string, a any) WhereClausePair {
	return cUtil.NewPairOf(s, a)
}

func PreloadPath(s string) PreloadPathString {
	return s
}

func prepareLoadQuery(g *gorm.DB, args ...any) *gorm.DB {
	query := g
	cUtil.ForEach(args, func(arg any) {
		switch typedArg := arg.(type) {
		case WhereClausePair:
			query = query.Where(typedArg.GetFirst(), typedArg.GetSecond())
			break
		case PreloadPathString:
			query = query.Preload(typedArg)
			break
		default:
			log.Printf("prepareLoadQuery: unexpected arg type: %s\n", reflect.TypeOf(typedArg).String())
		}
	})
	return query
}

func LoadEntitiesFromDB[T any](g *gorm.DB, args ...any) cUtil.Result[[]T] {
	entities := make([]T, 0)
	if err := prepareLoadQuery(g, args...).Find(&entities).Error; err != nil {
		return cUtil.NewFailureResult[[]T](nil)
	}
	return cUtil.NewSuccessResult(entities)
}

func LoadEntityFromDB[T any](g *gorm.DB, args ...any) cUtil.Result[T] {
	var entity T
	if err := prepareLoadQuery(g, args...).First(&entity).Error; err != nil {
		return cUtil.NewFailureResult[T](nil)
	}
	return cUtil.NewSuccessResult(entity)
}

func PersistEntityIntoDB[T any](g *gorm.DB, entityReference *T) error {
	return g.Save(entityReference).Error
}

func DoesSuchEntityExist[T any](g *gorm.DB, whereClausePairs ...WhereClausePair) cUtil.Result[bool] {
	entityCount := int64(0)
	query := g.Model(new(T))
	cUtil.ForEach(whereClausePairs, func(whereClausePair WhereClausePair) {
		query = query.Where(whereClausePair.GetFirst(), whereClausePair.GetSecond())
	})
	if err := query.Count(&entityCount).Error; err != nil {
		return cUtil.NewFailureResult[bool](err)
	}
	return cUtil.NewSuccessResult[bool](entityCount > 0)
}

func DeleteCertainEntity[T any](g *gorm.DB, id uint32) error {
	return g.Delete(new(T), id).Error
}

func DeleteEntities[T any](g *gorm.DB, whereClausePairs ...WhereClausePair) error {
	query := g
	cUtil.ForEach(whereClausePairs, func(whereClausePair WhereClausePair) {
		query = query.Where(whereClausePair.GetFirst(), whereClausePair.GetSecond())
	})
	return query.Delete(new(T)).Error
}
