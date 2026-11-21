-- El puntaje de un torneo pasa a contar SOLO las capturas verificadas.
--
-- Antes RecordCatch sumaba el peso al cargar la captura, sin importar `verified`.
-- Ahora la captura entra con verified=false y recién al verificarse (por el
-- organizador) suma score y catches_count; al desverificar, se resta.
--
-- Es idempotente: recomputa desde el estado actual de tournament_catches, así que
-- es un no-op si se vuelve a correr.

UPDATE tournament_participants p
SET score = COALESCE(c.s, 0),
    catches_count = COALESCE(c.n, 0)
FROM (
    SELECT participant_id, SUM(weight) AS s, COUNT(*) AS n
    FROM tournament_catches
    WHERE verified = TRUE
    GROUP BY participant_id
) c
WHERE p.id = c.participant_id;

UPDATE tournament_participants p
SET score = 0,
    catches_count = 0
WHERE NOT EXISTS (
    SELECT 1 FROM tournament_catches c
    WHERE c.participant_id = p.id AND c.verified = TRUE
);
