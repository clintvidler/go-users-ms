package data

import (
	"time"

	"github.com/stretchr/testify/assert"
)

func (s *DataTestSuite) TestReadOneToken() {
	_, err := s.store.ReadOneToken("", "")
	assert.Error(s.T(), err, "Invalid email and/or token should fail validation.")

	err = s.store.Login(User{Email: "test@test.com", Password: []byte("password")}, "xyz", time.Now().Add(time.Minute))
	assert.NoError(s.T(), err, "Operation should succeed.")
}

func (s *DataTestSuite) TestDeleteOneToken() {
	err := s.store.DeleteOneToken("")
	assert.Error(s.T(), err, "Email address cannot be blank.")

	s.store.Login(User{Email: "test@test.com", Password: []byte("password")}, "xyz", time.Now().Add(time.Minute))
	err = s.store.DeleteOneToken("test@test.com")
	assert.NoError(s.T(), err, "Operation should succeed.")
}
