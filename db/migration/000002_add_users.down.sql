DROP INDEX "accounts_user_name_currency_idx";

DROP INDEX "accounts_user_name_idx";

ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "accounts_user_name_fkey";
DROP TABLE IF EXISTS "users";