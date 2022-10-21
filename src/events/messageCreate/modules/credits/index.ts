import { ChannelType, Message } from "discord.js";
import { message as CooldownMessage } from "../../../../handlers/cooldown";
import prisma from "../../../../handlers/database";
import logger from "../../../../middlewares/logger";

export default {
  execute: async (message: Message) => {
    const { guild, author, content, channel } = message;

    if (!guild) return;
    if (author.bot) return;
    if (channel.type !== ChannelType.GuildText) return;

    const createGuildMember = await prisma.guildMember.upsert({
      where: {
        userId_guildId: {
          userId: author.id,
          guildId: guild.id,
        },
      },
      update: {},
      create: {
        user: {
          connectOrCreate: {
            create: {
              id: author.id,
            },
            where: {
              id: author.id,
            },
          },
        },
        guild: {
          connectOrCreate: {
            create: {
              id: guild.id,
            },
            where: {
              id: guild.id,
            },
          },
        },
      },
      include: {
        user: true,
        guild: true,
      },
    });

    logger.silly(createGuildMember);

    if (content.length < createGuildMember.guild.creditsMinimumLength) return;

    const isOnCooldown = await CooldownMessage(
      message,
      createGuildMember.guild.creditsTimeout,
      "messageCreate-credits"
    );
    if (isOnCooldown) return;

    const updateGuildMember = await prisma.guildMember.update({
      where: {
        userId_guildId: {
          userId: author.id,
          guildId: guild.id,
        },
      },
      data: {
        creditsEarned: {
          increment: createGuildMember.guild.creditsRate,
        },
      },
    });

    logger.silly(updateGuildMember);

    if (!updateGuildMember)
      throw new Error("Failed to update guildMember object");
  },
};