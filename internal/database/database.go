package database

import (
	"log"
	"time"

	"github.com/knadh/koanf/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Guild struct {
	ID        string `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Users    []GuildMember
	Settings GuildSettings `gorm:"constraint:OnDelete:CASCADE"`
	Quotes   []Quote
}

type User struct {
	ID        string `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Guilds       []GuildMember
	Quotes       []Quote `gorm:"foreignKey:AuthorID"`
	PostedQuotes []Quote `gorm:"foreignKey:PosterID"`

	// Instead of a separate table
	ReputationPositive int `gorm:"default:0"`
	ReputationNegative int `gorm:"default:0"`
}

type GuildMember struct {
	GuildID string `gorm:"primaryKey"`
	UserID  string `gorm:"primaryKey"`

	Guild *Guild `gorm:"foreignKey:GuildID;constraint:OnDelete:CASCADE"`
	User  *User  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`

	Balance int

	CreatedAt time.Time
	UpdatedAt time.Time
}

type GuildSettings struct {
	GuildID string `gorm:"primaryKey"`

	CreatedAt time.Time
	UpdatedAt time.Time

	// ---- Credits settings ----
	WorkBonusChance    int `gorm:"default:30"`
	WorkPenaltyChance  int `gorm:"default:10"`
	DailyBonusAmount   int `gorm:"default:25"`
	WeeklyBonusAmount  int `gorm:"default:50"`
	MonthlyBonusAmount int `gorm:"default:150"`

	// ---- Quotes settings ----
	QuotesEnabled   bool `gorm:"default:false"`
	QuotesChannelID string
}

type Quote struct {
	ID        string `gorm:"primaryKey;default:uuid_generate_v4()"`
	CreatedAt time.Time
	UpdatedAt time.Time

	GuildID string
	Guild   Guild `gorm:"constraint:OnDelete:CASCADE"`

	AuthorID string // The person quoted
	Author   User   `gorm:"constraint:OnDelete:SET NULL"`

	Message string

	PosterID string // The user who saved the quote
	Poster   User   `gorm:"constraint:OnDelete:SET NULL"`
}

type CoooldownType string

const (
	CooldownTypeUser        = "user"
	CooldownTypeGuild       = "guild"
	CooldownTypeGuildMember = "guilduser"
)

type Cooldown struct {
	ID        string `gorm:"primaryKey;default:uuid_generate_v4()"`
	CreatedAt time.Time
	UpdatedAt time.Time

	ExpiresAt time.Time
	Item      string

	GuildID *string
	UserID  *string

	Guild *Guild
	User  *User
}

// Token represents an API token
type Token struct {
	ID        string    `gorm:"primaryKey;type:uuid"`
	Hash      string    `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	Revoked   bool      `gorm:"default:false"`
}

// TokenPermission represents a direct permission assigned to a token
type TokenPermission struct {
	ID       string `gorm:"primaryKey;type:uuid"`
	TokenID  string `gorm:"index;not null"`
	Endpoint string `gorm:"not null"`
	Method   string `gorm:"not null"`
}

// Group represents a permission group
type Group struct {
	ID   string `gorm:"primaryKey;type:uuid"`
	Name string `gorm:"uniqueIndex;not null"`
}

// GroupPermission represents a permission assigned to a group
type GroupPermission struct {
	ID       string `gorm:"primaryKey;type:uuid"`
	GroupID  string `gorm:"index;not null"`
	Endpoint string `gorm:"not null"`
	Method   string `gorm:"not null"`
}

// TokenGroup maps a token to a group
type TokenGroup struct {
	ID      string `gorm:"primaryKey;type:uuid"`
	TokenID string `gorm:"index;not null"`
	GroupID string `gorm:"index;not null"`
}

func Open(k *koanf.Koanf) *gorm.DB {
	log.Println(k.String("database/url"))
	db, err := gorm.Open(postgres.Open(k.String("database/url")), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func Migrate(db *gorm.DB) {
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	db.AutoMigrate(
		&Guild{},
		&User{},
		&GuildMember{},
		&GuildSettings{},
		&Quote{},
		&Cooldown{},
		&Token{},
		&TokenPermission{},
		&Group{},
		&GroupPermission{},
		&TokenGroup{},
	)
}
