package database

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/liyansasongko/bumdes-be/internal/entity"
)

func Seed(db *gorm.DB) error {
	if err := seedPermissions(db); err != nil {
		return err
	}
	if err := seedRoles(db); err != nil {
		return err
	}
	if err := seedUsers(db); err != nil {
		return err
	}
	return nil
}

func seedPermissions(db *gorm.DB) error {
	permissions := []entity.Permission{
		{Code: "user.view", Name: "View Users", Description: "Melihat daftar akun/user"},
		{Code: "user.manage", Name: "Manage Users", Description: "Mengelola akun/user"},
		{Code: "role.view", Name: "View Roles", Description: "Melihat daftar role"},
		{Code: "role.manage", Name: "Manage Roles", Description: "Mengelola role dan permission"},
	}

	for _, p := range permissions {
		if err := db.Where("code = ?", p.Code).FirstOrCreate(&p).Error; err != nil {
			return err
		}
	}

	return nil
}

func seedRoles(db *gorm.DB) error {
	var allPermissions []entity.Permission
	if err := db.Find(&allPermissions).Error; err != nil {
		return err
	}

	var viewPermissions []entity.Permission
	if err := db.Where("code IN ?", []string{"user.view", "role.view"}).Find(&viewPermissions).Error; err != nil {
		return err
	}

	roleDefs := []struct {
		Name        string
		Description string
		Permissions []entity.Permission
	}{
		{Name: "superadmin", Description: "Superadmin dengan akses ke seluruh permission", Permissions: allPermissions},
		{Name: "admin", Description: "Administrator dengan akses penuh", Permissions: allPermissions},
		{Name: "member", Description: "Pengguna biasa dengan akses terbatas", Permissions: viewPermissions},
		{Name: "Bendahara", Description: "Mengelola keuangan"},
		{Name: "Petugas Parkir", Description: "Mengelola area parkir"},
		{Name: "Petugas Tiket", Description: "Mengelola penjualan tiket"},
		{Name: "Petugas Wahana", Description: "Mengelola wahana"},
		{Name: "Viewer", Description: "Hanya dapat melihat data, tanpa akses mengubah"},
	}

	for _, def := range roleDefs {
		var role entity.Role
		err := db.Where("name = ?", def.Name).First(&role).Error
		if err == nil {
			continue
		}
		if err != gorm.ErrRecordNotFound {
			return err
		}

		role = entity.Role{
			Name:        def.Name,
			Description: def.Description,
			Permissions: def.Permissions,
		}
		if err := db.Create(&role).Error; err != nil {
			return err
		}
	}

	return nil
}

func seedUsers(db *gorm.DB) error {
	userDefs := []struct {
		Name     string
		Email    string
		Password string
		RoleName string
	}{
		{Name: "Superadmin", Email: "superadmin@bumdes.local", Password: "password123", RoleName: "superadmin"},
		{Name: "Administrator", Email: "admin@bumdes.local", Password: "password123", RoleName: "admin"},
	}

	for _, def := range userDefs {
		var count int64
		if err := db.Model(&entity.User{}).Where("email = ?", def.Email).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			continue
		}

		var role entity.Role
		if err := db.Where("name = ?", def.RoleName).First(&role).Error; err != nil {
			return err
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(def.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		user := entity.User{
			Name:     def.Name,
			Email:    def.Email,
			Password: string(hashed),
			RoleID:   role.ID,
			IsActive: true,
		}
		if err := db.Create(&user).Error; err != nil {
			return err
		}
	}

	return nil
}
