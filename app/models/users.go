package models

import (
	"strconv"
	"time"

	"gorm.io/gorm"
)

type IUser interface {
	SetAcct(acct string) IUser
	SetPwd(password string) IUser
	SetFullname(fullname string) IUser
	//
	Create() error
	Get() (User, error)
	GetAll() ([]User, error)
	GetAllCount() (int64, error)
	Delete() (rowsAffected int64, err error)
	Update(user User) (rowsAffected int64, err error)
	//
	ListBy(paging, sorting string, page_size int) ([]User, error)
	//
	Or(users ...User) IUser
}

type User struct {
	tx *gorm.DB
	//
	Acct      string    `gorm:"column:acct;primaryKey"`
	Pwd       string    `gorm:"column:pwd"`
	Fullname  string    `gorm:"column:fullname"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// TableName -
func (User) TableName() string {
	return "user"
}

func NewUser() IUser {
	return &User{
		tx: DB,
	}
}

func (model *User) set_db() *gorm.DB {
	if model.tx == nil {
		model.tx = DB
	}
	return model.tx
}

func (model *User) SetAcct(acct string) IUser {
	model.Acct = acct
	return model
}

func (model *User) SetPwd(password string) IUser {
	model.Pwd = password
	return model
}

func (model *User) SetFullname(fullname string) IUser {
	model.Fullname = fullname
	return model
}

// Create -
func (model *User) Create() error {
	return model.set_db().Create(&model).Error
}

func (model *User) Or(users ...User) IUser {
	tx := DB
	for _, user := range users {
		tx = tx.Or(user)
	}
	model.tx = model.set_db().Where(tx)
	return model
}

func (model *User) Get() (User, error) {
	output := User{}
	tx := model.set_db().Where(model).Take(&output)
	return output, tx.Error
}

func (model *User) GetAll() ([]User, error) {
	output := []User{}
	tx := model.set_db().Where(model).Find(&output)
	return output, tx.Error
}

func (model *User) GetAllCount() (int64, error) {
	var total int64 = 0
	tx := model.set_db().Model(&User{}).Count(&total)
	return total, tx.Error
}

// Delete -
func (model *User) Delete() (rowsAffected int64, err error) {
	tx := model.set_db().Delete(&model)
	return tx.RowsAffected, tx.Error
}

// Update -
func (model *User) Update(user User) (rowsAffected int64, err error) {
	tx := model.set_db().Model(&model).Where(model).Updates(user)
	return tx.RowsAffected, tx.Error
}

// ListAll -
func (model *User) ListBy(paging, sorting string, page_size int) ([]User, error) {
	output := []User{}
	Paginate := func(paging string) func(db *gorm.DB) *gorm.DB {
		return func(db *gorm.DB) *gorm.DB {
			page, _ := strconv.Atoi(paging)
			if page == 0 {
				page = 1
			}

			offset := (page - 1) * page_size
			return db.Offset(offset).Limit(page_size)
		}
	}
	if sorting == "asc" {
		tx := model.set_db().Scopes(Paginate(paging)).Order("acct ASC").Find(&output)
		return output, tx.Error
	}
	tx := model.set_db().Scopes(Paginate(paging)).Order("acct DESC").Find(&output)
	return output, tx.Error
}
