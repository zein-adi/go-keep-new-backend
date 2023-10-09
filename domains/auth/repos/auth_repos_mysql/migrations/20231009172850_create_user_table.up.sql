CREATE TABLE `users` (
    `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    `username` VARCHAR(100) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci',
    `password` VARCHAR(100) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci',
    `nama` VARCHAR(150) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci',
    `role_ids` LONGTEXT NOT NULL DEFAULT '[]' COLLATE 'utf8mb4_bin',
    PRIMARY KEY (`id`) USING BTREE,
    CONSTRAINT `role_ids` CHECK (json_valid(`role_ids`))
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;
