BEGIN;

ALTER TABLE "car" DROP CONSTRAINT IF EXISTS fk_car_brand;
DROP TABLE IF EXISTS "car";
DROP TABLE IF EXISTS "brand";

COMMIT;