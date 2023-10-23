CREATE TABLE `keep_lokasi` (
    `id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    `nama` VARCHAR(200) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci',
    `last_update` TIMESTAMP NOT NULL DEFAULT current_timestamp(),
    PRIMARY KEY (`id`) USING BTREE
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;
