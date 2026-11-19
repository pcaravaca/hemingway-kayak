-- 005_incopesca_y_capturas.sql
-- 1) INCOPESCA: talla legal de primera captura (TLPC) por especie.
--    Fuente: INCOPESCA, Acuerdo AJDIP/026-2018 (tallas legales de primera captura,
--    respondiendo a la talla de primera madurez sexual) y Plan de Manejo del
--    Recurso Pesquero del RNVS Barra del Colorado (SINAC-INCOPESCA, 2021).
--    Idempotente: cada UPDATE va limitado por el nombre científico.
-- 2) Alta de Tarpon (necesaria para que las capturas del seed resuelvan taxonimia).
-- 3) Las capturas de demostracion pasan de nombres en ingles a los nombres comunes
--    de la DB CR. Tambien idempotente.

-- ---------- 1) TLPC / notas INCOPESCA ----------

-- AJDIP/026-2018, Cuadro 1 (peces): valores oficiales.
UPDATE species SET incopesca = 'TLPC 30 cm LT (tolerancia 10%) - AJDIP/026-2018'
 WHERE scientific = 'Lutjanus guttatus';

UPDATE species SET incopesca = 'TLPC 50 cm LT (tolerancia 10%) - AJDIP/026-2018 (especie "Corvina Reina")'
 WHERE scientific = 'Cynoscion albus';

UPDATE species SET incopesca = 'TLPC 50 cm LT (tolerancia 10%) - AJDIP/026-2018'
 WHERE scientific = 'Scomberomorus sierra';

UPDATE species SET incopesca = 'TLPC 80 cm LT (tolerancia 10%) - AJDIP/026-2018'
 WHERE scientific = 'Coryphaena hippurus';

UPDATE species SET incopesca = 'TLPC 38 cm LT (tolerancia 10%) - AJDIP/026-2018'
 WHERE scientific = 'Sphyraena ensis';

-- Tiburones: AJDIP/026-2018, Cuadro 4 (tiburones).
UPDATE species SET incopesca = 'AJDIP/026-2018 por especie: gris 100 cm LT (tol. 20%), martillo 85 cm, azul 46 cm, mamón 43 cm. Aletas adheridas obligatorio.'
 WHERE scientific = 'Rhizoprionodon/Carcharhinus spp.';

-- Agua dulce: el Plan de Manejo del RNVS Barra del Colorado (2021) fija talla para
-- guapote PINTO (Parachromis managuensis), no para P. dovii; se deja constancia.
UPDATE species SET incopesca = 'Sin TLPC nacional. Plan Barra del Colorado: guapote pinto (P. managuensis) 25-40 cm LT.'
 WHERE scientific = 'Parachromis dovii';

-- Robalo: AJDIP/026-2018 sólo regula Centropomus unionensis (gualaje), no C. undecimalis.
UPDATE species SET incopesca = 'Sin TLPC para C. undecimalis; AJDIP/026-2018 fija 28 cm LT para C. unionensis (gualaje).'
 WHERE scientific = 'Centropomus undecimalis';

UPDATE species SET incopesca = 'Sin TLPC para C. undecimalis; AJDIP/026-2018 fija 28 cm LT para C. unionensis (gualaje).'
 WHERE scientific = 'Centropomus spp. (juvenil)';

-- Sin TLPC oficial publicado por INCOPESCA.
UPDATE species SET incopesca = 'Sin TLPC oficial (INCOPESCA no la incluye en AJDIP/026-2018).'
 WHERE scientific IN ('Brycon guatemalensis', 'Caquetaia/Amphilophus spp.', 'Rhamdia spp.',
                      'Caranx caballus', 'Mugil cephalus', 'Epinephelus spp.', 'Caranx sexfasciatus');

UPDATE species SET incopesca = 'Sin TLPC; especie invasora, sin veda.'
 WHERE scientific = 'Oreochromis niloticus';

UPDATE species SET incopesca = 'Sin TLPC; pesca deportiva. Anzuelo circular obligatorio si se usa carnada natural.'
 WHERE scientific = 'Istiophorus platypterus';

-- ---------- 2) Tarpon / Atlantic Tarpon (la pedía el seed "Fall Tarpon Shootout") ----------

INSERT INTO species (common_cr, scientific, family, taxonomic_order, habitat, size_cm, gear, edible, invasive, lookalikes, distinguishing, incopesca)
VALUES ('Tarpon / Atlantic Tarpon', 'Megalops atlanticus', 'Megalopidae', 'Elopiformes', 'mar',
        '60-200', 'carnada viva, jig pesado', TRUE, FALSE, '[]'::jsonb,
        'plateado, escama grande, boca con mandíbula superior larga; salta mucho', 'Sin TLPC oficial.')
ON CONFLICT (scientific) DO NOTHING;

-- ---------- 3) Capturas del seed a nombres de la DB CR ----------

UPDATE tournament_catches SET species = 'Macarela / Sierra' WHERE species = 'King Mackerel';
UPDATE tournament_catches SET species = 'Macarela / Sierra' WHERE species = 'Spanish Mackerel';
UPDATE tournament_catches SET species = 'Tarpon / Atlantic Tarpon' WHERE species = 'Tarpon';
UPDATE tournament_catches SET species = 'Pargo' WHERE species = 'Snapper';
UPDATE tournament_catches SET species = 'Robalo / Snook' WHERE species = 'Robalo';
UPDATE tournament_catches SET species = 'Robalo / Snook' WHERE species = 'Snook';
UPDATE tournament_catches SET species = 'Jurel' WHERE species = 'Jack Crevalle';
