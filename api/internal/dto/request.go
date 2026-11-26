package dto

// LoginRequest digunakan untuk request login (mendukung login via NIP atau Username).
type LoginRequest struct {
	NIP        string `json:"nip"`
	Username   string `json:"username"`
	Password   string `json:"password"   binding:"required"`
	RememberMe bool   `json:"remember_me"`
	ClientIP   string `json:"-"`
	UserAgent  string `json:"-"`
}

// Identifier mengembalikan NIP jika diisi, atau Username sebagai fallback.
func (r LoginRequest) Identifier() string {
	if r.NIP != "" {
		return r.NIP
	}
	return r.Username
}

// StoreUserRequest digunakan untuk request tambah user baru.
type StoreUserRequest struct {
	Username             string `json:"username"              binding:"omitempty,max=100"`
	Name                 string `json:"name"                  binding:"required,max=255"`
	NIP                  string `json:"nip"                   binding:"required,max=30"`
	Email                string `json:"email"                 binding:"required,email,max=255"`
	Password             string `json:"password"              binding:"omitempty,min=8"`
	PasswordConfirmation string `json:"password_confirmation" binding:"omitempty,eqfield=Password"`
	Role                 string `json:"role"                  binding:"required,max=50"`
	UnitKerja            string `json:"unit_kerja"`
	Jabatan              string `json:"jabatan"`
	IsSSO                *bool  `json:"is_sso"`
	IsActive             *bool  `json:"is_active"`
}

// UpdateUserRequest digunakan untuk request update user.
type UpdateUserRequest struct {
	Username             string `json:"username"              binding:"omitempty,max=100"`
	Name                 string `json:"name"                  binding:"required,max=255"`
	NIP                  string `json:"nip"                   binding:"required,max=30"`
	Email                string `json:"email"                 binding:"required,email,max=255"`
	Password             string `json:"password"              binding:"omitempty,min=8"`
	PasswordConfirmation string `json:"password_confirmation" binding:"omitempty,eqfield=Password"`
	Role                 string `json:"role"                  binding:"required,max=50"`
	UnitKerja            string `json:"unit_kerja"`
	Jabatan              string `json:"jabatan"`
	IsSSO                *bool  `json:"is_sso"`
	IsActive             *bool  `json:"is_active"`
}

// UserListRequest digunakan untuk parameter query list user.
type UserListRequest struct {
	Page      int    `form:"page"       binding:"omitempty,min=1"`
	PerPage   int    `form:"per_page"   binding:"omitempty,min=1,max=100"`
	Search    string `form:"search"`
	Role      string `form:"role"       binding:"omitempty,max=50"`
	Status    string `form:"status"     binding:"omitempty,oneof=active inactive"`
	Trashed   string `form:"trashed"    binding:"omitempty,oneof=only with"`
	SortBy    string `form:"sort_by"`
	SortOrder string `form:"sort_order" binding:"omitempty,oneof=asc desc"`
}

// ChangePasswordRequest digunakan untuk request ubah kata sandi mandiri pengguna.
type ChangePasswordRequest struct {
	OldPassword             string `json:"old_password" binding:"required"`
	NewPassword             string `json:"new_password" binding:"required,min=8"`
	NewPasswordConfirmation string `json:"new_password_confirmation" binding:"required,eqfield=NewPassword"`
}
