CREATE TABLE IF NOT EXISTS "users" (
	"user_id"	INTEGER NOT NULL,
	"username"	TEXT NOT NULL UNIQUE,
	"password"	TEXT NOT NULL,
	"email"	TEXT NOT NULL UNIQUE,
	"created"	TEXT NOT NULL,
	PRIMARY KEY("user_id" AUTOINCREMENT)
);

CREATE TABLE IF NOT EXISTS "posts" (
	"post_id"	INTEGER NOT NULL,
	"title"	TEXT NOT NULL,
	"content"	TEXT NOT NULL,
	"user_id"	INTEGER NOT NULL,
	"category_name"	TEXT NOT NULL,
	"created"	TEXT NOT NULL,
	FOREIGN KEY("category_name") REFERENCES "category"("category_name") ON UPDATE NO ACTION,
	FOREIGN KEY("user_id") REFERENCES "users"("user_id"),
	PRIMARY KEY("post_id" AUTOINCREMENT)
);

CREATE TABLE IF NOT EXISTS "comments" (
	"comment_id"	INTEGER NOT NULL,
	"post_id"	INTEGER NOT NULL,
	"user_id"	INTEGER NOT NULL,
	"comment_body"	TEXT NOT NULL,
	"created"	TEXT NOT NULL,
	PRIMARY KEY("comment_id" AUTOINCREMENT),
	FOREIGN KEY("post_id") REFERENCES "posts"("post_id"),
	FOREIGN KEY("user_id") REFERENCES "users"("user_id")
);

CREATE TABLE IF NOT EXISTS "category" (
	"category_id"	INTEGER NOT NULL,
	"category_name"	TEXT NOT NULL UNIQUE,
	PRIMARY KEY("category_id" AUTOINCREMENT)
);

CREATE TABLE IF NOT EXISTS "votes" (
	"id"	INTEGER NOT NULL,
	"user_id"	INTEGER NOT NULL,
	"post_id"	INTEGER NOT NULL,
	"value"	BLOB NOT NULL,
	FOREIGN KEY("user_id") REFERENCES "users"("user_id"),
	FOREIGN KEY("post_id") REFERENCES "posts"("post_id"),
	PRIMARY KEY("id" AUTOINCREMENT)
);