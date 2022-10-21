import { SlashCommandSubcommandBuilder } from "@discordjs/builders";
import axios from "axios";
import { CommandInteraction, EmbedBuilder } from "discord.js";
import { command as CooldownCommand } from "../../../../handlers/cooldown";
import deferReply from "../../../../handlers/deferReply";
import getEmbedConfig from "../../../../helpers/getEmbedData";

export default {
  builder: (command: SlashCommandSubcommandBuilder) => {
    return command.setName("meme").setDescription("Get a meme from r/memes)");
  },

  execute: async (interaction: CommandInteraction) => {
    await deferReply(interaction, false);

    await CooldownCommand(interaction, 15);

    const { guild } = interaction;

    const embedConfig = await getEmbedConfig(guild);

    await axios
      .get("https://www.reddit.com/r/memes/random/.json")
      .then(async (res) => {
        const response = res.data[0].data.children;
        const content = response[0].data;

        const embed = new EmbedBuilder()
          .setAuthor({
            name: content.title,
            iconURL:
              "https://www.redditinc.com/assets/images/site/reddit-logo.png",
            url: `https://reddit.com${content.permalink}`,
          })
          .setTitle("[:sweat_smile:] Meme")
          .addFields([
            {
              name: "Author",
              value: `[${content.author}](https://reddit.com/user/${content.author})`,
              inline: true,
            },
            {
              name: "Votes",
              value: `${content.ups}/${content.downs}`,
              inline: true,
            },
          ])
          .setTimestamp(new Date())
          .setImage(content.url)
          .setFooter({
            text: embedConfig.footerText,
            iconURL: embedConfig.footerIcon,
          })
          .setColor(embedConfig.successColor);

        await interaction.editReply({ embeds: [embed] });
        return;
      })
      .catch((error) => {
        throw new Error(error.message);
      });
  },
};