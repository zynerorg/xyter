package commands

import (
	"fmt"
	"time"

	"git.zyner.org/meta/xyter/internal/database"
	"git.zyner.org/meta/xyter/internal/database/dal"
	"git.zyner.org/meta/xyter/internal/types/err"
	"git.zyner.org/meta/xyter/internal/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/dromara/carbon/v2"
	"github.com/knadh/koanf/v2"
	"gorm.io/gorm"
)

// EchoCommand defines the Discord application command for /echo.
var CreditsCommand = &discordgo.ApplicationCommand{
	Name:        "credits",
	Description: "Manage your credits.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "bonus",
			Description: "Get bonuses",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "daily",
					Description: "Claim your daily treasure!",
				},
				{
					Name:        "weekly",
					Description: "Claim your weekly treasure!",
				},
				{
					Name:        "monthly",
					Description: "Claim your monthly treasure!",
				},
			},
			Type: discordgo.ApplicationCommandOptionSubCommandGroup,
		},
	},
}

func init() {
	RegisterCommand("credits", "bonus", bonusDailyCmd.Name, bonusDailyCmd, BonusDailyHandler)
	RegisterCommand("credits", "bonus", bonusWeeklyCmd.Name, bonusWeeklyCmd, BonusWeeklyHandler)
	RegisterCommand("credits", "bonus", bonusMonthlyCmd.Name, bonusMonthlyCmd, BonusMonthlyHandler)
}

func claimBonus(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	db *gorm.DB,
	item string,
	bonusAmount int,
	cooldownUntil time.Time,
	title string,
) error {

	utils.DeferResponse(s, i)

	if i.Member.User.Bot {
		return err.ErrBotUser
	}

	guildID := i.GuildID
	userID := i.Member.User.ID

	if guildID == "" || userID == "" {
		utils.SendErrorEmbed(s, i, "Guild or User not found.")
		return nil
	}

	// --- Check Cooldown ---
	var cd database.Cooldown
	err := db.First(&cd,
		"guild_id = ? AND user_id = ? AND item = ?",
		guildID, userID, item,
	).Error

	if err == nil {
		// Cooldown exists -> check if expired
		if time.Now().Before(cd.ExpiresAt) {
			remaining := cd.ExpiresAt.Unix()
			utils.SendErrorEmbed(s, i,
				fmt.Sprintf("You're on cooldown!\nRemaining time: <t:%d:R>", remaining),
			)
			return nil
		}

		// Expired -> delete old cooldown
		db.Delete(&cd)
	}

	// --- Grant Bonus ---
	if err := dal.GiveCredits(db, guildID, userID, bonusAmount); err != nil {
		utils.SendErrorEmbed(s, i, "Failed to grant credits.")
		return nil
	}

	balance, err := dal.GetBalance(db, guildID, userID)
	if err != nil {
		utils.SendErrorEmbed(s, i, "Failed to fetch balance.")
		return nil
	}

	// --- Send Embed Response ---
	embed := &discordgo.MessageEmbed{
		Title:     title,
		Color:     0x00ffcc,
		Author:    &discordgo.MessageEmbedAuthor{Name: title},
		Thumbnail: &discordgo.MessageEmbedThumbnail{URL: i.Member.User.AvatarURL("")},
		Description: fmt.Sprintf(
			"You've received **%d credits**! 🎉\n\n💰 **Your balance**: %d credits",
			bonusAmount, balance,
		),
		Footer: &discordgo.MessageEmbedFooter{
			Text:    fmt.Sprintf("Claimed by %s", i.Member.User.Username),
			IconURL: i.Member.User.AvatarURL(""),
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	utils.SendEmbed(s, i, embed)

	// --- Create New Cooldown ---
	db.Create(&database.Cooldown{
		Item:      item,
		GuildID:   &guildID,
		UserID:    &userID,
		ExpiresAt: cooldownUntil,
	})

	return nil
}

var bonusDailyCmd = &discordgo.ApplicationCommandOption{
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	Name:        "daily",
	Description: "Claim your daily treasure!",
}

func BonusDailyHandler(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	k *koanf.Koanf,
	db *gorm.DB,
) error {

	guildID := i.GuildID
	var settings database.GuildSettings

	if err := db.First(&settings, "guild_id = ?", guildID).Error; err != nil {
		utils.SendErrorEmbed(s, i, "Guild settings not found.")
		return nil
	}

	return claimBonus(
		s,
		i,
		db,
		"daily_bonus",
		settings.DailyBonusAmount,
		carbon.Now().AddDay().StartOfDay().StdTime(),
		"⭐ Daily Treasure Claimed!",
	)
}

var bonusWeeklyCmd = &discordgo.ApplicationCommandOption{
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	Name:        "weekly",
	Description: "Claim your weekly treasure!",
}

func BonusWeeklyHandler(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	k *koanf.Koanf,
	db *gorm.DB,
) error {

	guildID := i.GuildID
	var settings database.GuildSettings

	if err := db.First(&settings, "guild_id = ?", guildID).Error; err != nil {
		utils.SendErrorEmbed(s, i, "Guild settings not found.")
		return nil
	}

	return claimBonus(
		s,
		i,
		db,
		"weekly_bonus",
		settings.WeeklyBonusAmount,
		carbon.Now().AddWeek().StartOfDay().StdTime(),
		"⭐ Weekly Treasure Claimed!",
	)
}

var bonusMonthlyCmd = &discordgo.ApplicationCommandOption{
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	Name:        "monthly",
	Description: "Claim your monthly treasure!",
}

func BonusMonthlyHandler(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	k *koanf.Koanf,
	db *gorm.DB,
) error {

	guildID := i.GuildID
	var settings database.GuildSettings

	if err := db.First(&settings, "guild_id = ?", guildID).Error; err != nil {
		utils.SendErrorEmbed(s, i, "Guild settings not found.")
		return nil
	}

	return claimBonus(
		s,
		i,
		db,
		"monthly_bonus",
		settings.MonthlyBonusAmount,
		carbon.Now().AddMonth().StartOfDay().StdTime(),
		"⭐ Monthly Treasure Claimed!",
	)
}
