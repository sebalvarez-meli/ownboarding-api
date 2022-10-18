package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_NewRepository_ReadTimeoutSet_Success(t *testing.T) {
	ass := assert.New(t)

	readTimeout := 1 * time.Second

	config := Config{
		ConnReadTimeout: &readTimeout,
	}
	service, err := NewRepository(config)
	ass.Nil(err)
	ass.NotNil(service)
}

func Test_NewRepository_WriteTimeoutSet_Success(t *testing.T) {
	ass := assert.New(t)

	writeTimeout := 1 * time.Second

	config := Config{
		ConnWriteTimeout: &writeTimeout,
	}
	service, err := NewRepository(config)
	ass.Nil(err)
	ass.NotNil(service)
}

func Test_NewRepository_TimeoutSet_Success(t *testing.T) {
	ass := assert.New(t)

	timeout := 1 * time.Second

	config := Config{
		ConnTimeout: &timeout,
	}
	service, err := NewRepository(config)
	ass.Nil(err)
	ass.NotNil(service)
}

func Test_NewRepository_SetWriteAndTimeout_Success(t *testing.T) {
	ass := assert.New(t)

	timeout := 1 * time.Second
	writeTimeout := 1 * time.Second

	config := Config{
		ConnTimeout:      &timeout,
		ConnWriteTimeout: &writeTimeout,
	}
	service, err := NewRepository(config)
	ass.Nil(err)
	ass.NotNil(service)
}
