CREATE SCHEMA IF NOT EXISTS "public";

CREATE TABLE "public"."users" (
    "id" bigserial NOT NULL,
    "username" varchar(100) NOT NULL UNIQUE,
    "password" text NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    CONSTRAINT "pk_users_id" PRIMARY KEY ("id")
);

CREATE TABLE "public"."gifts" (
    "id" bigserial NOT NULL,
    "title" text,
    "description" text,
    "stock" int,
    "redeem_point" int NOT NULL,
    "image_url" text,
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    CONSTRAINT "pk_gifts_id" PRIMARY KEY ("id")
);

CREATE TABLE "public"."gift_ratings" (
    "id" bigserial NOT NULL,
    "user_id" bigserial NOT NULL,
    "gift_id" bigserial NOT NULL,
    "rate" numeric,
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    CONSTRAINT "pk_gift_ratings_id" PRIMARY KEY ("id")
);

-- Foreign key constraints
-- Schema: public
ALTER TABLE "public"."gift_ratings" ADD CONSTRAINT "fk_gift_ratings_user_id_users_id" FOREIGN KEY("user_id") REFERENCES "public"."users"("id");
ALTER TABLE "public"."gift_ratings" ADD CONSTRAINT "fk_gift_ratings_gift_id_gifts_id" FOREIGN KEY("gift_id") REFERENCES "public"."gifts"("id");