CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL,
  "hashed_password" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "role" varchar NOT NULL DEFAULT 'buyer',
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "books" (
  "id" bigserial PRIMARY KEY,
  "title" varchar UNIQUE NOT NULL,
  "author" varchar NOT NULL,
  "price" bigint NOT NULL,
  "description" varchar NOT NULL,
  "cover_image_url" varchar NOT NULL,
  "published_date" timestamptz NOT NULL,
  "user_id" bigint,
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "users" ("name");

CREATE INDEX ON "users" ("email");

CREATE INDEX ON "books" ("title");

ALTER TABLE "books" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
