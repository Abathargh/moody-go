package db

import (
	"context"
	"github.com/Abathargh/moody-go/models"

	"github.com/smallnest/gen/dbmeta"
)

// GetAllNode is a function to get a slice of record(s) from node table in the main database
// params - page     - page requested (defaults to 0)
// params - pagesize - number of records in a page  (defaults to 20)
// params - order    - db sort order column
// error - NotFound, db Find error
func GetAllNodes(ctx context.Context, page, pagesize int64, order string) (nodes []*models.Node, totalRows int, err error) {

	nodes = []*models.Node{}

	nodes_orm := DB.Model(&models.Node{})
	nodes_orm.Count(&totalRows)

	if page > 0 {
		offset := (page - 1) * pagesize
		nodes_orm = nodes_orm.Offset(offset).Limit(pagesize)
	} else {
		nodes_orm = nodes_orm.Limit(pagesize)
	}

	if order != "" {
		nodes_orm = nodes_orm.Order(order)
	}

	if err = nodes_orm.Find(&nodes).Error; err != nil {
		err = NotFound
		return nil, -1, err
	}

	return nodes, totalRows, nil
}

// GetNode is a function to get a single record to node table in the main database
// error - NotFound, db Find error
func GetNode(ctx context.Context, mac string) (*models.Node, error) {
	var record models.Node
	if err := DB.Where("node_macaddress = ?", mac).First(&record).Error; err != nil {
		return nil, err
	}

	return &record, nil
}

// AddNode is a function to add a single record to node table in the main database
// error - InsertFailedError, db save call failed
func AddNode(ctx context.Context, node *models.Node) (err error) {

	if err = DB.Save(node).Error; err != nil {
		err = InsertFailedError
		return
	}

	return nil
}

// UpdateNode is a function to update a single record from node table in the main database
// error - NotFound, db record for id not found
// error - UpdateFailedError, db meta data copy failed or db.Save call failed
func UpdateNode(ctx context.Context, id interface{}, updated *models.Node) (err error) {

	node := &models.Node{}
	if err = DB.First(node, id).Error; err != nil {
		err = NotFound
		return
	}

	if err = dbmeta.Copy(node, updated); err != nil {
		err = UpdateFailedError
		return
	}

	if err = DB.Save(node).Error; err != nil {
		err = UpdateFailedError
		return
	}

	return nil
}

// DeleteNode is a function to delete a single record from node table in the main database
// error - NotFound, db Find error
// error - DeleteFailedError, db Delete failed error
func DeleteNode(ctx context.Context, id interface{}) (err error) {

	node := &models.Node{}

	if DB.First(node, id).Error != nil {
		err = NotFound
		return
	}
	if err = DB.Delete(node).Error; err != nil {
		err = DeleteFailedError
		return
	}

	return nil
}
