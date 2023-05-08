CREATE TYPE "subscription_status" AS ENUM (
  'active',
  'inactive'
);

CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "stripe_id" varchar UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "packages" (
  "id" varchar PRIMARY KEY,
  "name" varchar NOT NULL,
  "description" varchar NOT NULL,
  "price" bigint NOT NULL,
  "stripe_price_id" varchar NOT NULL
);

CREATE TABLE "users_packages" (
  "id" varchar PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "package_id" varchar NOT NULL,
  "status" subscription_status NOT NULL,
  "start_date" timestamptz NOT NULL DEFAULT (now()),
  "end_date" timestamptz NOT NULL DEFAULT (now())
);

CREATE UNIQUE INDEX ON "users_packages" ("user_id", "package_id");

ALTER TABLE "users_packages" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "users_packages" ADD FOREIGN KEY ("package_id") REFERENCES "packages" ("id");