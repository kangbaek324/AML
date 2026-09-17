-- Kronex Server(Source Database)의 스키마 중 AML에서 참조하는 테이블만 옮겨온 것입니다.
-- 원본: https://github.com/KRONEX-Stock-Exchange/kronex-server/blob/dev/prisma/schema.prisma
-- 이 스키마는 sqlc 타입 생성을 위한 참고용이며, AML이 직접 관리하지 않습니다.

CREATE TABLE users (
    id         INT             NOT NULL PRIMARY KEY AUTO_INCREMENT,
    username   VARCHAR(20)     NOT NULL,
    email      VARCHAR(50)     NOT NULL,
    created_at DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3)
);

CREATE TABLE accounts (
    id                 INT             NOT NULL PRIMARY KEY AUTO_INCREMENT,
    user_id            INT             NOT NULL,
    account_number     INT             NOT NULL,
    balance            BIGINT UNSIGNED NOT NULL,
    available_balance  BIGINT UNSIGNED NOT NULL,
    status             ENUM('PENDING','ACTIVE') NOT NULL DEFAULT 'PENDING',
    published_at       DATETIME(3)     NULL,
    created_at         DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    INDEX idx_accounts_user_id (user_id)
);

CREATE TABLE stocks (
    id             INT             NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name           VARCHAR(30)     NOT NULL,
    price          BIGINT UNSIGNED NOT NULL,
    listing_price  BIGINT UNSIGNED NOT NULL,
    status         ENUM('PENDING','LISTED','SUSPENDED','DELISTED') NOT NULL DEFAULT 'PENDING',
    published_at   DATETIME(3)     NULL,
    created_at     DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at     DATETIME(3)     NULL
);

CREATE TABLE user_stocks (
    account_id          INT             NOT NULL,
    stock_id            INT             NOT NULL,
    quantity            BIGINT UNSIGNED NOT NULL,
    available_quantity  BIGINT UNSIGNED NOT NULL,
    average              BIGINT UNSIGNED NOT NULL,
    total_buy_amount    BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (account_id, stock_id)
);
