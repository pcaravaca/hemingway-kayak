-- Normaliza el estado de las salidas: la semilla usa mayúsculas (PLANNED,
-- LIVE, COMPLETED, CANCELLED) pero CreateTrip dejaba el default de la columna,
-- que estaba en minúscula. Se deja todo en mayúscula y se cambia el default.

UPDATE trips SET status = UPPER(status) WHERE status <> UPPER(status);

ALTER TABLE trips ALTER COLUMN status SET DEFAULT 'PLANNED';
