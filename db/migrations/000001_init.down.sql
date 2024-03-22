-- Drop unique index on USER table for username
DROP INDEX idx_user_username;

-- Drop index on TAGS table for product_id
DROP INDEX idx_tags_product_id;

-- Drop index on PRODUCT table for price
DROP INDEX idx_product_price;

-- Drop TABLES

DROP TABLE IF EXISTS BANK_ACCOUNT;
DROP TABLE IF EXISTS PAYMENT;
DROP TABLE IF EXISTS TAGS;
DROP TABLE IF EXISTS PRODUCT;
DROP TABLE IF EXISTS USER;
