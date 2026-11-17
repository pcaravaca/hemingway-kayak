-- 003: queda un trip del seed viejo con coordenadas de Miami (24.5554, -81.78)
-- y un titulo que decia "Crandon Park Sprint". Mismo espiritu que 002: solo
-- toca filas con los valores viejos, es un no-op en una base recien sembrada.
CREATE TEMP TABLE _coord_map (
	old_lat DOUBLE PRECISION,
	old_lon DOUBLE PRECISION,
	new_lat DOUBLE PRECISION,
	new_lon DOUBLE PRECISION
);

INSERT INTO _coord_map (old_lat, old_lon, new_lat, new_lon) VALUES
	(24.5554, -81.7800, 9.8456, -85.1736);

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

UPDATE live_tracking t
   SET latitude = m.new_lat, longitude = m.new_lon
  FROM _coord_map m
 WHERE t.latitude = m.old_lat AND t.longitude = m.old_lon;

UPDATE fishing_spots t
   SET latitude = m.new_lat, longitude = m.new_lon
  FROM _coord_map m
 WHERE t.latitude = m.old_lat AND t.longitude = m.old_lon;

UPDATE trips SET title = 'Playa Hermosa Sprint' WHERE title = 'Crandon Park Sprint';
UPDATE trips SET description = 'Carrera de velocidad en la bahia.' WHERE description = 'Carrera de velocidad en el parque.';

DROP TABLE _coord_map;
