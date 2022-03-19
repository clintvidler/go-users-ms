package data

import (
	"time"

	"github.com/stretchr/testify/assert"
)

func (s *DataTestSuite) TestSetPassword() {
	user, _ := s.store.ReadOneUser("x@x.x")

	err := user.SetPassword("xx")

	assert.NoError(s.T(), err, "Failed to 'SetPassword'.")
}

func (s *DataTestSuite) TestComparePassword() {
	user, _ := s.store.ReadOneUser("x@x.x")

	err := user.ComparePassword("x")
	assert.NoError(s.T(), err, "Passwords should match.")

	err = user.ComparePassword("xx")
	assert.Error(s.T(), err, "Passwords should not match.")
}

func (s *DataTestSuite) TestCreateOneUser() {
	var user User

	user.Email = "xx@x.x"
	user.FirstName = "yy"
	user.LastName = "zz"
	user.SetPassword("x")

	err := s.store.CreateOneUser(&user)
	assert.NoError(s.T(), err, "'CreateOneUser' should succeed.")

	err = s.store.CreateOneUser(&user)
	assert.Error(s.T(), err, "'CreateOneUser' should not fail due to non-unique email.")

	user2, _ := s.store.ReadOneUser(user.Email)

	assert.Equal(s.T(), user.FirstName, user2.FirstName)
	assert.Equal(s.T(), user.LastName, user2.LastName)
	assert.NotEmpty(s.T(), user.Password, user2.Password)
}

func (s *DataTestSuite) TestReadOneUser() {
	user, err := s.store.ReadOneUser("x@x.x")

	assert.NoError(s.T(), err, "Failed to 'ReadOneUser'.")

	assert.Equal(s.T(), user.FirstName, "y")
	assert.Equal(s.T(), user.LastName, "z")
	assert.NotEmpty(s.T(), user.Password)
}

func (s *DataTestSuite) TestLogin() {
	err := s.store.Login(User{}, "", time.Now())
	assert.Error(s.T(), err, "User must have an email.")

	err = s.store.Login(User{Email: "test@test.com"}, "", time.Now())
	assert.Error(s.T(), err, "User must have a password.")

	err = s.store.Login(User{Email: "test@test.com", Password: []byte("password")}, "", time.Now())
	assert.Error(s.T(), err, "Token is required.")

	err = s.store.Login(User{Email: "test@test.com", Password: []byte("password")}, "xyz", time.Now().Add(time.Minute*-1))
	assert.Error(s.T(), err, "Expires must be in the future.")

	err = s.store.Login(User{Email: "test@test.com", Password: []byte("password")}, "xyz", time.Now().Add(time.Minute))
	assert.NoError(s.T(), err, "Operation should succeed.")
}

func (s *DataTestSuite) TestLogout() {
	err := s.store.Logout(User{})
	assert.Error(s.T(), err, "User must have an email.")

	err = s.store.Logout(User{Email: "404@test.com"})
	assert.Error(s.T(), err, "User must exist.")

	s.store.CreateOneUser(&User{Email: "test@test.com", Password: []byte("password")})

	err = s.store.Logout(User{Email: "test@test.com", Password: []byte("password")})
	assert.NoError(s.T(), err, "Operation should succeed.")
}
