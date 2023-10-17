CREATE TABLE `pos` (
    `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
    `nama` VARCHAR(100) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci',
    `urutan` SMALLINT(6) NOT NULL DEFAULT '0',
    `saldo` DECIMAL(15,0) NOT NULL DEFAULT '0',
    `parent_id` INT(11) NULL DEFAULT NULL,
    `level` TINYINT(4) NOT NULL DEFAULT '0',
    `is_show` TINYINT(4) NOT NULL DEFAULT '0',
    `is_leaf` TINYINT(4) NOT NULL DEFAULT '0',
    `status` VARCHAR(10) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci',
    PRIMARY KEY (`id`) USING BTREE
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;
