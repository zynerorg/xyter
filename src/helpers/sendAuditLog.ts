import { ChannelType, EmbedBuilder, Guild } from "discord.js";
import prisma from "../handlers/prisma";
import logger from "../middlewares/logger";
import getEmbedConfig from "./getEmbedConfig";

export default async (guild: Guild, embed: EmbedBuilder) => {
  const getGuildConfigAudits = await prisma.guildConfigAudits.findUnique({
    where: { id: guild.id },
  });
  if (!getGuildConfigAudits) {
    logger.verbose("Guild not found");
    return;
  }

  if (getGuildConfigAudits.status !== true) return;
  if (!getGuildConfigAudits.channelId) {
    throw new Error("Channel not found");
  }

  const embedConfig = await getEmbedConfig(guild);

  embed
    .setTimestamp(new Date())
    .setFooter({
      text: embedConfig.footerText,
      iconURL: embedConfig.footerIcon,
    })
    .setColor(embedConfig.successColor);

  const channel = guild.client.channels.cache.get(
    getGuildConfigAudits.channelId
  );

  if (!channel) throw new Error("Channel not found");
  if (channel.type !== ChannelType.GuildText) {
    throw new Error("Channel must be a text channel");
  }

  await channel
    .send({
      embeds: [embed],
    })
    .then(() => {
      logger.debug(`Audit log sent for event guildMemberAdd`);
    })
    .catch(() => {
      throw new Error("Audit log failed to send");
    });
};