package repository

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/guregu/null.v3"
)

type Namespace struct {
	CreatedAt null.Time `gorm:"column:created_at"`
	ID        int64     `gorm:"column:id;primary_key"`
	Name      string    `gorm:"column:name"`
	NameEn    string    `gorm:"column:name_en"`
	UpdatedAt null.Time `gorm:"column:updated_at"`
}

// TableName sets the insert table name for this struct type
func (n *Namespace) TableName() string {
	return "namespaces"
}

type NamespaceRepository interface {
	Find(name string) (res *Namespace, err error)
	Create(ns *Namespace) error
}

type namespace struct {
	db *gorm.DB
}

func NewNamespaceRepository(db *gorm.DB) NamespaceRepository {
	return &namespace{db: db}
}

func (c *namespace) Find(name string) (res *Namespace, err error) {
	var ns Namespace
	if err = c.db.First(&ns, "name_en = ?", name).Error; err != nil {
		return
	}
	return &ns, nil
}

func (c *namespace) Create(ns *Namespace) error {
	return c.db.Save(ns).Error
}
