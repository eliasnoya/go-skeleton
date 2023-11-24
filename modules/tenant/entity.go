package tenant

import (
	"time"

	"gorm.io/gorm"
)

type Tenant struct {
	ID                uint   `gorm:"primaryKey"`
	Domain            string `gorm:"type:VARCHAR(100)"`
	BillingCompany    string `gorm:"type:VARCHAR(150)"`
	BillingTaxNumber  string `gorm:"type:VARCHAR(40)"`
	BillingTaxSubject string `gorm:"type:VARCHAR(50)"`
	BillingAddress    string `gorm:"type:VARCHAR(200)"`
	PersonID          uint   `gorm:"type:BIGINT"`
	Plan              string `gorm:"type:ENUM('free','core','premium');default:'core'"`
	DefaultTheme      string `gorm:"type:VARCHAR(6);default:'light'"`
	LightSetup        string `gorm:"type:JSON"`
	DarkSetup         string `gorm:"type:JSON"`
	Logo              string `gorm:"type:MEDIUMTEXT"`
	Approved          bool   `gorm:"type:TINYINT;default:0"`
	Configured        bool   `gorm:"type:TINYINT;default:0"`
	SetupStep         int    `gorm:"type:TINYINT;default:0"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt
}

func (t *Tenant) TableName() string {
	return "tenants"
}
