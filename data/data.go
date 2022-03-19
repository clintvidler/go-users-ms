package data

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func NewStore(host, username, password, name string) (s *Store, err error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=5432 TimeZone=Australia/Sydney",
		host,
		username,
		password,
		name)

	gdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	gdb.AutoMigrate(&User{}, &Token{})

	s = &Store{db: gdb}

	return
}

// drop data, populate with seed data
func (s *Store) Populate() {
	s.db.Migrator().DropTable(&User{})
	s.db.Migrator().DropTable(&Token{})
	s.db.AutoMigrate(&User{}, &Token{})

	var us []User

	us = append(us, User{FirstName: "y", LastName: "z", Email: "x@x.x"})

	for _, u := range us {
		u.SetPassword("x")
		s.db.Create(&u)
	}
}
