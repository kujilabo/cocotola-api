package domain

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewSystemOwner(t *testing.T) {
	model, err := NewModel(1, 1, time.Now(), time.Now(), 1, 1)
	assert.NoError(t, err)
	appUser, err := NewAppUser(nil, model, 1, "LOGIN_ID", "USERNAME", nil, nil)
	assert.NoError(t, err)
	systemOwner, err := NewSystemOwner(nil, appUser)
	assert.NoError(t, err)
	log.Println(systemOwner)
}
