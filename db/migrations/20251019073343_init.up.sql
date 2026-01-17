DO $$ BEGIN
  CREATE TYPE user_type_t AS ENUM ('SELLER','COLLECTOR','COMPANY');
EXCEPTION
  WHEN duplicate_object THEN null;
END $$;

CREATE TABLE "Address" (
  id BIGSERIAL,
  street_address TEXT NOT NULL,
  city TEXT NOT NULL,
  regency TEXT,
  province TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  PRIMARY KEY (id)
);

CREATE TABLE "User" (
  id BIGSERIAL,
  address_id BIGINT UNIQUE,
  username VARCHAR(128) NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  user_type user_type_t NOT NULL DEFAULT 'SELLER',

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  PRIMARY KEY (id),
  FOREIGN KEY (address_id) REFERENCES "Address"(id) ON DELETE SET NULL
);

CREATE TABLE "Seller" (
  id BIGSERIAL,
  user_id BIGINT NOT NULL UNIQUE,
  seller_name TEXT NOT NULL,

  PRIMARY KEY (id),
  FOREIGN KEY (user_id) REFERENCES "User"(id) ON DELETE CASCADE
);

CREATE TABLE "Collector" (
  id BIGSERIAL,
  user_id BIGINT NOT NULL UNIQUE,
  collector_name TEXT NOT NULL,

  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),

  PRIMARY KEY (id),
  FOREIGN KEY (user_id) REFERENCES "User"(id) ON DELETE CASCADE
);

CREATE TABLE "Company" (
  id BIGSERIAL,
  user_id BIGINT NOT NULL UNIQUE,
  company_name TEXT NOT NULL,

  PRIMARY KEY (id),
  FOREIGN KEY (user_id) REFERENCES "User"(id) ON DELETE CASCADE
);

CREATE TABLE "Oil" (
  id BIGSERIAL,
  collector_id BIGINT NOT NULL UNIQUE,
  total_volume DECIMAL(10,2) NOT NULL DEFAULT 0,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  PRIMARY KEY (id),
  FOREIGN KEY (collector_id) REFERENCES "Collector"(id) ON DELETE CASCADE
);

CREATE TABLE "SellTransaction" (
  id BIGSERIAL,
  seller_id BIGINT NOT NULL,
  collector_id BIGINT NOT NULL,
  volume DECIMAL(10, 2) NOT NULL,
  price DECIMAL(10, 2) NOT NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  PRIMARY KEY (id),
  FOREIGN KEY (seller_id) REFERENCES "Seller"(id) ON DELETE CASCADE,
  FOREIGN KEY (collector_id) REFERENCES "Collector"(id) ON DELETE CASCADE
);

CREATE TABLE "DistributeTransaction" (
  id BIGSERIAL,
  collector_id BIGINT NOT NULL,
  company_id BIGINT NOT NULL,
  volume DECIMAL(10, 2) NOT NULL,
  price DECIMAL(10, 2) NOT NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  PRIMARY KEY (id),
  FOREIGN KEY (collector_id) REFERENCES "Collector"(id) ON DELETE CASCADE,
  FOREIGN KEY (company_id) REFERENCES "Company"(id) ON DELETE CASCADE
);

-- Indexes for Oil table
CREATE INDEX idx_oil_collector_id ON "Oil"(collector_id);

-- Indexes for User table
CREATE INDEX idx_user_created_at ON "User"(created_at DESC);
CREATE INDEX idx_user_type ON "User"(user_type);

-- Indexes for transactions
CREATE INDEX idx_sell_transaction_seller_id ON "SellTransaction"(seller_id);
CREATE INDEX idx_sell_transaction_collector_id ON "SellTransaction"(collector_id);
CREATE INDEX idx_sell_transaction_created_at ON "SellTransaction"(created_at DESC);

CREATE INDEX idx_distribute_transaction_collector_id ON "DistributeTransaction"(collector_id);
CREATE INDEX idx_distribute_transaction_company_id ON "DistributeTransaction"(company_id);
CREATE INDEX idx_distribute_transaction_created_at ON "DistributeTransaction"(created_at DESC);

CREATE OR REPLACE FUNCTION create_user_profile()
RETURNS TRIGGER AS $$
BEGIN
  IF NEW.user_type = 'SELLER' THEN
    INSERT INTO "Seller" (user_id, seller_name)
    VALUES (NEW.id, NEW.username);

  ELSIF NEW.user_type = 'COLLECTOR' THEN
    INSERT INTO "Collector" (user_id, collector_name)
    VALUES (NEW.id, NEW.username);

  ELSIF NEW.user_type = 'COMPANY' THEN
    INSERT INTO "Company" (user_id, company_name)
    VALUES (NEW.id, NEW.username);
  END IF;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_create_user_profile
AFTER INSERT ON "User"
FOR EACH ROW
EXECUTE FUNCTION create_user_profile();

-- Function to create Oil inventory when Collector is created
CREATE OR REPLACE FUNCTION create_oil_inventory()
RETURNS TRIGGER AS $$
BEGIN
  INSERT INTO "Oil" (collector_id, total_volume)
  VALUES (NEW.id, 0)
  ON CONFLICT (collector_id) DO NOTHING;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_create_oil_inventory
AFTER INSERT ON "Collector"
FOR EACH ROW
EXECUTE FUNCTION create_oil_inventory();

-- Function to update Oil inventory when SellTransaction is created
CREATE OR REPLACE FUNCTION update_oil_on_sell()
RETURNS TRIGGER AS $$
BEGIN
  -- Add volume to collector's inventory
  UPDATE "Oil"
  SET total_volume = total_volume + NEW.volume,
      updated_at = NOW()
  WHERE collector_id = NEW.collector_id;

  -- If no row exists, create it
  IF NOT FOUND THEN
    INSERT INTO "Oil" (collector_id, total_volume)
    VALUES (NEW.collector_id, NEW.volume);
  END IF;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_oil_on_sell
AFTER INSERT ON "SellTransaction"
FOR EACH ROW
EXECUTE FUNCTION update_oil_on_sell();

-- Function to update Oil inventory when DistributeTransaction is created
CREATE OR REPLACE FUNCTION update_oil_on_distribute()
RETURNS TRIGGER AS $$
BEGIN
  -- Subtract volume from collector's inventory
  UPDATE "Oil"
  SET total_volume = total_volume - NEW.volume,
      updated_at = NOW()
  WHERE collector_id = NEW.collector_id;

  -- Check if collector has enough oil
  IF (SELECT total_volume FROM "Oil" WHERE collector_id = NEW.collector_id) < 0 THEN
    RAISE EXCEPTION 'Insufficient oil inventory for collector_id %', NEW.collector_id;
  END IF;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_oil_on_distribute
BEFORE INSERT ON "DistributeTransaction"
FOR EACH ROW
EXECUTE FUNCTION update_oil_on_distribute();
