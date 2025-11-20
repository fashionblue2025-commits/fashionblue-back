package entities

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// UserRole representa los roles de usuario en el sistema
type UserRole string

const (
	RoleSuperAdmin UserRole = "SUPER_ADMIN"
	RoleSeller     UserRole = "SELLER"
)

// User representa un usuario del sistema (entidad de dominio pura)
type User struct {
	ID        uint
	Email     string
	Password  string
	FirstName string
	LastName  string
	Role      UserRole
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// HashPassword encripta la contraseña del usuario
func (u *User) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword verifica si la contraseña es correcta
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// FullName retorna el nombre completo del usuario
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// IsSuperAdmin verifica si el usuario es super admin
func (u *User) IsSuperAdmin() bool {
	return u.Role == RoleSuperAdmin
}

// IsSeller verifica si el usuario es vendedor
func (u *User) IsSeller() bool {
	return u.Role == RoleSeller
}
