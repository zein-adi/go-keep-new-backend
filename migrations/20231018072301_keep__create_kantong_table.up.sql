CREATE TABLE `keep_kantong` (
    `id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    `nama` VARCHAR(100) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci',
    `urutan` SMALLINT(6) NOT NULL DEFAULT '0',
    `saldo` DECIMAL(15,0) NOT NULL DEFAULT '0',
    `saldo_mengendap` DECIMAL(15,0) NOT NULL DEFAULT '0',
    `pos_id` INT(10) UNSIGNED NULL DEFAULT NULL,
    `is_show` TINYINT(4) NOT NULL DEFAULT '0',
    `status` VARCHAR(50) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci',
    PRIMARY KEY (`id`) USING BTREE
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;
