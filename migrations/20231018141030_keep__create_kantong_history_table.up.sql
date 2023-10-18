CREATE TABLE `keep_kantong_history` (
    `id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    `kantong_id` INT(10) UNSIGNED NOT NULL,
    `jumlah` DECIMAL(15,0) NOT NULL DEFAULT '0',
    `uraian` VARCHAR(100) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci',
    `waktu` TIMESTAMP NOT NULL DEFAULT current_timestamp(),
    PRIMARY KEY (`id`) USING BTREE
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;
