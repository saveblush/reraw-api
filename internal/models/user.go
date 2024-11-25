package models

type Timestamp int64

type User struct {
	Pubkey       string    `json:"pubkey" gorm:"primaryKey;type:varchar(64)"`
	CreatedAt    Timestamp `json:"created_at" gorm:"type:integer"`
	UpdatedAt    Timestamp `json:"updated_at" gorm:"type:integer"`
	DeletedAt    Timestamp `json:"deleted_at" gorm:"type:integer"`
	Name         string    `json:"name"`
	LightningURL string    `json:"lightning_url"`
}

func (User) TableName() string {
	return "users"
}
