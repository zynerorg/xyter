package bot

import (
	"log"

	"git.zyner.org/meta/xyter/internal/client"
	"git.zyner.org/meta/xyter/internal/commands"
	"git.zyner.org/meta/xyter/internal/database"
	"git.zyner.org/meta/xyter/internal/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/knadh/koanf/v2"
	"gorm.io/gorm/clause"
)

// Start initializes the Discord session, registers interaction handlers, and creates commands.
func Start(k *koanf.Koanf, db *client.Client) (*discordgo.Session, error) {
	// Create a new Discord session using the provided bot token.
	s, err := discordgo.New("Bot " + k.String("bot/token"))
	if err != nil {
		return nil, err
	}

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		guildID := i.GuildID
		userID := i.Member.User.ID

		// --- Ensure Guild Exists ---
		guild := database.Guild{ID: guildID}
		if err := db.FirstOrCreate(&guild).Error; err != nil {
			utils.SendErrorEmbed(s, i, "Failed to create or find guild.")
			return
		}

		// --- Ensure User Exists ---
		user := database.User{ID: userID}
		if err := db.FirstOrCreate(&user).Error; err != nil {
			utils.SendErrorEmbed(s, i, "Failed to create or find user.")
			return
		}

		// --- Ensure GuildMember Exists (recommended) ---
		member := database.GuildMember{GuildID: guildID, UserID: userID}
		if err := db.FirstOrCreate(&member).Error; err != nil {
			utils.SendErrorEmbed(s, i, "Failed to register guild member.")
			return
		}

		// --- UPSERT GuildSettings (merged table) ---
		settings := database.GuildSettings{
			GuildID: guildID,

			// Defaults (only used on insert)
			DailyBonusAmount:   25,
			WeeklyBonusAmount:  50,
			MonthlyBonusAmount: 150,
			WorkBonusChance:    30,
			WorkPenaltyChance:  10,
		}

		err := db.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "guild_id"}}, // conflict target
			DoUpdates: clause.AssignmentColumns([]string{
				"daily_bonus_amount",
				"weekly_bonus_amount",
				"monthly_bonus_amount",
				"work_bonus_chance",
				"work_penalty_chance",
				"updated_at",
			}),
		}).Create(&settings).Error

		if err != nil {
			utils.SendErrorEmbed(s, i, "Failed to upsert guild settings.")
			return
		}

		// --- Dispatch to actual command handler ---
		handler := commands.ResolveHandler(i)
		if handler != nil {
			handler(s, i, k, db)
			return
		}

		// --- Unknown command fallback ---
		log.Printf("Couldn't find handler for command %q. Removing the command.",
			commands.GetFullCommandName(i),
		)

		err = s.ApplicationCommandDelete(
			k.String("bot/app-id"),
			i.GuildID,
			i.ApplicationCommandData().ID,
		)

		if err != nil {
			log.Fatalf("Failed to remove unknown command. This should be impossible.\n%v", err)
		}

		log.Printf("Command was successfully removed!")
	})

	log.Println("Registering slash commands...")
	for _, cmd := range commands.GetTopLevelCommands() {
		_, err := s.ApplicationCommandCreate(k.String("bot/app-id"), k.String("bot/guild-id"), cmd)
		_, err = s.ApplicationCommandCreate(k.String("bot/app-id"), "", cmd)

		if err != nil {
			log.Printf("Error creating command %q: %v", cmd.Name, err)
		}
	}

	// Open a connection to Discord.
	if err := s.Open(); err != nil {
		return nil, err
	}

	log.Println("Bot is up and running.")
	return s, nil
}
