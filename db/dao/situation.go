package dao

import (
	"context"

	"github.com/Abathargh/moody-go/db/model"

	"github.com/smallnest/gen/dbmeta"
)

// GetAllSituation is a function to get a slice of record(s) from situation table in the main database
// params - page     - page requested (defaults to 0)
// params - pagesize - number of records in a page  (defaults to 20)
// params - order    - db sort order column
// error - NotFound, db Find error
func GetAllSituations(_ context.Context, page, pagesize int64, order string) (situations []*model.Situation, totalRows int, err error) {

	situations = []*model.Situation{}

	situations_orm := DB.Model(&model.Situation{})
	situations_orm.Count(&totalRows)

	if page > 0 {
		offset := (page - 1) * pagesize
		situations_orm = situations_orm.Offset(offset).Limit(pagesize)
	} else {
		situations_orm = situations_orm.Limit(pagesize)
	}

	if order != "" {
		situations_orm = situations_orm.Order(order)
	}

	if err = situations_orm.Find(&situations).Error; err != nil {
		err = NotFound
		return nil, -1, err
	}

	return situations, totalRows, nil
}

// GetSituation is a function to get a single record to situation table in the main database
// error - NotFound, db Find error
func GetSituation(ctx context.Context, id interface{}) (record *model.Situation, err error) {
	if err = DB.First(record, id).Error; err != nil {
		err = NotFound
		return nil, err
	}

	return record, nil
}

// AddSituation is a function to add a single record to situation table in the main database
// error - InsertFailedError, db save call failed
func AddSituation(ctx context.Context, situation *model.Situation) (err error) {

	if err = DB.Save(situation).Error; err != nil {
		err = InsertFailedError
		return
	}

	return nil
}

// UpdateSituation is a function to update a single record from situation table in the main database
// error - NotFound, db record for id not found
// error - UpdateFailedError, db meta data copy failed or db.Save call failed
func UpdateSituation(ctx context.Context, id interface{}, updated *model.Situation) (err error) {

	situation := &model.Situation{}
	if err = DB.First(situation, id).Error; err != nil {
		err = NotFound
		return
	}

	if err = dbmeta.Copy(situation, updated); err != nil {
		err = UpdateFailedError
		return
	}

	if err = DB.Save(situation).Error; err != nil {
		err = UpdateFailedError
		return
	}

	return nil
}

// DeleteSituation is a function to delete a single record from situation table in the main database
// error - NotFound, db Find error
// error - DeleteFailedError, db Delete failed error
func DeleteSituation(ctx context.Context, id interface{}) (err error) {

	situation := &model.Situation{}

	if DB.First(situation, id).Error != nil {
		err = NotFound
		return
	}
	if err = DB.Delete(situation).Error; err != nil {
		err = DeleteFailedError
		return
	}

	return nil
}
