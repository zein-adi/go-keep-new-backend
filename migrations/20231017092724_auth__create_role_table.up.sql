CREATE TABLE `roles` (
    `id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    `nama` VARCHAR(100) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci',
    `deskripsi` TEXT NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci',
    `level` SMALLINT(6) UNSIGNED NOT NULL DEFAULT '0',
    `permissions` LONGTEXT NOT NULL DEFAULT '[]' COLLATE 'utf8mb4_bin',
    PRIMARY KEY (`id`) USING BTREE,
    CONSTRAINT `permissions` CHECK (json_valid(`permissions`))
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;
