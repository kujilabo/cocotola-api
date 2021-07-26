package domain

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewSystemOwner(t *testing.T) {
	model := NewModel(0, 0, time.Now(), time.Now(), 0, 0)
	appUser, err := NewAppUser(nil, model, 1, "LOGIN_ID", "USERNAME", nil, nil)
	assert.NoError(t, err)
	systemOwner, err := NewSystemOwner(nil, appUser)
	assert.NoError(t, err)
	fmt.Println(systemOwner)
}
