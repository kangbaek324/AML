CREATE TABLE users (
    id             INT UNSIGNED    NOT NULL PRIMARY KEY,
    average_asset  DECIMAL(18,2)   NOT NULL DEFAULT 0,
    risk_level     ENUM('LOW','MEDIUM','HIGH') NOT NULL DEFAULT 'LOW',
    asset_tier     ENUM('LOW','MEDIUM','HIGH') NOT NULL DEFAULT 'LOW',
    updated_at     DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE alerts (
    id           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    userId       INT UNSIGNED    NOT NULL,                             
    type         ENUM('CROSS_TRADING','LARGE_TRANSACTION','ABNORMAL_TRANSFER') NOT NULL,
    reason       TEXT            NOT NULL,
    status       ENUM('PENDING','NORMAL','ABNORMAL') NOT NULL DEFAULT 'PENDING',
    alerted_at   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    processed_at DATETIME        NULL,                                
    INDEX idx_alerts_userId (userId),
    INDEX idx_alerts_status (status)
);

CREATE TABLE alert_trades (
    id         INT UNSIGNED    NOT NULL AUTO_INCREMENT PRIMARY KEY,
    alertId    BIGINT UNSIGNED NOT NULL,
    tradeId    BIGINT UNSIGNED NOT NULL,
    created_at DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_alert_trades_alertId (alertId),
    INDEX idx_alert_trades_tradeId (tradeId)
);

CREATE TABLE alert_transfers (
    id         INT UNSIGNED    NOT NULL AUTO_INCREMENT PRIMARY KEY,
    alertId    BIGINT UNSIGNED NOT NULL,
    transferId BIGINT UNSIGNED NOT NULL,
    created_at DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_alert_transfers_alertId (alertId),
    INDEX idx_alert_transfers_transferId (transferId)
);

CREATE TABLE cross_trading_count (
    id         INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    userId     INT UNSIGNED NOT NULL,
    date       DATE         NOT NULL,
    count      INT UNSIGNED NOT NULL DEFAULT 0,
    alerted    BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uq_user_date (userId, date)
);

CREATE TABLE cursors (
    id        INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    type      VARCHAR(50)  NOT NULL,
    timestamp DATETIME(3)  NOT NULL,
    UNIQUE KEY uq_cursors_type (type)
);
