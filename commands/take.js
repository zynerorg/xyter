const { SlashCommandBuilder } = require('@discordjs/builders');
const { Permissions } = require('discord.js');

const db = require('quick.db');

const credits = new db.table('credits');

module.exports = {
  data: new SlashCommandBuilder()
    .setName('take')
    .setDescription('Take credits from a user (ADMIN).')
    .addUserOption((option) =>
      option
        .setName('user')
        .setDescription('The user you want to take credits from.')
        .setRequired(true)
    )
    .addIntegerOption((option) =>
      option.setName('amount').setDescription('The amount you will take.').setRequired(true)
    ),
  async execute(interaction) {
    if (!interaction.member.permissions.has(Permissions.FLAGS.MANAGE_GUILD)) {
      const embed = {
        title: 'Take',
        description: 'You need to have permission to manage this guild (MANAGE_GUILD)',
        color: 0xbb2124,
        timestamp: new Date(),
        footer: { text: 'Zyner Bot' },
      };
      return await interaction.reply({ embeds: [embed], ephemeral: true });
    }
    const user = await interaction.options.getUser('user');
    const amount = await interaction.options.getInteger('amount');

    if (amount <= 0) {
      const embed = {
        title: 'Take',
        description: "You can't take zero or below.",
        color: 0xbb2124,
        timestamp: new Date(),
        footer: { text: 'Zyner Bot' },
      };
      return await interaction.reply({ embeds: [embed], ephemeral: true });
    } else {
      await credits.subtract(user.id, amount);

      const embed = {
        title: 'Take',
        description: `You took ${
          amount <= 1 ? `${amount} credit` : `${amount} credits`
        } to ${user}.`,
        color: 0x22bb33,
        timestamp: new Date(),
        footer: { text: 'Zyner Bot' },
      };
      return await interaction.reply({ embeds: [embed], ephemeral: true });
    }
  },
};
