package repository

type User struct {
	ID       int    `gorm:"primaryKey"`
	FullName string `gorm:"not null"`
	Email    string `gorm:"index;not null"`
	Password string `gorm:"not null"`
	Balance  int    `gorm:"not null;default:0"`
	Author   bool   `gorm:"not null;default:false"`
	Books    []Book `gorm:"foreignKey:AuthorID"`
	Rents    []Rent `gorm:"foreignKey:UserID"`
}

type Category struct {
	ID    int    `gorm:"primaryKey"`
	Name  string `gorm:"not null"`
	Books []Book `gorm:"foreignKey:CategoryID"`
}

type Book struct {
	ID          int      `gorm:"primaryKey"`
	Title       string   `gorm:"not null"`
	Description string   `gorm:"not null"`
	AuthorID    int      `gorm:"index"`
	Author      User     `gorm:"foreignKey:AuthorID"`
	CategoryID  int      `gorm:"index"`
	Category    Category `gorm:"foreignKey:CategoryID"`
	Deposit     int      `gorm:"not null"`
	Available   bool     `gorm:"not null;default:true"`
}

type Rent struct {
	ID     int  `gorm:"primaryKey"`
	BookID int  `gorm:"index"`
	Book   Book `gorm:"foreignKey:BookID"`
	UserID int  `gorm:"index"`
	User   User `gorm:"foreignKey:UserID"`
}
