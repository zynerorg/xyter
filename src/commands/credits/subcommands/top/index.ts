import { GuildMemberCredits } from "@prisma/client";
import {
  CommandInteraction,
  SlashCommandSubcommandBuilder,
  userMention,
} from "discord.js";

import prisma from "../../../../handlers/database";
import deferReply from "../../../../handlers/deferReply";
import { success as BaseEmbedSuccess } from "../../../../helpers/baseEmbeds";
import upsertGuildMember from "../../../../helpers/upsertGuildMember";
import logger from "../../../../middlewares/logger";

// 1. Export a builder function.
export const builder = (command: SlashCommandSubcommandBuilder) => {
  return command.setName("top").setDescription(`View the top users`);
};

// 2. Export an execute function.
export const execute = async (interaction: CommandInteraction) => {
  // 1. Defer reply as permanent.
  await deferReply(interaction, false);

  // 2. Destructure interaction object.
  const { guild, client, user } = interaction;
  if (!guild) throw new Error("Guild not found");
  if (!client) throw new Error("Client not found");

  await upsertGuildMember(guild, user);

  // 3. Create base embeds.
  const EmbedSuccess = await BaseEmbedSuccess(guild, "[:dollar:] Top");

  // 4. Get the top 10 users.
  const topTen = await prisma.guildMemberCredits.findMany({
    where: {
      guildId: guild.id,
    },
    orderBy: {
      balance: "desc",
    },
    take: 10,
  });
  logger.silly(topTen);

  // 5. Create the top 10 list.
  const entry = (guildMemberCredits: GuildMemberCredits, index: number) =>
    `${index + 1}. ${userMention(guildMemberCredits.userId)} | :coin: ${
      guildMemberCredits.balance
    }`;

  // 6. Send embed
  return interaction.editReply({
    embeds: [
      EmbedSuccess.setDescription(
        `The top 10 users in this server are:\n\n${topTen
          .map(entry)
          .join("\n")}`
      ),
    ],
  });
};
