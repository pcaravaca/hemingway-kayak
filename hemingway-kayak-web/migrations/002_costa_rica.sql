-- 002: la semilla de demostracion paso de Miami (Florida) a Costa Rica.
-- Es idempotente: cada UPDATE solo toca las filas que todavia tienen los valores
-- de Miami, asi que en una base recien sembrada (que ya nacio en Costa Rica) es
-- un no-op. No borra nada.

CREATE TEMP TABLE _coord_map (
	old_lat DOUBLE PRECISION,
	old_lon DOUBLE PRECISION,
	new_lat DOUBLE PRECISION,
	new_lon DOUBLE PRECISION
);
INSERT INTO _coord_map (old_lat, old_lon, new_lat, new_lon) VALUES
	(25.7617, -80.1918, 9.9382, -84.8371),
	(25.8900, -80.3600, 10.1200, -85.0200),
	(25.5100, -80.2000, 10.6517, -85.6610),
	(25.6500, -80.2500, 10.7200, -85.7500),
	(25.5300, -80.2200, 10.6900, -85.7000),
	(25.9500, -80.2800, 9.7125, -85.0056),
	(25.9700, -80.3000, 9.7800, -85.0600),
	(25.5250, -80.2050, 10.6600, -85.6700),
	(25.5200, -80.2100, 10.6550, -85.6650),
	(25.5260, -80.2150, 10.6720, -85.6800),
	(25.5180, -80.2080, 10.6480, -85.6580),
	(25.5340, -80.2250, 10.7010, -85.7050),
	(25.5205, -80.2105, 10.6555, -85.6655),
	(25.5210, -80.2110, 10.6560, -85.6660),
	(25.5305, -80.2205, 10.6905, -85.7005),
	(25.5265, -80.2155, 10.6725, -85.6805),
	(25.5185, -80.2085, 10.6485, -85.6585),
	(25.5150, -80.2000, 10.6450, -85.6610),
	(25.5155, -80.2005, 10.6455, -85.6615),
	(25.8505, -80.3505, 9.9805, -84.9805),
	(25.5350, -80.2150, 10.7050, -85.7000),
	(25.9505, -80.2805, 9.7155, -85.0055),
	(25.8500, -80.3500, 9.9800, -84.9800);

UPDATE trips t
   SET start_lat = m.new_lat, start_lon = m.new_lon
  FROM _coord_map m
 WHERE t.start_lat = m.old_lat AND t.start_lon = m.old_lon;
UPDATE trips t
   SET end_lat = m.new_lat, end_lon = m.new_lon
  FROM _coord_map m
 WHERE t.end_lat = m.old_lat AND t.end_lon = m.old_lon;
UPDATE trip_points t
   SET latitude = m.new_lat, longitude = m.new_lon
  FROM _coord_map m
 WHERE t.latitude = m.old_lat AND t.longitude = m.old_lon;
UPDATE tournaments t
   SET start_lat = m.new_lat, start_lon = m.new_lon
  FROM _coord_map m
 WHERE t.start_lat = m.old_lat AND t.start_lon = m.old_lon;
UPDATE tournament_participants t
   SET current_lat = m.new_lat, current_lon = m.new_lon
  FROM _coord_map m
 WHERE t.current_lat = m.old_lat AND t.current_lon = m.old_lon;
UPDATE tournament_catches t
   SET latitude = m.new_lat, longitude = m.new_lon
  FROM _coord_map m
 WHERE t.latitude = m.old_lat AND t.longitude = m.old_lon;
UPDATE tracking_sessions t
   SET start_lat = m.new_lat, start_lon = m.new_lon
  FROM _coord_map m
 WHERE t.start_lat = m.old_lat AND t.start_lon = m.old_lon;
UPDATE tracking_sessions t
   SET end_lat = m.new_lat, end_lon = m.new_lon
  FROM _coord_map m
 WHERE t.end_lat = m.old_lat AND t.end_lon = m.old_lon;
UPDATE live_tracking t
   SET latitude = m.new_lat, longitude = m.new_lon
  FROM _coord_map m
 WHERE t.latitude = m.old_lat AND t.longitude = m.old_lon;
UPDATE fishing_spots t
   SET latitude = m.new_lat, longitude = m.new_lon
  FROM _coord_map m
 WHERE t.latitude = m.old_lat AND t.longitude = m.old_lon;

UPDATE trips SET title = 'Travesia al Golfo de Nicoya' WHERE title = 'Miami Everglades Expedition';
UPDATE trips SET description = 'Travesia por el golfo buscando snook.' WHERE description = 'Travesia por los canales del Everglades buscando snook.';
UPDATE trips SET title = 'Fin de semana en Isla Chira' WHERE title = 'Keys Weekend Adventure';
UPDATE trips SET description = 'Fin de semana en la isla con campamento.' WHERE description = 'Fin de semana en los cayos con campamento.';
UPDATE trips SET title = 'Manana en Playa Hermosa' WHERE title = 'Biscayne Bay Morning';
UPDATE trips SET description = 'Sesion corta de pesca en la bahia de Hermosa.' WHERE description = 'Sesion corta de pesca en la bahia.';
UPDATE trips SET title = 'En vivo: Cruce por el Golfo' WHERE title = 'Live: Florida Bay Crossing';
UPDATE trips SET description = 'Travesia en vivo por el Golfo de Nicoya.' WHERE description = 'Travesia en vivo por Florida Bay.';
UPDATE users SET bio = 'Campeona de kayak racing en Costa Rica.' WHERE bio = 'Campeona de kayak racing en Florida.';
UPDATE tournaments SET name = 'Campeonato Nacional de Kayak Fishing 2026' WHERE name = 'Spring Kingfish Championship 2026';
UPDATE tournaments SET description = 'El clasico campeonato de kingfish de la costa del Golfo de Nicoya.' WHERE description = 'El clasico campeonato de kingfish de la costa de Florida.';
UPDATE tournaments SET name = 'Liga de Carreras del Golfo' WHERE name = 'Kayak Racing League - Miami';
UPDATE tournaments SET description = 'Liga de carreras de kayak por equipos en Playa Hermosa.' WHERE description = 'Liga de carreras de kayak por equipos en Crandon Park.';
UPDATE tournaments SET location_name = 'Marina de Herradura' WHERE location_name = 'Miami Marina';
UPDATE tournaments SET location_name = 'Playa Hermosa' WHERE location_name = 'Crandon Park';
UPDATE tournaments SET location_name = 'Boca del Rio Tempisque' WHERE location_name = 'Biscayne Inlet';
UPDATE tournaments SET status = 'UPCOMING', start_time = '2026-12-05 08:00:00', end_time = '2026-12-05 14:00:00' WHERE name = 'Liga de Carreras del Golfo' AND status = 'LIVE';
UPDATE fishing_spots SET name = 'Los Patos Flats' WHERE name = 'Crandon Flats';

DROP TABLE _coord_map;
