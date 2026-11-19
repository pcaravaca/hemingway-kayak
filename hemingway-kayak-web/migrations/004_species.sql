-- 004_species.sql — DB de peces de Costa Rica (rio / estuario / mar).
-- Fuente: Vault Pesca/cr-fish-field-guide -> references/cr_fish_db.json (18 especies).
-- Idempotente: CREATE IF NOT EXISTS + ON CONFLICT (scientific) DO NOTHING.

CREATE TABLE IF NOT EXISTS species (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	common_cr TEXT NOT NULL,
	scientific TEXT NOT NULL,
	family TEXT,
	taxonomic_order TEXT,
	habitat TEXT NOT NULL DEFAULT 'mar',
	size_cm TEXT,
	gear TEXT,
	edible BOOLEAN DEFAULT FALSE,
	invasive BOOLEAN DEFAULT FALSE,
	lookalikes JSONB DEFAULT '[]'::jsonb,
	distinguishing TEXT,
	incopesca TEXT DEFAULT 'pendiente',
	created_at TIMESTAMPTZ DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS species_scientific_key ON species (scientific);
CREATE INDEX IF NOT EXISTS species_habitat_key ON species (habitat);

INSERT INTO species (common_cr, scientific, family, taxonomic_order, habitat, size_cm, gear, edible, invasive, lookalikes, distinguishing, incopesca)
VALUES
	('Guapote / Rainbow bass', 'Parachromis dovii', 'Cichlidae', 'Perciformes', 'rio', 'hasta 70', 'UL/BFS: jig, crank 2-5g, soft plastic', TRUE, FALSE, '["tilapia", "mojarra"]'::jsonb, 'boca grande, cuerpo alargado, barras verticales oscuras de joven, agresivo', 'pendiente'),
	('Machaca', 'Brycon guatemalensis', 'Characidae', 'Characiformes', 'rio', '30-60', 'spinner, small jig', TRUE, FALSE, '["robalo juvenil", "sabalo"]'::jsonb, 'escamas grandes, boca terminal, salta al pegarle', 'pendiente'),
	('Tilapia', 'Oreochromis niloticus', 'Cichlidae', 'Perciformes', 'rio', '20-40', 'artificial pequeno, masa, gusano', TRUE, TRUE, '["mojarra", "guapote chico"]'::jsonb, 'cuerpo redondo, boca chica, dorsal espinosa larga', 'pendiente'),
	('Mojarra (dulce)', 'Caquetaia/Amphilophus spp.', 'Cichlidae', 'Perciformes', 'rio', '10-25', 'artificial pequeno, masa', TRUE, FALSE, '["tilapia", "guapote"]'::jsonb, 'ciclida pequena, patrones irregulares', 'pendiente'),
	('Robalo de rio / Robalillo', 'Centropomus spp. (juvenil)', 'Centropomidae', 'Carangiformes', 'estuario', '15-40', 'alacran de mar (erizo), jig', TRUE, FALSE, '["machaca", "barracuda"]'::jsonb, 'linea lateral curva hacia abajo tras cabeza, boca oblicua', 'pendiente'),
	('Bagre / Barbudo', 'Rhamdia spp.', 'Heptapteridae', 'Siluriformes', 'rio', '20-50', 'carnada de fondo', TRUE, FALSE, '[]'::jsonb, 'bigotes, sin escamas, piel viscosa', 'pendiente'),
	('Robalo / Snook', 'Centropomus undecimalis', 'Centropomidae', 'Carangiformes', 'estuario', '30-100+', 'alacran de mar (erizo), jig', TRUE, FALSE, '["machaca", "barracuda", "sabalo"]'::jsonb, 'linea lateral CURVA hacia abajo tras cabeza, boca grande oblicua, aleta anal larga', 'pendiente'),
	('Pargo', 'Lutjanus guttatus', 'Lutjanidae', 'Perciformes', 'estuario', '25-80', 'spinnerbait rojo ~1/4oz, carnada natural', TRUE, FALSE, '["mojarra de mar", "corvina", "jurel"]'::jsonb, 'colmillos caninos, hocico puntiagudo, color rojo/amarillo', 'pendiente'),
	('Corvina', 'Cynoscion albus', 'Sciaenidae', 'Perciformes', 'estuario', '40-120', 'carnada viva, jig', TRUE, FALSE, '["pargo", "macarela"]'::jsonb, 'barbillas en la barbilla, escamas pequenas, sin colmillos', 'pendiente'),
	('Macarela / Sierra', 'Scomberomorus sierra', 'Scombridae', 'Perciformes', 'mar', '30-80', 'cucharilla, pez artificial rapido', TRUE, FALSE, '["corvina"]'::jsonb, 'cuerpo fusiforme, manchas oscuras redondas, cola ahorquillada', 'pendiente'),
	('Pez gallo / Roosterfish', 'Caranx caballus', 'Carangidae', 'Carangiformes', 'mar', '50-150', 'carnada viva, jig pesado', TRUE, FALSE, '[]'::jsonb, '7 espinas largas en la frente (cresta de gallo)', 'pendiente'),
	('Barracuda', 'Sphyraena ensis', 'Sphyraenidae', 'Perciformes', 'estuario', '40-100', 'carnada, jig', TRUE, FALSE, '["robalo", "sabalo"]'::jsonb, 'cuerpo muy alargado, dientes grandes, sin curva de linea lateral marcada', 'pendiente'),
	('Lisa / Mule', 'Mugil cephalus', 'Mugilidae', 'Mugiliformes', 'estuario', '30-60', 'carnada, masa', TRUE, FALSE, '[]'::jsonb, 'plateado, ojo con banda grasa, boca chica', 'pendiente'),
	('Cabrilla / Grouper', 'Epinephelus spp.', 'Serranidae', 'Perciformes', 'mar', '30-100', 'carnada, jig pesado', TRUE, FALSE, '[]'::jsonb, 'cuerpo robusto, boca grande, cola redondeada', 'pendiente'),
	('Jurel', 'Caranx sexfasciatus', 'Carangidae', 'Carangiformes', 'mar', '30-70', 'carnada, jig', TRUE, FALSE, '["pargo", "corvina"]'::jsonb, 'cuerpo comprimido, escudo lateral dentado', 'pendiente'),
	('Dorado / Mahi-mahi', 'Coryphaena hippurus', 'Coryphaenidae', 'Perciformes', 'mar', '60-150', 'carnada, pez artificial', TRUE, FALSE, '[]'::jsonb, 'frente vertical (macho), verde/oro/azul, cola ahorquillada', 'pendiente'),
	('Pez vela / Sailfish', 'Istiophorus platypterus', 'Istiophoridae', 'Istiopteriformes', 'mar', '150-300', 'carnada grande (kayak offshore)', TRUE, FALSE, '[]'::jsonb, 'pico largo, aleta dorsal tipo vela', 'pendiente'),
	('Tiburon', 'Rhizoprionodon/Carcharhinus spp.', 'Carcharhinidae', 'Carcharhiniformes', 'mar', '80-200', 'carnada (kayak)', FALSE, FALSE, '[]'::jsonb, 'branquias en hendiduras, cola asimetrica', 'pendiente')
ON CONFLICT (scientific) DO NOTHING;
