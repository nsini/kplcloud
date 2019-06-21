package repository

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/guregu/null.v3"
)

type Member struct {
	CreatedAt null.Time `gorm:"column:created_at"`
	Email     string    `gorm:"column:email"`
	ID        int64     `gorm:"column:id;primary_key"`
	Openid    string    `gorm:"column:openid"`
	Phone     string    `gorm:"column:phone"`
	State     int64     `gorm:"column:state"`
	UpdatedAt null.Time `gorm:"column:updated_at"`
	Username  string    `gorm:"column:username"`
}

// TableName sets the insert table name for this struct type
func (m *Member) TableName() string {
	return "members"
}

type MemberRepository interface {
	Find(name string) (res *Member, err error)
}

type member struct {
	db *gorm.DB
}

func NewMemberRepository(db *gorm.DB) MemberRepository {
	return &member{db: db}
}

func (c *member) Find(name string) (res *Member, err error) {

	return
}
