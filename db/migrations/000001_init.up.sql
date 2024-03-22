BEGIN TRANSACTION;

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) UNIQUE,
    name VARCHAR(255),
    password VARCHAR(255)
);

CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255),
    price INT,
    image_url VARCHAR(255),
    stock INT,
    user_id UUID REFERENCES users(id),
    condition VARCHAR(255) CHECK (condition IN ('new', 'second')),
    is_purchasable BOOLEAN
);

CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    product_id UUID REFERENCES products(id) ON DELETE CASCADE,
    tag VARCHAR(255)
);

CREATE TABLE bank_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    bank_name VARCHAR(255),
    bank_account_name VARCHAR(255),
    bank_account_number VARCHAR(255)
);

CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    product_id UUID REFERENCES products(id) ON DELETE CASCADE,
    bank_account_id UUID REFERENCES bank_accounts(id) ON DELETE SET NULL,
    payment_proof_image_url VARCHAR(255),
    quantity INT
);

CREATE INDEX idx_product_price ON products (price);

CREATE INDEX idx_tags_product_id ON tags (product_id);

CREATE UNIQUE INDEX idx_user_username ON users (username);

COMMIT TRANSACTION;