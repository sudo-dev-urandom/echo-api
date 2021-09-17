package models

import (
	"goecho/core"
)

type (
	User struct {
		core.Model
		Name    string `json:"name" gorm:"column:name"`
		Address string `json:"address" gorm:"column:address"`
	}
)

func (User) TableName() string {
	return "users"
}

func (p *User) Create() error {
	err := core.Create(&p)
	return err
}

func (p *User) Save() error {
	err := core.Save(&p)
	return err
}

func (p *User) Delete() error {
	err := core.Delete(&p)
	return err
}

func (p *User) FindbyID(id int) error {
	err := core.FindbyID(&p, id)
	return err
}

func (p *User) Find(filter interface{}) error {
	err := core.Find(&p, filter)
	return err
}

func (b *User) PagedFilterSearch(page int, rows int, orderby string, sort string, filter interface{}) (result core.PagedFindResult, err error) {
	Question := []User{}
	orders := []string{orderby}
	sorts := []string{sort}
	result, err = core.PagedFindFilter(&Question, page, rows, orders, sorts, filter, []string{})

	return result, err
}
