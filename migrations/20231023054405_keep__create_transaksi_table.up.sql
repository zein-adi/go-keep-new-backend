CREATE TABLE `keep_transaksi` (
    `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
    `waktu` TIMESTAMP NOT NULL DEFAULT current_timestamp(),
    `jenis` VARCHAR(15) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci',
    `jumlah` DECIMAL(15,0) NOT NULL DEFAULT '0',
    `pos_asal_id` INT(11) UNSIGNED NULL DEFAULT NULL,
    `pos_asal_nama` VARCHAR(200) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci',
    `pos_tujuan_id` INT(11) UNSIGNED NULL DEFAULT NULL,
    `pos_tujuan_nama` VARCHAR(200) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci',
    `kantong_asal_id` INT(11) UNSIGNED NULL DEFAULT NULL,
    `kantong_asal_nama` VARCHAR(200) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci',
    `kantong_tujuan_id` INT(11) UNSIGNED NULL DEFAULT NULL,
    `kantong_tujuan_nama` VARCHAR(200) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci',
    `uraian` VARCHAR(200) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci',
    `keterangan` TEXT NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci',
    `lokasi` VARCHAR(100) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci',
    `url_foto` TEXT NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci',
    `created_at` TIMESTAMP NOT NULL DEFAULT current_timestamp(),
    `updated_at` TIMESTAMP NOT NULL DEFAULT current_timestamp(),
    `details` LONGTEXT NOT NULL DEFAULT '[]' COLLATE 'utf8mb4_bin',
    `status` VARCHAR(10) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci',
    PRIMARY KEY (`id`) USING BTREE
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;
