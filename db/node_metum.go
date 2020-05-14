package db

import (
	"context"
	"github.com/Abathargh/moody-go/models"

	"github.com/smallnest/gen/dbmeta"
)

// GetAllNodeMetum is a function to get a slice of record(s) from node_meta table in the main database
// params - page     - page requested (defaults to 0)
// params - pagesize - number of records in a page  (defaults to 20)
// params - order    - db sort order column
// error - NotFound, db Find error
func GetAllNodeMeta(ctx context.Context, page, pagesize int64, order string) (nodemeta []*models.NodeMeta, totalRows int, err error) {

	nodemeta = []*models.NodeMeta{}

	nodemeta_orm := DB.Model(&models.NodeMeta{})
	nodemeta_orm.Count(&totalRows)

	if page > 0 {
		offset := (page - 1) * pagesize
		nodemeta_orm = nodemeta_orm.Offset(offset).Limit(pagesize)
	} else {
		nodemeta_orm = nodemeta_orm.Limit(pagesize)
	}

	if order != "" {
		nodemeta_orm = nodemeta_orm.Order(order)
	}

	if err = nodemeta_orm.Find(&nodemeta).Error; err != nil {
		err = NotFound
		return nil, -1, err
	}

	return nodemeta, totalRows, nil
}

// GetNodeMetum is a function to get a single record to node_meta table in the main database
// error - NotFound, db Find error
func GetNodeMetum(ctx context.Context, id interface{}) (record *models.NodeMeta, err error) {
	if err = DB.First(record, id).Error; err != nil {
		err = NotFound
		return nil, err
	}

	return record, nil
}

// AddNodeMetum is a function to add a single record to node_meta table in the main database
// error - InsertFailedError, db save call failed
func AddNodeMetum(ctx context.Context, nodemetum *models.NodeMeta) (err error) {

	if err = DB.Save(nodemetum).Error; err != nil {
		err = InsertFailedError
		return
	}

	return nil
}

// UpdateNodeMetum is a function to update a single record from node_meta table in the main database
// error - NotFound, db record for id not found
// error - UpdateFailedError, db meta data copy failed or db.Save call failed
func UpdateNodeMetum(ctx context.Context, id interface{}, updated *models.NodeMeta) (err error) {

	nodemetum := &models.NodeMeta{}
	if err = DB.First(nodemetum, id).Error; err != nil {
		err = NotFound
		return
	}

	if err = dbmeta.Copy(nodemetum, updated); err != nil {
		err = UpdateFailedError
		return
	}

	if err = DB.Save(nodemetum).Error; err != nil {
		err = UpdateFailedError
		return
	}

	return nil
}

// DeleteNodeMetum is a function to delete a single record from node_meta table in the main database
// error - NotFound, db Find error
// error - DeleteFailedError, db Delete failed error
func DeleteNodeMetum(ctx context.Context, id interface{}) (err error) {

	nodemetum := &models.NodeMeta{}

	if DB.First(nodemetum, id).Error != nil {
		err = NotFound
		return
	}
	if err = DB.Delete(nodemetum).Error; err != nil {
		err = DeleteFailedError
		return
	}

	return nil
}
