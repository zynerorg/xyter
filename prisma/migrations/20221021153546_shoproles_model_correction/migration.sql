-- RedefineTables
PRAGMA foreign_keys=OFF;
CREATE TABLE "new_GuildShopRoles" (
    "guildId" TEXT NOT NULL,
    "roleId" TEXT NOT NULL,
    "userId" TEXT NOT NULL,
    "pricePerHour" INTEGER NOT NULL DEFAULT 5,
    "lastPayed" DATETIME NOT NULL,
    "createdAt" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" DATETIME NOT NULL,
    CONSTRAINT "GuildShopRoles_guildId_fkey" FOREIGN KEY ("guildId") REFERENCES "Guild" ("id") ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT "GuildShopRoles_userId_fkey" FOREIGN KEY ("userId") REFERENCES "User" ("id") ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT "GuildShopRoles_userId_guildId_fkey" FOREIGN KEY ("userId", "guildId") REFERENCES "GuildMember" ("userId", "guildId") ON DELETE RESTRICT ON UPDATE CASCADE
);
INSERT INTO "new_GuildShopRoles" ("createdAt", "guildId", "lastPayed", "pricePerHour", "roleId", "updatedAt", "userId") SELECT "createdAt", "guildId", "lastPayed", "pricePerHour", "roleId", "updatedAt", "userId" FROM "GuildShopRoles";
DROP TABLE "GuildShopRoles";
ALTER TABLE "new_GuildShopRoles" RENAME TO "GuildShopRoles";
CREATE UNIQUE INDEX "GuildShopRoles_guildId_userId_roleId_key" ON "GuildShopRoles"("guildId", "userId", "roleId");
PRAGMA foreign_key_check;
PRAGMA foreign_keys=ON;