package models

// Situation struct is a row record of the situation table in the main database
type Situation struct {
	SituationId   uint64 `gorm:"column:id;primary_key" json:"id"`
	SituationName string `gorm:"column:name;unique" json:"name" validate:"nonzero"`
}

// TableName sets the insert table name for this struct type
func (s *Situation) TableName() string {
	return "situation"
}

func GetAllSituations() (situations []*Situation, totalRows int64, err error) {
	situations = []*Situation{}
	situationOrm := DB.Model(&Situation{})
	situationOrm.Count(&totalRows)

	if err = situationOrm.Find(&situations).Error; err != nil {
		err = NotFound
		return nil, -1, err
	}

	return situations, totalRows, nil
}

// GetSituation is a function to get a single record to situation table in the main database
// error - NotFound, models Find error
func GetSituation(id uint64) (*Situation, error) {
	var situation Situation
	if err := DB.First(&situation, id).Error; err != nil {
		return nil, err
	}

	return &situation, nil
}

// AddSituation is a function to add a single record to situation table in the main database
// error - InsertFailedError, models save call failed
func AddSituation(situation *Situation) (err error) {
	if err = DB.Save(situation).Error; err != nil {
		err = InsertFailedError
		return
	}

	return nil
}

// DeleteSituation is a function to delete a single record from situation table in the main database
// error - NotFound, models Find error
// error - DeleteFailedError, models Delete failed error
func DeleteSituation(s *Situation) (err error) {
	situation := &Situation{}
	if err := DB.First(&situation, s.SituationId).Error; err != nil {
		return NotFound
	}
	if err := DB.Delete(situation).Error; err != nil {
		return DeleteFailedError
	}

	return nil
}
