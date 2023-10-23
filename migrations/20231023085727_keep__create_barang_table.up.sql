CREATE TABLE `keep_barang` (
    `id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    `nama` VARCHAR(200) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci',
    `harga` DECIMAL(15,0) NOT NULL DEFAULT '0',
    `diskon` DECIMAL(15,0) NOT NULL DEFAULT '0',
    `satuan_nama` VARCHAR(20) NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci',
    `satuan_jumlah` DECIMAL(15,2) NOT NULL DEFAULT '0.00',
    `satuan_harga` DECIMAL(15,2) NOT NULL DEFAULT '0.00',
    `keterangan` TEXT NOT NULL DEFAULT '' COLLATE 'utf8mb4_general_ci',
    `last_update` TIMESTAMP NOT NULL DEFAULT current_timestamp(),
    `details` LONGTEXT NOT NULL DEFAULT '[]' COLLATE 'utf8mb4_bin',
    PRIMARY KEY (`id`) USING BTREE,
    CONSTRAINT `details` CHECK (json_valid(`details`))
)
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;
