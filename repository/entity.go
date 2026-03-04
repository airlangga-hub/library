package repository

type User struct {
	ID       int    `gorm:"primaryKey"`
	Email    string `gorm:"index"`
	Password string
	Balance  int
	Author   bool
	Books    []Book `gorm:"foreignKey:AuthorID"`
	Rents    []Rent `gorm:"foreignKey:UserID"`
}

type Category struct {
	ID    int `gorm:"primaryKey"`
	Name  string
	Books []Book `gorm:"foreignKey:CategoryID"`
}

type Book struct {
	ID          int `gorm:"primaryKey"`
	Title       string
	Description string
	AuthorID    int      `gorm:"index"`
	Author      User     `gorm:"foreignKey:AuthorID"`
	CategoryID  int      `gorm:"index"`
	Category    Category `gorm:"foreignKey:CategoryID"`
	Deposit     int
	Available   bool
}

type Rent struct {
	ID     int  `gorm:"primaryKey"`
	BookID int  `gorm:"index"`
	Book   Book `gorm:"foreignKey:BookID"`
	UserID int  `gorm:"index"`
	User   User `gorm:"foreignKey:UserID"`
}
