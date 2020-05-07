package dao

import (
	"context"

	"github.com/Abathargh/moody-go/db/model"

	"github.com/smallnest/gen/dbmeta"
)

// GetAllAction is a function to get a slice of record(s) from action table in the main database
// params - page     - page requested (defaults to 0)
// params - pagesize - number of records in a page  (defaults to 20)
// params - order    - db sort order column
// error - NotFound, db Find error
func GetAllActions(ctx context.Context, page, pagesize int64, order string) (actions []*model.Action, totalRows int, err error) {

	actions = []*model.Action{}

	actions_orm := DB.Model(&model.Action{})
	actions_orm.Count(&totalRows)

	if page > 0 {
		offset := (page - 1) * pagesize
		actions_orm = actions_orm.Offset(offset).Limit(pagesize)
	} else {
		actions_orm = actions_orm.Limit(pagesize)
	}

	if order != "" {
		actions_orm = actions_orm.Order(order)
	}

	if err = actions_orm.Find(&actions).Error; err != nil {
		err = NotFound
		return nil, -1, err
	}

	return actions, totalRows, nil
}

// GetAction is a function to get a single record to action table in the main database
// error - NotFound, db Find error
func GetAction(ctx context.Context, id interface{}) (record *model.Action, err error) {
	if err = DB.First(record, id).Error; err != nil {
		err = NotFound
		return nil, err
	}

	return record, nil
}

// AddAction is a function to add a single record to action table in the main database
// error - InsertFailedError, db save call failed
func AddAction(ctx context.Context, action *model.Action) (err error) {

	if err = DB.Save(action).Error; err != nil {
		err = InsertFailedError
		return
	}

	return nil
}

// UpdateAction is a function to update a single record from action table in the main database
// error - NotFound, db record for id not found
// error - UpdateFailedError, db meta data copy failed or db.Save call failed
func UpdateAction(ctx context.Context, id interface{}, updated *model.Action) (err error) {

	action := &model.Action{}
	if err = DB.First(action, id).Error; err != nil {
		err = NotFound
		return
	}

	if err = dbmeta.Copy(action, updated); err != nil {
		err = UpdateFailedError
		return
	}

	if err = DB.Save(action).Error; err != nil {
		err = UpdateFailedError
		return
	}

	return nil
}

// DeleteAction is a function to delete a single record from action table in the main database
// error - NotFound, db Find error
// error - DeleteFailedError, db Delete failed error
func DeleteAction(ctx context.Context, id interface{}) (err error) {

	action := &model.Action{}

	if DB.First(action, id).Error != nil {
		err = NotFound
		return
	}
	if err = DB.Delete(action).Error; err != nil {
		err = DeleteFailedError
		return
	}

	return nil
}
