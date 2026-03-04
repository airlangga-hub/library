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

type MailjetRequest struct {
	Messages []MessageRequest `json:"Messages"`
}

type MessageRequest struct {
	From     Person   `json:"From"`
	To       []Person `json:"To"`
	Subject  string   `json:"Subject"`
	TextPart string   `json:"TextPart"`
	HTMLPart string   `json:"HTMLPart"`
}

type Person struct {
	Email string `json:"Email"`
	Name  string `json:"Name"`
}

type MailjetResponse struct {
	Messages     []MessageResponse `json:"Messages"`
	ErrorMessage string            `json:"ErrorMessage"`
	StatusCode   int               `json:"StatusCode"`
}

type MessageResponse struct {
	Status   string `json:"Status"`
	CustomID string `json:"CustomID"`
	To       []To   `json:"To"`
	Cc       []any  `json:"Cc"`
	Bcc      []any  `json:"Bcc"`
}

type To struct {
	Email       string `json:"Email"`
	MessageUUID string `json:"MessageUUID"`
	MessageID   int64  `json:"MessageID"`
	MessageHref string `json:"MessageHref"`
}
