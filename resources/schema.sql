-- Database init
CREATE USER smartwitter WITH UNENCRYPTED PASSWORD 'smartwitter';
CREATE DATABASE "smartwitter";
GRANT ALL ON DATABASE "smartwitter" TO "smartwitter";

-- Switch to the smartwitter db
\connect "smartwitter";
set role "smartwitter";

-- User

CREATE TABLE "user" (
    "uid" text NOT NULL,

    "twitter_secret" text DEFAULT NULL,
    "twitter_token" text DEFAULT NULL,
    "twitter_id" text DEFAULT NULL,
    "twitter_username" text DEFAULT NULL,
    "twitter_name" text DEFAULT NULL,

    "session_token" text DEFAULT NULL,

    -- emailing
    "unsubscribe_token" text DEFAULT '',

    -- payment
    "stripe_token" text DEFAULT '',

    "creation_time" timestamp with time zone NOT NULL default now(),
    "last_login" timestamp with time zone NOT NULL default now()
);

CREATE UNIQUE INDEX ON "user" ("uid");

-- Tweet

CREATE TABLE "tweet" (
    "uid" text NOT NULL default '',
    "text" text NOT NULL default '',

    "twitter_user_uid" text NOT NULL,
    "twitter_id" text NOT NULL,
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

-- Twitter Keywords Watcher

CREATE TABLE "twitter_keywords_watcher" (
    "uid" text NOT NULL,
    "user_uid" text NOT NULL,
    "position" int NOT NULL DEFAULT 0,
    "keywords" text NOT NULL,

    "last_run" timestamp with time zone,

    "creation_time" timestamp with time zone NOT NULL DEFAULT now(),
    "last_update" timestamp with time zone NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX ON "twitter_keywords_watcher" ("uid")
CREATE INDEX ON "twitter_keywords_watcher" ("user_uid");
ALTER TABLE "twitter_keywords_watcher" ADD CONSTRAINT tkw_user FOREIGN KEY ("user_uid") REFERENCES "user" ("uid") MATCH FULL;

-- Twitter user

CREATE TABLE "twitter_user" (
    "uid" text NOT NULL default '',

    "twitter_id" text NOT NULL,
    "description" text NOT NULL,
    "screen_name" text NOT NULL,
    "name" text NOT NULL,
    "avatar" text DEFAULT 'https://abs.twimg.com/sticky/default_profile_images/default_profile_normal.png',
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

CREATE TABLE "twitter_planned_action" (
    "uid" text NOT NULL default '',

    "type" text NOT NULL,
    "user_uid" text NOT NULL, -- author of the action
    "tweet_id" text NOT NULL,

    -- time
    "creation_time" timestamp with time zone NOT NULL DEFAULT now(),
    "execution_time" timestamp with time zone NOT NULL DEFAULT now() + interval '12 hour',
    "done" timestamp with time zone DEFAULT NULL
);

ALTER TABLE "twitter_planned_action" ADD CONSTRAINT tpa_user FOREIGN KEY ("user_uid") REFERENCES "user" ("uid") MATCH FULL;

-- Tweet State

CREATE TABLE "tweet_done_action" (
    "user_uid" text NOT NULL,
    "tweet_id" text NOT NULL,
    "ignored_time" timestamp with time zone DEFAULT NULL,
    "liked_time" timestamp with time zone DEFAULT NULL,
    "retweeted_time" timestamp with time zone DEFAULT NULL
);

CREATE UNIQUE INDEX ON "tweet_done_action" ("user_uid", "tweet_id");
CREATE INDEX ON "tweet_done_action" ("tweet_id");
ALTER TABLE "tweet_done_action" ADD CONSTRAINT tda_user FOREIGN KEY ("user_uid") REFERENCES "user" ("uid") MATCH FULL;

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
