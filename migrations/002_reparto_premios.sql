-- Reparto del premio entre los participantes de una quiniela ganadora.
--
-- El reparto se congela en el momento en que se genera. Estas tablas guardan
-- los importes ya calculados para que consultar el reparto devuelva siempre
-- las mismas cifras, aunque despues se registren mas aportes o se actualice
-- el premio. No recalcular nada: esto es el comprobante de lo que se agreed a
-- pagar.
--
-- Idempotente: se puede volver a ejecutar sin romper nada.

CREATE TABLE IF NOT EXISTS premio_repartos (
    id_reparto                INT AUTO_INCREMENT PRIMARY KEY,
    id_premio                 INT          NOT NULL,
    id_quiniela               INT          NOT NULL,

    -- Copia del premio tal como estaba al repartir. Se duplica a proposito:
    -- la tabla premios puede actualizarse despues y el comprobante no.
    monto_bruto               DECIMAL(12, 2) NOT NULL,
    porcentaje_retencion      DECIMAL(5, 4)  NOT NULL,
    monto_neto                DECIMAL(12, 2) NOT NULL,

    -- Fotografia de la quiniela en el momento del reparto.
    total_recaudado           DECIMAL(12, 2) NOT NULL,
    monto_meta                DECIMAL(12, 2) NOT NULL,
    porcentaje_meta_alcanzado DECIMAL(7, 4)  NOT NULL,
    total_participantes       INT          NOT NULL,

    fecha_reparto             DATETIME     NOT NULL,

    -- Una sola fila por quiniela: es la garantia de que el premio no se
    -- reparte (y por lo tanto no se paga) dos veces.
    UNIQUE KEY uq_premio_repartos_quiniela (id_quiniela),
    KEY idx_premio_repartos_premio (id_premio)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE IF NOT EXISTS premio_reparto_participantes (
    id_detalle                INT AUTO_INCREMENT PRIMARY KEY,
    id_reparto                INT          NOT NULL,

    id_cliente                INT          NOT NULL,

    -- El nombre tambien se copia: si el cliente cambia de nombre despues, el
    -- comprobante debe seguir mostrando como se llamaba al repartir.
    nombre_cliente            VARCHAR(150) NOT NULL,
    monto_aporte              DECIMAL(12, 2) NOT NULL,
    porcentaje_aporte         DECIMAL(7, 4)  NOT NULL,
    porcentaje_sobre_meta     DECIMAL(7, 4)  NOT NULL,
    porcentaje_recompensa     DECIMAL(7, 4)  NOT NULL,

    -- Lo que le toca a este cliente.
    monto_bruto_asignado      DECIMAL(12, 2) NOT NULL,
    monto_retencion           DECIMAL(12, 2) NOT NULL,
    monto_neto_asignado       DECIMAL(12, 2) NOT NULL,

    CONSTRAINT fk_reparto_participantes_reparto
        FOREIGN KEY (id_reparto) REFERENCES premio_repartos (id_reparto)
        ON DELETE CASCADE,

    -- Un cliente no puede aparecer dos veces en el mismo reparto.
    UNIQUE KEY uq_reparto_participantes_cliente (id_reparto, id_cliente)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;
