import { Message } from "discord.js";

import logger from "@logger";
import counterSchema from "@schemas/counter";

export default {
  execute: async (message: Message) => {
    const { guild, author, content, channel } = message;

    if (guild == null) return;
    if (author.bot) return;
    if (channel?.type !== "GUILD_TEXT") return;

    const { id: guildId } = guild;
    const { id: channelId } = channel;

    const counter = await counterSchema.findOne({
      guildId,
      channelId,
    });

    if (!counter) {
      logger.verbose(
        `No counter found for guild ${guildId} and channel ${channelId}`
      );
      return;
    }

    if (content !== counter.word) {
      logger.verbose(
        `Counter word ${counter.word} does not match message ${content}`
      );

      await message.delete();

      return;
    }

    counter.counter += 1;
    await counter
      .save()
      .then(async () => {
        logger.verbose(
          `Counter for guild ${guildId} and channel ${channelId} is now ${counter.counter}`
        );
      })
      .catch(async (err) => {
        logger.error(
          `Error saving counter for guild ${guildId} and channel ${channelId}`,
          err
        );
      });

    logger.verbose(
      `Counter word ${counter.word} was found in message ${content} from ${author.tag} (${author.id}) in guild: ${guild?.name} (${guild?.id})`
    );
  },
};