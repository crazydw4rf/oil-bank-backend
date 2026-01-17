DROP TRIGGER IF EXISTS trg_create_user_profile ON "User";
DROP FUNCTION IF EXISTS create_user_profile();

DROP INDEX IF EXISTS idx_oil_user_id;
DROP INDEX IF EXISTS idx_user_created_at;
DROP INDEX IF EXISTS idx_user_type;
DROP INDEX IF EXISTS idx_oil_created_at;

DROP TABLE IF EXISTS "DistributeTransaction";
DROP TABLE IF EXISTS "SellTransaction";
DROP TABLE IF EXISTS "Oil";
DROP TABLE IF EXISTS "Company";
DROP TABLE IF EXISTS "Collector";
DROP TABLE IF EXISTS "Seller";
DROP TABLE IF EXISTS "User";

DROP TYPE IF EXISTS user_type_t;
