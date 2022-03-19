package data

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestProdDB(t *testing.T) {
	_, err := NewStore(os.Getenv("DB_PROD"), os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	assert.NoError(t, err, "Failed to connect 'prod' database.")
}

type DataTestSuite struct {
	suite.Suite
	store *Store
}

// listen for 'go test' command --> run test methods
func TestDataTestSuite(t *testing.T) {
	suite.Run(t, new(DataTestSuite))
}

// run once, before test suite methods
func (s *DataTestSuite) SetupSuite() {
	store, err := NewStore(os.Getenv("DB_TEST"), os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	assert.NoError(s.T(), err, "Failed to connect 'test' database")

	store.PopulateTest()

	s.store = store
}

// run once, after test suite methods
func (s *DataTestSuite) TearDownSuite() {
	err := s.store.db.Migrator().DropTable(&User{})

	assert.NoError(s.T(), err, "Failed to drop 'User' table from the database.")
}

// run before each test
func (s *DataTestSuite) BeforeTest(suiteName, testName string) {
	s.store.db.Migrator().DropTable(&User{})

	s.store.PopulateTest()
}

func (s *Store) PopulateTest() {
	s.db.AutoMigrate(&User{})

	var us []User

	us = append(us, User{FirstName: "y", LastName: "z", Email: "x@x.x"})

	for _, u := range us {
		u.SetPassword("x")
		s.db.Create(&u)
	}
}
