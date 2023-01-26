import axios from "axios";
import { Guild } from "discord.js";
import prisma from "../../../handlers/prisma";
import encryption from "../../../helpers/encryption";

export default async (
  guild: Guild,
  code: string,
  amount: number,
  uses: number
) => {
  const upsertGuildConfigApisCpgg = await prisma.guildConfigApisCpgg.upsert({
    where: {
      id: guild.id,
    },
    update: {},
    create: {
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
      guild: true,
    },
  });

  if (!upsertGuildConfigApisCpgg.urlIv || !upsertGuildConfigApisCpgg.urlContent)
    throw new Error("No API url available");
  if (
    !upsertGuildConfigApisCpgg.tokenIv ||
    !upsertGuildConfigApisCpgg.tokenContent
  )
    throw new Error("No API token available");

  const url = encryption.decrypt({
    iv: upsertGuildConfigApisCpgg.urlIv,
    content: upsertGuildConfigApisCpgg.urlContent,
  });
  const api = axios?.create({
    baseURL: `${url}/api/`,
    headers: {
      Authorization: `Bearer ${encryption.decrypt({
        iv: upsertGuildConfigApisCpgg.tokenIv,
        content: upsertGuildConfigApisCpgg.tokenContent,
      })}`,
    },
  });
  const shopUrl = `${url}/store`;

  await api.post("vouchers", {
    uses,
    code,
    credits: amount,
    memo: `Generated by Discord Bot: ${guild.client.user.tag}`,
  });

  return { redeemUrl: `${shopUrl}?voucher=${code}` };
};
