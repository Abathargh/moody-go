package db

import (
	"context"
	"github.com/Abathargh/moody-go/gateway/models"

	"github.com/smallnest/gen/dbmeta"
)

// GetAllNeuralDatum is a function to get a slice of record(s) from neural_data table in the main database
// params - page     - page requested (defaults to 0)
// params - pagesize - number of records in a page  (defaults to 20)
// params - order    - db sort order column
// error - NotFound, db Find error
func GetAllNeuralData(ctx context.Context, page, pagesize int64, order string) (neuraldata []*models.NeuralDatum, totalRows int, err error) {

	neuraldata = []*models.NeuralDatum{}

	neuraldata_orm := DB.Model(&models.NeuralDatum{})
	neuraldata_orm.Count(&totalRows)

	if page > 0 {
		offset := (page - 1) * pagesize
		neuraldata_orm = neuraldata_orm.Offset(offset).Limit(pagesize)
	} else {
		neuraldata_orm = neuraldata_orm.Limit(pagesize)
	}

	if order != "" {
		neuraldata_orm = neuraldata_orm.Order(order)
	}

	if err = neuraldata_orm.Find(&neuraldata).Error; err != nil {
		err = NotFound
		return nil, -1, err
	}

	return neuraldata, totalRows, nil
}

// GetNeuralDatum is a function to get a single record to neural_data table in the main database
// error - NotFound, db Find error
func GetNeuralDatum(ctx context.Context, id interface{}) (record *models.NeuralDatum, err error) {
	if err = DB.First(record, id).Error; err != nil {
		err = NotFound
		return nil, err
	}

	return record, nil
}

// AddNeuralDatum is a function to add a single record to neural_data table in the main database
// error - InsertFailedError, db save call failed
func AddNeuralDatum(ctx context.Context, neuraldatum *models.NeuralDatum) (err error) {

	if err = DB.Save(neuraldatum).Error; err != nil {
		err = InsertFailedError
		return
	}

	return nil
}

// UpdateNeuralDatum is a function to update a single record from neural_data table in the main database
// error - NotFound, db record for id not found
// error - UpdateFailedError, db meta data copy failed or db.Save call failed
func UpdateNeuralDatum(ctx context.Context, id interface{}, updated *models.NeuralDatum) (err error) {

	neuraldatum := &models.NeuralDatum{}
	if err = DB.First(neuraldatum, id).Error; err != nil {
		err = NotFound
		return
	}

	if err = dbmeta.Copy(neuraldatum, updated); err != nil {
		err = UpdateFailedError
		return
	}

	if err = DB.Save(neuraldatum).Error; err != nil {
		err = UpdateFailedError
		return
	}

	return nil
}

// DeleteNeuralDatum is a function to delete a single record from neural_data table in the main database
// error - NotFound, db Find error
// error - DeleteFailedError, db Delete failed error
func DeleteNeuralDatum(ctx context.Context, id interface{}) (err error) {

	neuraldatum := &models.NeuralDatum{}

	if DB.First(neuraldatum, id).Error != nil {
		err = NotFound
		return
	}
	if err = DB.Delete(neuraldatum).Error; err != nil {
		err = DeleteFailedError
		return
	}

	return nil
}
