package db

import (
	"context"
	"github.com/Abathargh/moody-go/gateway/models"

	"github.com/smallnest/gen/dbmeta"
)

// GetAllNeuralMetum is a function to get a slice of record(s) from neural_meta table in the main database
// params - page     - page requested (defaults to 0)
// params - pagesize - number of records in a page  (defaults to 20)
// params - order    - db sort order column
// error - NotFound, db Find error
func GetAllNeuralMeta(ctx context.Context, page, pagesize int64, order string) (neuralmeta []*models.NeuralMetum, totalRows int, err error) {

	neuralmeta = []*models.NeuralMetum{}

	neuralmeta_orm := DB.Model(&models.NeuralMetum{})
	neuralmeta_orm.Count(&totalRows)

	if page > 0 {
		offset := (page - 1) * pagesize
		neuralmeta_orm = neuralmeta_orm.Offset(offset).Limit(pagesize)
	} else {
		neuralmeta_orm = neuralmeta_orm.Limit(pagesize)
	}

	if order != "" {
		neuralmeta_orm = neuralmeta_orm.Order(order)
	}

	if err = neuralmeta_orm.Find(&neuralmeta).Error; err != nil {
		err = NotFound
		return nil, -1, err
	}

	return neuralmeta, totalRows, nil
}

// GetNeuralMetum is a function to get a single record to neural_meta table in the main database
// error - NotFound, db Find error
func GetNeuralMetum(ctx context.Context, id interface{}) (record *models.NeuralMetum, err error) {
	if err = DB.First(record, id).Error; err != nil {
		err = NotFound
		return nil, err
	}

	return record, nil
}

// AddNeuralMetum is a function to add a single record to neural_meta table in the main database
// error - InsertFailedError, db save call failed
func AddNeuralMetum(ctx context.Context, neuralmetum *models.NeuralMetum) (err error) {

	if err = DB.Save(neuralmetum).Error; err != nil {
		err = InsertFailedError
		return
	}

	return nil
}

// UpdateNeuralMetum is a function to update a single record from neural_meta table in the main database
// error - NotFound, db record for id not found
// error - UpdateFailedError, db meta data copy failed or db.Save call failed
func UpdateNeuralMetum(ctx context.Context, id interface{}, updated *models.NeuralMetum) (err error) {

	neuralmetum := &models.NeuralMetum{}
	if err = DB.First(neuralmetum, id).Error; err != nil {
		err = NotFound
		return
	}

	if err = dbmeta.Copy(neuralmetum, updated); err != nil {
		err = UpdateFailedError
		return
	}

	if err = DB.Save(neuralmetum).Error; err != nil {
		err = UpdateFailedError
		return
	}

	return nil
}

// DeleteNeuralMetum is a function to delete a single record from neural_meta table in the main database
// error - NotFound, db Find error
// error - DeleteFailedError, db Delete failed error
func DeleteNeuralMetum(ctx context.Context, id interface{}) (err error) {

	neuralmetum := &models.NeuralMetum{}

	if DB.First(neuralmetum, id).Error != nil {
		err = NotFound
		return
	}
	if err = DB.Delete(neuralmetum).Error; err != nil {
		err = DeleteFailedError
		return
	}

	return nil
}
