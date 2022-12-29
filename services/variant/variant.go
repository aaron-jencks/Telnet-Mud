package variant

import (
	"mud/entities"
	"mud/utils/crud"
	"mud/utils/io/db"
)

func variantToArr(rs interface{}) []interface{} {
	rec := rs.(entities.TileVariant)
	return []interface{}{
		rec.Id,
		rec.Name,
		rec.Icon,
	}
}

func variantFromArr(arr []interface{}) interface{} {
	return entities.TileVariant{
		Id:   arr[1].(int),
		Name: arr[2].(string),
		Icon: arr[3].(string),
	}
}

func createVariantFunc(table *db.TableDefinition, args ...interface{}) []interface{} {
	if len(args) < 3 {
		id := 0
		if table.CSV.LineCount > 0 {
			id = table.RetrieveLine(table.CSV.LineCount - 1)[1].(int) + 1
		}

		return []interface{}{id, args[0], args[1]}
	} else {
		return []interface{}{args[0], args[1], args[2]}
	}
}

var CRUD crud.Crud = crud.CreateCrud("variants", variantToArr, variantFromArr, createVariantFunc)

func GetAllVariants(vid int) []entities.TileVariant {
	variants := CRUD.RetrieveAll(vid)

	var result []entities.TileVariant = make([]entities.TileVariant, len(variants))
	for vi, variant := range variants {
		result[vi] = variant.(entities.TileVariant)
	}

	return result
}

func GetSpecificVariant(vid int, name string) entities.TileVariant {
	variants := GetAllVariants(vid)
	for _, variant := range variants {
		if variant.Name == name {
			return variant
		}
	}

	return entities.TileVariant{
		Id:   -1,
		Name: "Not Found",
	}
}
