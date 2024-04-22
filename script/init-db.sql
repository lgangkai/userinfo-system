USE userinfo;

CREATE TABLE `profile_tab`
(
    `id`          bigint unsigned NOT NULL AUTO_INCREMENT,
    `user_id`     bigint unsigned NOT NULL,
    `username`   varchar(255) NOT NULL DEFAULT '',
    `birthday`    DATE,
    `email`       varchar(255) NOT NULL DEFAULT '',
    `avatar_url`  varchar(255) NOT NULL DEFAULT '',
    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY   (`id`),
    KEY           `idx_user_id` (`user_id`),
    UNIQUE KEY    `email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `user_tab`
(
    `id`          bigint unsigned NOT NULL AUTO_INCREMENT,
    `name`        varchar(255) NOT NULL DEFAULT '',
    `password`    varchar(255) NOT NULL DEFAULT '' COMMENT 'encrypted by md5',
    `email`       varchar(255) NOT NULL DEFAULT '',
    `status`      tinyint(3) unsigned NOT NULL DEFAULT 0 COMMENT '0-available, 1-suspended, 2-deleted',
    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY   (`id`),
    UNIQUE KEY    `email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;