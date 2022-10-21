// Dependencies
import { SlashCommandSubcommandBuilder } from "@discordjs/builders";
import { CommandInteraction, EmbedBuilder } from "discord.js";
// Configurations
import getEmbedConfig from "../../../../helpers/getEmbedData";

// Function
export default {
  metadata: { guildOnly: false, ephemeral: false },

  builder: (command: SlashCommandSubcommandBuilder) => {
    return command.setName("ping").setDescription("Ping this bot");
  },
  execute: async (interaction: CommandInteraction) => {
    const { successColor, footerText, footerIcon } = await getEmbedConfig(
      interaction.guild
    );

    const interactionEmbed = new EmbedBuilder()
      .setTitle("[:tools:] Ping")
      .addFields(
        {
          name: "📦 Deliver Latency",
          value: `${Math.abs(Date.now() - interaction.createdTimestamp)} ms`,
          inline: true,
        },
        {
          name: "🤖 API Latency",
          value: `${Math.round(interaction.client.ws.ping)} ms`,
          inline: true,
        }
      )
      .setTimestamp()
      .setColor(successColor)
      .setFooter({ text: footerText, iconURL: footerIcon });

    await interaction.editReply({
      embeds: [interactionEmbed],
    });
  },
};