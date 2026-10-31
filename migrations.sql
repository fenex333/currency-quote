CREATE TABLE prices;

CREATE TABLE rates (
  id     VARCHAR(20)  NOT NULL,
  pair   VARCHAR(7)   NOT NULL,
  time   BIGINT       NOT NULL,
  price  DECIMAL(8,8) NOT NULL,
  PRIMARY KEY (`id`)
);
