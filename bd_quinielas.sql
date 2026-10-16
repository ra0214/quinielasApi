CREATE DATABASE IF NOT EXISTS `quinielas` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
USE `quinielas`;

-- 1. Tabla Clientes
CREATE TABLE IF NOT EXISTS `clientes` (
  `id_cliente` INT NOT NULL AUTO_INCREMENT,
  `nombre` VARCHAR(100) NOT NULL,
  `telefono` VARCHAR(20) DEFAULT NULL,
  `fecha_registro` TIMESTAMP NULL DEFAULT (NOW()),
  PRIMARY KEY (`id_cliente`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 2. Tabla Ediciones
CREATE TABLE IF NOT EXISTS `ediciones` (
  `id_edicion` INT NOT NULL AUTO_INCREMENT,
  `tipo_edicion` VARCHAR(50) NOT NULL,
  `nombre_edicion` VARCHAR(100) NOT NULL,
  `fecha_inicio` TIMESTAMP NULL DEFAULT (NOW()),
  `fecha_limite` TIMESTAMP NULL DEFAULT NULL,
  PRIMARY KEY (`id_edicion`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 3. Tabla Quinielas
CREATE TABLE IF NOT EXISTS `quinielas` (
  `id_quiniela` INT NOT NULL AUTO_INCREMENT,
  `id_edicion` INT NOT NULL,
  `nombre_variante` VARCHAR(50) NOT NULL,
  `monto_meta` DECIMAL(12,2) NOT NULL DEFAULT '0.00',
  `fecha_limite` TIMESTAMP NOT NULL,
  `estado` VARCHAR(20) DEFAULT 'EN_JUEGO',
  `tipo_quiniela` VARCHAR(30) DEFAULT NULL,
  `id_cliente` INT DEFAULT NULL,
  `monto_pagado` DECIMAL(12,2) NOT NULL DEFAULT '0.00',
  `analista` VARCHAR(100) DEFAULT NULL,
  `tipo_pago` VARCHAR(30) DEFAULT NULL,
  `canal_venta` VARCHAR(20) DEFAULT NULL,
  `sujeto_cancelacion` TINYINT(1) NOT NULL DEFAULT 1,
  `estado_fallo` TINYINT NOT NULL DEFAULT 0,
  `monto_meta_anterior` DECIMAL(12,2) DEFAULT NULL,
  `fecha_meta_anterior` DATETIME DEFAULT NULL,
  PRIMARY KEY (`id_quiniela`),
  KEY `id_edicion` (`id_edicion`),
  KEY `idx_quinielas_id_cliente` (`id_cliente`),
  CONSTRAINT `quinielas_ibfk_1` FOREIGN KEY (`id_edicion`) REFERENCES `ediciones` (`id_edicion`),
  CONSTRAINT `fk_quinielas_cliente` FOREIGN KEY (`id_cliente`) REFERENCES `clientes` (`id_cliente`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 4. Tabla Aportes
CREATE TABLE IF NOT EXISTS `aportes` (
  `id_aporte` INT NOT NULL AUTO_INCREMENT,
  `id_cliente` INT NOT NULL,
  `id_quiniela` INT NOT NULL,
  `monto_total_acumulado` DECIMAL(12,2) NOT NULL DEFAULT '0.00',
  `porcentaje_participacion` DECIMAL(7,4) NOT NULL DEFAULT '0.0000',
  PRIMARY KEY (`id_aporte`),
  KEY `id_cliente` (`id_cliente`),
  KEY `id_quiniela` (`id_quiniela`),
  CONSTRAINT `aportes_ibfk_1` FOREIGN KEY (`id_cliente`) REFERENCES `clientes` (`id_cliente`),
  CONSTRAINT `aportes_ibfk_2` FOREIGN KEY (`id_quiniela`) REFERENCES `quinielas` (`id_quiniela`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 5. Tabla Premios
CREATE TABLE IF NOT EXISTS `premios` (
  `id_premio` INT NOT NULL AUTO_INCREMENT,
  `id_quiniela` INT NOT NULL,
  `monto_bruto` DECIMAL(12,2) NOT NULL DEFAULT '0.00',
  `porcentaje_retencion` DECIMAL(5,4) NOT NULL DEFAULT '0.0000',
  `monto_neto` DECIMAL(12,2) NOT NULL DEFAULT '0.00',
  `fecha_registro` TIMESTAMP NULL DEFAULT (NOW()),
  PRIMARY KEY (`id_premio`),
  UNIQUE KEY `id_quiniela` (`id_quiniela`),
  CONSTRAINT `premios_ibfk_1` FOREIGN KEY (`id_quiniela`) REFERENCES `quinielas` (`id_quiniela`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 6. Tabla Saldos
CREATE TABLE IF NOT EXISTS `saldos` (
  `id_cliente` INT NOT NULL,
  `saldo_favor` DECIMAL(12,2) NOT NULL DEFAULT '0.00',
  `saldo_deuda` DECIMAL(12,2) NOT NULL DEFAULT '0.00',
  PRIMARY KEY (`id_cliente`),
  CONSTRAINT `saldos_ibfk_1` FOREIGN KEY (`id_cliente`) REFERENCES `clientes` (`id_cliente`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 7. Tabla Movimientos
CREATE TABLE IF NOT EXISTS `movimientos` (
  `id_movimiento` INT NOT NULL AUTO_INCREMENT,
  `id_cliente` INT NOT NULL,
  `id_quiniela` INT DEFAULT NULL,
  `tipo` VARCHAR(30) NOT NULL,
  `monto` DECIMAL(12,2) NOT NULL DEFAULT '0.00',
  `descripcion` TEXT,
  `fecha` TIMESTAMP NULL DEFAULT (NOW()),
  PRIMARY KEY (`id_movimiento`),
  KEY `id_cliente` (`id_cliente`),
  KEY `id_quiniela` (`id_quiniela`),
  CONSTRAINT `movimientos_ibfk_1` FOREIGN KEY (`id_cliente`) REFERENCES `clientes` (`id_cliente`),
  CONSTRAINT `movimientos_ibfk_2` FOREIGN KEY (`id_quiniela`) REFERENCES `quinielas` (`id_quiniela`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 8. Tabla Premio Repartos
CREATE TABLE IF NOT EXISTS `premio_repartos` (
  `id_reparto` INT NOT NULL AUTO_INCREMENT,
  `id_premio` INT NOT NULL,
  `id_quiniela` INT NOT NULL,
  `monto_bruto` DECIMAL(12,2) NOT NULL,
  `porcentaje_retencion` DECIMAL(5,4) NOT NULL,
  `monto_neto` DECIMAL(12,2) NOT NULL,
  `total_recaudado` DECIMAL(12,2) NOT NULL,
  `monto_meta` DECIMAL(12,2) NOT NULL,
  `porcentaje_meta_alcanzado` DECIMAL(7,4) NOT NULL,
  `total_participantes` INT NOT NULL,
  `fecha_reparto` DATETIME NOT NULL,
  PRIMARY KEY (`id_reparto`),
  UNIQUE KEY `uq_premio_repartos_quiniela` (`id_quiniela`),
  KEY `idx_premio_repartos_premio` (`id_premio`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 9. Tabla Premio Reparto Participantes
CREATE TABLE IF NOT EXISTS `premio_reparto_participantes` (
  `id_detalle` INT NOT NULL AUTO_INCREMENT,
  `id_reparto` INT NOT NULL,
  `id_cliente` INT NOT NULL,
  `nombre_cliente` VARCHAR(150) NOT NULL,
  `monto_aporte` DECIMAL(12,2) NOT NULL,
  `porcentaje_aporte` DECIMAL(7,4) NOT NULL,
  `porcentaje_sobre_meta` DECIMAL(7,4) NOT NULL,
  `porcentaje_recompensa` DECIMAL(7,4) NOT NULL,
  `monto_bruto_asignado` DECIMAL(12,2) NOT NULL,
  `monto_retencion` DECIMAL(12,2) NOT NULL,
  `monto_neto_asignado` DECIMAL(12,2) NOT NULL,
  PRIMARY KEY (`id_detalle`),
  UNIQUE KEY `uq_reparto_participantes_cliente` (`id_reparto`, `id_cliente`),
  CONSTRAINT `fk_reparto_participantes_reparto` FOREIGN KEY (`id_reparto`) REFERENCES `premio_repartos` (`id_reparto`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;