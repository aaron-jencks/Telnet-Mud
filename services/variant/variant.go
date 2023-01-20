package variant

import (
	"database/sql"
	"fmt"
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

// TODO something will have to be done here to handle using the id
func createVariantFunc(table db.TableDefinition, args ...interface{}) []interface{} {
	return []interface{}{args[0], args[1], args[2]}
}

func variantSelector(args []interface{}) string {
	return fmt.Sprintf("Id=%d and Name=\"%s\"", args[0].(int), args[1].(string))
}

func variantScanner(row *sql.Rows) (interface{}, error) {
	result := entities.TileVariant{}
	err := row.Scan(&result.Id, &result.Name, &result.Icon)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func variantUpdateFunc(oldValue, newValue interface{}) []crud.RowModStruct {
	ois := oldValue.(entities.TileVariant)
	nis := newValue.(entities.TileVariant)

	var result []crud.RowModStruct

	if ois.Id != nis.Id {
		result = append(result, crud.RowModStruct{
			Column:   "Id",
			NewValue: nis.Id,
		})
	}
	if ois.Name != nis.Name {
		result = append(result, crud.RowModStruct{
			Column:   "Name",
			NewValue: fmt.Sprintf("\"%s\"", nis.Name),
		})
	}
	if ois.Icon != nis.Icon {
		result = append(result, crud.RowModStruct{
			Column:   "Icon",
			NewValue: fmt.Sprintf("\"%s\"", nis.Icon),
		})
	}

	return result
}

var CRUD crud.Crud = crud.CreateCrud("variants", variantSelector, variantToArr, variantScanner, variantFromArr, createVariantFunc, variantUpdateFunc)

// TODO redo these to make use of queries

func GetAllVariants(vid int) []entities.TileVariant {
	table := CRUD.FetchTable()
	variants := table.QueryData(fmt.Sprintf("Id=%d", vid), variantScanner)

	var result []entities.TileVariant = make([]entities.TileVariant, len(variants))
	for vi, variant := range variants {
		result[vi] = variant.(entities.TileVariant)
	}

	return result
}
