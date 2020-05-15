CREATE TABLE IF NOT EXISTS category (
  `id` INTEGER NOT NULL AUTO_INCREMENT,
  `title` VARCHAR(100) NOT NULL,
  `pos` INTEGER NOT NULL,
  `image_url` VARCHAR(512) NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT NOW(),
  `updated_at` TIMESTAMP NOT NULL DEFAULT NOW(),
  PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS product (
  `id` VARCHAR(36) NOT NULL,
  `category_id` INT,
  `title` VARCHAR(100) NOT NULL,
  `image_url` VARCHAR(512) NOT NULL,
  `price` FLOAT NOT NULL,
  `description` TEXT,
  `created_at` TIMESTAMP NOT NULL DEFAULT NOW(),
  `updated_at` TIMESTAMP NOT NULL DEFAULT NOW(),
  PRIMARY KEY (`id`)
);

ALTER TABLE product ADD CONSTRAINT fk_product_category_id FOREIGN KEY (category_id) REFERENCES category(id);

CREATE TABLE IF NOT EXISTS `user` (
  `id` INTEGER NOT NULL AUTO_INCREMENT,
  `username` VARCHAR(100) NOT NULL,
  `password` VARCHAR(60) NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT NOW(),
  PRIMARY KEY (`id`)
);
