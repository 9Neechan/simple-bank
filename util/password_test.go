package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// поверяет, что хеши двух одинаковых паролей разные 
// + проверяет функции HashPassword и CheckPassword
func TestPassword(t *testing.T) {
    password := RandomString(6)
    hashedPassword1, err := HashPassword(password)
    require.NoError(t, err)
    require.NotEmpty(t, hashedPassword1)

    err = CheckPassword(password, hashedPassword1)
    require.NoError(t, err)

    hashedPassword2, err := HashPassword(password)
    require.NoError(t, err)
    require.NotEmpty(t, hashedPassword2)

    require.NotEqual(t, hashedPassword1, hashedPassword2) // тк генерируется разная соль при хешировании
}

func TestWrongPassword(t *testing.T) {
	password := RandomString(6)

	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	wrongPassword := RandomString(6)
	err = CheckPassword(wrongPassword, hashedPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}
