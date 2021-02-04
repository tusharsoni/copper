package cacl

import (
	"gorm.io/gorm"
)

func New(db *gorm.DB) Svc {
	return newSvcImpl(newSQLRepo(db))
}
