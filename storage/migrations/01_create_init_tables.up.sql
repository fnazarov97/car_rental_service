BEGIN;

CREATE TABLE IF NOT EXISTS "brand" (
	"brand_id" CHAR(36) NOT NULL PRIMARY KEY,
	"name" VARCHAR(255) UNIQUE NOT NULL,
	"discription" TEXT NOT NULL,
	"year" CHAR(36) NOT NULL,
	"created_at" TIMESTAMP DEFAULT now(),
	"updated_at" TIMESTAMP,
	"deleted_at" TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "car" (
	"car_id" CHAR(36) NOT NULL PRIMARY KEY,
	"model" VARCHAR(255) NOT NULL,
	"color" VARCHAR(255) NOT NULL,
	"year"  VARCHAR(255) NOT NULL,
	"mileage" VARCHAR(255) NOT NULL,
	"brand_id" CHAR(36) REFERENCES "brand" (brand_id)
	"created_at" TIMESTAMP DEFAULT now(),
	"updated_at" TIMESTAMP,
	"deleted_at" TIMESTAMP
);

COMMIT;