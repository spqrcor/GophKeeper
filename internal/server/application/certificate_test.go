package application

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_createCertificate(t *testing.T) {
	err := createCertificate()
	assert.Nil(t, err)
}
func Test_initCertificate(t *testing.T) {
	err := initCertificate()
	assert.Nil(t, err)
}
