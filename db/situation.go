package db

import (
	"github.com/Abathargh/moody-go/models"
)

// GetAllSituation is a function to get a slice of record(s) from situation table in the main database
// params - page     - page requested (defaults to 0)
// params - pagesize - number of records in a page  (defaults to 20)
// error - NotFound, db Find error
func GetAllSituations() (situations []*models.Situation, totalRows int64, err error) {
	situations = []*models.Situation{}
	situationOrm := DB.Model(&models.Situation{})
	situationOrm.Count(&totalRows)

	if err = situationOrm.Find(&situations).Error; err != nil {
		err = NotFound
		return nil, -1, err
	}

	return situations, totalRows, nil
}

// GetSituation is a function to get a single record to situation table in the main database
// error - NotFound, db Find error
func GetSituation(name string) (*models.Situation, error) {
	var situation models.Situation
	if err := DB.Where("name = ?", name).First(&situation).Error; err != nil {
		return nil, err
	}

	return &situation, nil
}

// AddSituation is a function to add a single record to situation table in the main database
// error - InsertFailedError, db save call failed
func AddSituation(situation *models.Situation) (err error) {
	if err = DB.Save(situation).Error; err != nil {
		err = InsertFailedError
		return
	}

	return nil
}

// DeleteSituation is a function to delete a single record from situation table in the main database
// error - NotFound, db Find error
// error - DeleteFailedError, db Delete failed error
func DeleteSituation(situation *models.Situation) (err error) {
	if err = DB.Delete(situation).Error; err != nil {
		err = DeleteFailedError
		return
	}
	return nil
}
