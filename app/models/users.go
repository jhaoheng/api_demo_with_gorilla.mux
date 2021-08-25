package models

import (
	"strconv"
	"time"

	"gorm.io/gorm"
)

type USERS []USER

type USER struct {
	Acct      string    `gorm:"column:acct;primaryKey"`
	Pwd       string    `gorm:"column:pwd"`
	Fullname  string    `gorm:"column:fullname"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// TableName -
func (USER) TableName() string {
	return "users"
}

func (u *USER) FindUser(account, password string) error {
	return DB.Where("acct=? and pwd=?", account, password).Take(&u).Error
}

// ListAll -
func (users *USERS) ListBy(paging, sorting string, pageSize int, total *int64) (err error) {
	Paginate := func(paging string) func(db *gorm.DB) *gorm.DB {
		return func(db *gorm.DB) *gorm.DB {
			page, _ := strconv.Atoi(paging)
			if page == 0 {
				page = 1
			}

			offset := (page - 1) * pageSize
			return db.Offset(offset).Limit(pageSize)
		}
	}
	DB.Model(&USER{}).Count(total)
	if sorting == "asc" {
		return DB.Scopes(Paginate(paging)).Order("acct ASC").Find(&users).Error
	}
	return DB.Scopes(Paginate(paging)).Order("acct DESC").Find(&users).Error
}

// SearchByFullname -
func (u *USER) SearchByFullname(fullname string) error {
	return DB.Where("fullname=?", fullname).Take(&u).Error
}

// Create -
func (u *USER) Create() error {
	return DB.Create(&u).Error
}

// GetUserDetail -
func (u *USER) GetUserDetail(account string) error {
	return DB.Where("acct=?", account).Take(&u).Error
}

// Delete -
func (u *USER) Delete() (rowsAffected int64, err error) {
	result := DB.Delete(&u)
	return result.RowsAffected, result.Error
}

// Update -
func (u *USER) Update(pwd, fullname string) (rowsAffected int64, err error) {
	t := DB.Model(&u).Where("acct=?", u.Acct)
	updateUserData := USER{}
	if len(pwd) != 0 {
		updateUserData.Pwd = pwd
	}
	if len(fullname) != 0 {
		updateUserData.Fullname = fullname
	}
	result := t.Updates(updateUserData)
	return result.RowsAffected, result.Error
}
