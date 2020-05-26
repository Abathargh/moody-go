package db

import (
	"context"
	"github.com/Abathargh/moody-go/models"

	"github.com/smallnest/gen/dbmeta"
)

// GetAllSituationSetting is a function to get a slice of record(s) from situation_setting table in the main database
// params - page     - page requested (defaults to 0)
// params - pagesize - number of records in a page  (defaults to 20)
// params - order    - db sort order column
// error - NotFound, db Find error
func GetAllSituationSettings(ctx context.Context, page, pagesize int64, order string) (situationsettings []*models.SituationSetting, totalRows int, err error) {

	situationsettings = []*models.SituationSetting{}

	situationsettings_orm := DB.Model(&models.SituationSetting{})
	situationsettings_orm.Count(&totalRows)

	if page > 0 {
		offset := (page - 1) * pagesize
		situationsettings_orm = situationsettings_orm.Offset(offset).Limit(pagesize)
	} else {
		situationsettings_orm = situationsettings_orm.Limit(pagesize)
	}

	if order != "" {
		situationsettings_orm = situationsettings_orm.Order(order)
	}

	if err = situationsettings_orm.Find(&situationsettings).Error; err != nil {
		err = NotFound
		return nil, -1, err
	}

	return situationsettings, totalRows, nil
}

// GetSituationSetting is a function to get a single record to situation_setting table in the main database
// error - NotFound, db Find error
func GetSituationSetting(ctx context.Context, id interface{}) (record *models.SituationSetting, err error) {
	if err = DB.First(record, id).Error; err != nil {
		err = NotFound
		return nil, err
	}

	return record, nil
}

// AddSituationSetting is a function to add a single record to situation_setting table in the main database
// error - InsertFailedError, db save call failed
func AddSituationSetting(ctx context.Context, situationsetting *models.SituationSetting) (err error) {

	if err = DB.Save(situationsetting).Error; err != nil {
		err = InsertFailedError
		return
	}

	return nil
}

// UpdateSituationSetting is a function to update a single record from situation_setting table in the main database
// error - NotFound, db record for id not found
// error - UpdateFailedError, db meta data copy failed or db.Save call failed
func UpdateSituationSetting(ctx context.Context, id interface{}, updated *models.SituationSetting) (err error) {

	situationsetting := &models.SituationSetting{}
	if err = DB.First(situationsetting, id).Error; err != nil {
		err = NotFound
		return
	}

	if err = dbmeta.Copy(situationsetting, updated); err != nil {
		err = UpdateFailedError
		return
	}

	if err = DB.Save(situationsetting).Error; err != nil {
		err = UpdateFailedError
		return
	}

	return nil
}

// DeleteSituationSetting is a function to delete a single record from situation_setting table in the main database
// error - NotFound, db Find error
// error - DeleteFailedError, db Delete failed error
func DeleteSituationSetting(ctx context.Context, id interface{}) (err error) {

	situationsetting := &models.SituationSetting{}

	if DB.First(situationsetting, id).Error != nil {
		err = NotFound
		return
	}
	if err = DB.Delete(situationsetting).Error; err != nil {
		err = DeleteFailedError
		return
	}

	return nil
}
