-- Database init
CREATE USER smartwitter WITH UNENCRYPTED PASSWORD 'smartwitter';
CREATE DATABASE "smartwitter";
GRANT ALL ON DATABASE "smartwitter" TO "smartwitter";

-- Switch to the smartwitter db
\connect "smartwitter";
set role "smartwitter";

-- Tweet

CREATE TABLE "tweet" (
    "uid" text NOT NULL default '',
    "text" text NOT NULL default '',

    "twitter_user_uid" text NOT NULL,
    "twitter_id" bigint default 0,
    "twitter_creation_time" timestamp with time zone NOT NULL DEFAULT now(),

    "retweet_count" int default 0,
    "favorite_count" int default 0,

    "lang" text default 'en',
    "link" text default 'https://twitter.com',
    "keywords" text[] default array[]::text[],

    -- time
    "creation_time" timestamp with time zone NOT NULL DEFAULT now(),
    "last_update" timestamp with time zone NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX ON "tweet" ("uid");
CREATE INDEX ON "tweet" ("twitter_user_uid");
CREATE INDEX ON "tweet" ("twitter_id");

-- Twitter user

CREATE TABLE "twitter_user" (
    "uid" text NOT NULL default '',

    "twitter_id" bigint default 0,
    "description" text NOT NULL,
    "screen_name" text NOT NULL,
    "name" text NOT NULL,
    "timezone" text NOT NULL,
    "utc_offset" int NOT NULL,
    "followers_count" int default 0,

    -- time
    "creation_time" timestamp with time zone NOT NULL DEFAULT now(),
    "last_update" timestamp with time zone NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX ON "twitter_user" ("uid");
CREATE INDEX ON "twitter_user" ("twitter_id");
ALTER TABLE "tweet" ADD CONSTRAINT tweet_twitter_user FOREIGN KEY ("twitter_user_uid") REFERENCES "twitter_user" ("uid") MATCH FULL;

-- Twitter Future Action

CREATE TABLE "twitter_future_action" (
    "uid" text NOT NULL default '',

    "parameters" json default '[]',

    -- time
    "creation_time" timestamp with time zone NOT NULL DEFAULT now(),
    "execution_time" timestamp with time zone NOT NULL DEFAULT now() + interval '12 hour'
);

----------------------
-- DB Schema
----------------------

CREATE TABLE "db_schema" (
    "version" int NOT NULL DEFAULT 0,
    "update_time" timestamp with time zone NOT NULL DEFAULT now()
);

--

INSERT INTO "db_schema" VALUES (
    1,
    now()
);
