-- 006_capturas_robalo.sql
-- 005_incopesca_y_capturas.sql se aplicó a hemingway_dev antes de que se le agregara
-- la equivalencia 'Robalo' -> 'Robalo / Snook'. Las bases que ya corrieron 005 necesitan
-- ese renglón; las nuevas lo traen de 005. Idempotente (queda en 0 filas).
UPDATE tournament_catches SET species = 'Robalo / Snook' WHERE species = 'Robalo';
