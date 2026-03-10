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

	Users    []GuildUser
	Settings GuildSettings `gorm:"constraint:OnDelete:CASCADE"`
	Quotes   []Quote
}

type User struct {
	ID        string `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Guilds       []GuildUser
	Quotes       []Quote `gorm:"foreignKey:AuthorID"`
	PostedQuotes []Quote `gorm:"foreignKey:PosterID"`

	// Instead of a separate table
	ReputationPositive int `gorm:"default:0"`
	ReputationNegative int `gorm:"default:0"`
}

type GuildUser struct {
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
	CooldownTypeUser      = "user"
	CooldownTypeGuild     = "guild"
	CooldownTypeGuildUser = "guilduser"
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
		&GuildUser{},
		&GuildSettings{},
		&Quote{},
		&Cooldown{},
	)
}
