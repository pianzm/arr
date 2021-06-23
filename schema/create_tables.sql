DROP TABLE IF EXISTS "public"."members";
CREATE SEQUENCE article_id_seq;

CREATE TABLE "public"."members" (
  "id" int4 NOT NULL DEFAULT nextval('article_id_seq'::regclass) PRIMARY KEY,
  "email" varchar(255) COLLATE "pg_catalog"."default",
  "first_name" varchar(255) COLLATE "pg_catalog"."default",
  "last_name" varchar(255) COLLATE "pg_catalog"."default",
  "created_at" timestamp(6)
)
;
ALTER TABLE "public"."members" OWNER TO "postgres";
CREATE UNIQUE INDEX "id" ON "public"."members" USING btree (
  "id" "pg_catalog"."int4_ops" ASC NULLS LAST
);

