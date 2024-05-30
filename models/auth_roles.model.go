package models

type AuthRoles struct {
	AuthID uint `json:"auth_id" gorm:"type:uint;primaryKey"`
	RoleID uint `json:"role_id" gorm:"type:uint;primaryKey"`
}
