package model

import "time"

// TODO: ドメイン貧血症を治す
// TODO: JSON タグをドメインモデルではなく、ハンドラー層に定義した構造体に定義するように修正する
type Blog struct {
	ID        uint   `gorm:"primary_key;AUTO_INCREMENT"`
	Title     string `gorm:"title;unique;not null"`
	Count     *int   `gorm:"count;default:0;not null"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (b Blog) TableName() string {
	return "blog_posts"
}
