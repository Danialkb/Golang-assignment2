package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	info := "host=localhost user=postgres password=Dankb2131193* dbname=assignment2 port=5432"
	db, err := gorm.Open(postgres.Open(info), &gorm.Config{})
	return db, err
}

func ComparePasswords(hashedPwd string, plainPwd []byte) error {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		return err
	}
	return nil
}

type Registration interface {
	Register(*User) error
}

type RegistrationService struct{}

func (r *RegistrationService) Register(u *User) error {
	db, _ := InitDB()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.Password = string(hashedPassword)
	err := db.Create(u).Error
	db.AutoMigrate(&User{})
	return err
}

type Authorization interface {
	SignIn(username, password string) error
}

type AuthorizationService struct{}

func (a *AuthorizationService) SignIn(username, password string) (*User, error) {
	db, _ := InitDB()
	var user *User
	res := db.Where("username = ?", username).First(&user)
	if res.Error != nil {
		return nil, res.Error
	}

	err := ComparePasswords(user.Password, []byte(password))
	return user, err
}

type Searching interface {
	Search(string)
}

type SearchingService struct{}

func (s *SearchingService) Search(name string) []Item {
	db, _ := InitDB()
	var items []Item
	if name != "" {
		db.Where("name LIKE ?", name+"%").
			Joins("left join item_rating on items.id = item_rating.item_id").
			Select("items.*, avg(item_rating.rate) as rating").
			Group("items.id").
			Find(&items)
	} else {
		db.Joins("left join item_rating on items.id = item_rating.item_id").
			Select("items.*, avg(item_rating.rate) as rating").
			Group("items.id").
			Find(&items)
	}
	return items
}

type Filtering interface {
	FilterByName()
	FilterByPrice()
}

type FilteringService struct{}

func (f *FilteringService) FilterByPrice() []Item {
	var items []Item
	db, _ := InitDB()
	db.Joins("left join item_rating on items.id = item_rating.item_id").
		Order("price").
		Select("items.*, avg(item_rating.rate) as rating").
		Group("items.id").
		Find(&items)
	return items
}

func (f *FilteringService) FilterByName() []Item {
	var items []Item
	db, _ := InitDB()
	db.Joins("left join item_rating on items.id = item_rating.item_id").
		Order("name").
		Select("items.*, avg(item_rating.rate) as rating").
		Group("items.id").
		Find(&items)
	return items
}

type GiveRating interface {
	RateItem(useId, itemId uint, rate float64)
}

type RatingService struct{}

func (r *RatingService) RateItem(useId, itemId uint, rate float64) {
	db, _ := InitDB()
	res := db.Exec("SELECT * FROM item_rating WHERE user_id=? and item_id=?", useId, itemId)
	if res.RowsAffected == 0 {
		db.Exec("INSERT INTO item_rating (user_id, item_id, rate) VALUES (?, ?, ?)", useId, itemId, rate)
	} else {
		db.Exec("UPDATE item_rating SET rate=? WHERE user_id=? and item_id=?", rate, useId, itemId)
	}
}
