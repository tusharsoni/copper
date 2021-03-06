package phone

import "time"

type Credentials struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`

	UserUUID         string `gorm:"not null;uniqueIndex"`
	PhoneNumber      string `gorm:"not null;uniqueIndex"`
	Verified         bool   `gorm:"not null;default:false"`
	VerificationCode uint   `gorm:"not null"`
}

func (Credentials) TableName() string {
	return "cauth_phone_credentials"
}
