# Hemingway Kayak — Análisis profundo y plan de continuación

**Fecha:** 2026-11-15
**Autor:** Hermes Agent
**Alcance:** `D:\Workspace\Main Repo\hemingway-kayak` (app Android) y `D:\Workspace\Main Repo\hemingway-kayak-web` (backend Go + frontend Vue)
**Método:** inventario real del árbol, lectura de docs (README/AGENTS.md/manuales), verificación de código por 3 subagentes en paralelo + spot-checks propios con `file:line`. Todo lo que dice "verificado" tiene evidencia; lo demás está marcado como hipótesis.

> Nota de nombre: el canal se llama "hemmingway-kayak" (doble m) pero las carpetas y el repo son **hemingway-kayak** (una m). Uso la grafía real del código.

---

## 1. Resumen ejecutivo

- Son **dos proyectos separados y desbalanceados**, no uno: la app Android (`hemingway-kayak`, repo git en GitHub, rama `new-main`) y el backend+web (`hemingway-kayak-web`, **sin git**, 88 MB).
- **Ninguno de los dos funciona de punta a punta hoy.**
  - La app Android **no compila** (verificado): declaración duplicada de `KayakDatabase`, imports faltantes, clases inexistentes, dependencias no declaradas y recursos rotos. Además es un esqueleto: `MainActivity` de 56 líneas sin wiring, **cero ViewModels**, cero tests.
  - El backend **compila** (`go build ./...` y `go vet ./...` exit 0) pero **no puede conectarse a SQL Server** (bug de `Sprintf`: la contraseña no llega al driver) y usa columnas `type` donde el esquema crea `kaytype`/`triptype`/`tournamenttype`/`historytype` → el CRUD de kayaks/trips/tournaments/history está roto en runtime. **No existe middleware de auth**: todos los endpoints son abiertos (IDOR total).
  - El frontend **sí compila** (`vite build` exit 0, 102 módulos, dist commiteado coincide) pero tiene 4 enlaces a vistas inexistentes, auth rota (usa `/api/v1/users/me`, que no existe) y el WebSocket de torneos nunca conecta.
- **El trabajo bueno de Android existe, pero no está en el árbol:** está en `stash@{0}` + `stash@{0}^3` (fragments, ViewModels, nav_graph, 6 tests, `gradle.properties`, compileSdk 36, osmdroid). Recuperarlo es el paso cero.
- **Consolidación: sí es viable y recomendable**, pero no como `.sln` de Visual Studio (Gradle/Go/Vue no los construye un .sln). Lo correcto es un **monorepo** al estilo de `POSMaster`/`loteria-nacional-v2`: una sola raíz con `android/` y `web/`, un solo git, un solo issue tracker.
- Hay **secretos en disco** (API key real de OpenWeatherMap en `local.properties`, token de GitHub embebido en la URL del remote, credenciales SQL y JWT secret en código). Ver sección 8.

---

## 2. Inventario real

| | hemingway-kayak | hemingway-kayak-web |
|---|---|---|
| Stack | Kotlin/Android, Gradle, Room, Retrofit/Moshi, Hilt, Maps, DataStore | Go 1.25 + Fiber v2, MS SQL, JWT, bcrypt, Redis; Vue 3 + Vite + Pinia + Leaflet |
| LOC | 3.509 (.kt, 28 archivos) | 2.675 Go + 7 vistas Vue + 1 CSS |
| Git | Sí, `github.com/pcaravaca/hemingway-kayak`, rama `new-main` (HEAD `3a90e56`) | **No** (sin `.git`) |
| Build | **No compila** | Backend compila; frontend compila |
| Tests | 0 en el árbol (solo artifacts viejos en `app/build`) | 0 |
| CI/Docker | No | No |
| Tamaño | 84 MB `.git` + 78 MB `app` | 88 MB (74 MB `node_modules`) + binario de 14 MB commiteado |

Estructura real (sin `build/`, `.gradle/`, `gradle-7.5/`, `node_modules/`):

```
Main Repo/
├── hemingway-kayak/                 # app Android
│   ├── app/src/main/java/com/hemingway/kayak/
│   │   ├── KayakApplication.kt, MainActivity.kt
│   │   ├── billing/  data/  di/  emergency/  location/
│   │   ├── places/   receivers/  sensors/  sonar/
│   │   ├── ui/map/MapFragment.kt   utils/  voice/
│   ├── app/src/main/res/{layout,values,xml,drawable,mipmap-anydpi-v26}
│   ├── docs/{TECHNICAL_MANUAL.md,USER_MANUAL.md}
│   ├── build.gradle, settings.gradle, gradlew, gradle-7.5/, .gradle/
│   └── README.md
└── hemingway-kayak-web/
    ├── cmd/server/main.go
    ├── internal/{api/server.go, models/models.go, services/*.go, ws/hub.go}
    ├── pkg/database/database.go
    ├── migrations/                    # VACÍO
    ├── web/  (Vite+Vue: src/views, stores, composables, assets; dist/ commiteado)
    ├── AGENTS.md, go.mod, go.sum
    └── hemingway-server               # binario ELF de 14 MB, desactualizado
```

---

## 3. Estado por capa (verificado contra código)

### 3.1 Android — esqueleto + feature classes huérfanas

**Existe de verdad**
- Room real: `@Database` (versión 3), 3 DAOs, TypeConverters, migraciones 1→2 y 2→3 (`data/KayakDatabase.kt:16-133`).
- Retrofit/OpenWeatherMap real: endpoint, modelos Moshi, repositorio con caché y fallback (`data/WeatherApiService.kt:21-33`, `WeatherRepository.kt:30-107`).
- Hilt/DI real pero parcial (`di/DataModule.kt:27-119`): provee DB, DAOs, OkHttp, Moshi, WeatherApiService, WeatherRepository. **No provee `FishingSpotDao`.**
- Maps parcial: `MapFragment` con `OnMapReadyCallback`, marcadores y permisos (`ui/map/MapFragment.kt:103-246`); `MapUtils`.
- Sensores, voz, localización, DataStore, modelos de kayak y facturación: código real (no cableado).
- Estabilidad: umbrales y cálculo pitch/roll (`sensors/StabilityMonitor.kt:34-190`).

**Falta / stub / roto**
- **MVVM no existe**: no hay ningún ViewModel en el árbol; `MapFragment.kt:39` referencia `MapViewModel by viewModels()` (clase inexistente).
- `MainActivity` es esqueleto: `applyTheme()` vacío (`MainActivity.kt:52-55`); los 4 botones del layout no tienen listener; no hay navegación ni map host.
- **Nadie referencia** `MapFragment`, `WeatherRepository`, `VoiceCommandManager`, `EmergencyNotifier`, `SonarManager`, `BillingManager`, `PlacesManager`, `LocationTracker`, `StabilityMonitor`: están huérfanos.
- `SonarManager` es stub (imports de SDKs inexistentes, `SonarDatabase` indefinido, métodos vacíos).
- `PlacesManager` mayormente stub (`isKayakPlace` siempre false → búsqueda siempre vacía).
- `EmergencyNotifier`: email stub con `javax.mail` (no disponible en Android); `notifyCoastGuard/MarineRescue/...` vacíos; SMS real sin permiso `SEND_SMS`.
- `BootReceiver` no-op (solo log, código comentado).
- Sin notificaciones (no hay NotificationChannel, ni servicio en primer plano, ni Worker).
- Sin UI de clima/pesca/ajustes: solo `activity_main.xml` y `fragment_map.xml`.
- **0 tests** en el árbol.
- **No compila** (errores duros): `KayakDatabase` declarada dos veces (`KayakDatabase.kt:22` y `KayakModel.kt:59`); `LocationTracker` sin importar (`EmergencyNotifier.kt:24`); `EmergencyNotifier`/`EmergencyType` sin importar (`StabilityMonitor.kt:54,132`); `MapViewModel`/`SonarDatabase`/`PointDao` inexistentes.
- **Dependencias ausentes**: `com.android.billingclient`, `play-services-places`, `kotlinx-coroutines-play-services`, `javax.mail`.
- **Recursos rotos**: manifest usa `@style/Theme.HemingwayKayak` pero `themes.xml` define `Theme.Hemingway`; `fragment_map.xml` usa `@color/primaryColor` (no existe) y 3 strings inexistentes.
- Falta `gradle.properties` (`android.useAndroidX=true`) y `namespace`.
- **Build inviable**: wrapper pide Gradle 8.0 (`gradle-wrapper.properties:3`) con AGP 7.4.2 + Kotlin 1.8.0 (combinación no soportada); `local.properties` apunta a un SDK inexistente.
- Docs desalineadas: niveles de suscripción (código `ANGLER/CAPTAIN/MASTER` vs docs `NAVIGATOR/RESCUER/PROFESSIONAL`), NOAA API sin implementar, "moon phase" es simulación.

### 3.2 Backend Go — compila, pero no conecta y no autoriza

**Existe de verdad**
- 2.675 LOC; `go build ./...` y `go vet ./...` exit 0 (verificado).
- **Tabla de endpoints de `AGENTS.md:308-387` coincide 1:1 con `internal/api/server.go:77-143`.** No hay endpoints declarados ausentes ni registrados sin declarar.
- MS SQL (go-mssqldb), bcrypt, JWT (emisión), seed, WebSocket (hub), Redis opcional.

**Falta / stub / roto (crítico)**
1. **No hay middleware de auth.** `server.go:42-50` solo registra `recover`, `logger`, `cors`. Ningún handler lee `Authorization`. → IDOR total: cualquiera puede `PUT /users/:id`, borrar kayaks/trips/spots, modificar torneos, unirse con un `userId` ajeno, falsificar posiciones de tracking.
2. **La conexión a SQL está rota.** `database.go:23-26` (y `43-46`): el format string tiene 4 verbos y 5 argumentos, con `password=***` literal. Verificado ejecutando el `Sprintf`: `"password=***...%!(EXTRA string=[REDACTED])"`. La contraseña no llega al driver → `Connect()`/`EnsureDatabase()` no autentican.
3. **Columnas renombradas.** El esquema crea `kayaks.[kaytype]`, `trips.[triptype]`, `tournaments.[tournamenttype]`, `history_entries.[historytype]` (`database.go:129,146,182,250`), pero el código usa `type` (`user_service.go:231,295,337`; `trip_service.go:126,176`; `tournament_service.go:188,241`; `history_service.go:22,26`). `AGENTS.md:516-522` documenta justo la renabra que el código viola → kayaks/trips/tournaments/history fallan en runtime.
4. **`SCOPE_IDENTITY()` sobre tablas `UNIQUEIDENTIFIER` sin IDENTITY** → fallo garantizado en CreateTrip/CreateSpot/CreateTournament/StartTracking/RecordCatch.
5. **Historial siempre vacío**: `GetTripHistory` (10 columnas, 9 destinos de Scan) y `GetTournamentHistory` (9 columnas, 8 destinos) → error por fila y `continue`.
6. **Stubs en torneos**: `GetParticipants`, `GetRanking`, `GetCatches` declaran mapas sin escanear → devuelven `[{},{}]`; `GetRanking` filtra `status='active'` mientras el seed inserta `'ACTIVE'`.
7. **Filtros que no filtran**: `GetTrips` y `GetTournaments` construyen `@userId/@type/@status` pero no los pasan a `QueryContext`.
8. **WebSocket muerto**: `subscribe`/`unsubscribe` son no-op (`hub.go:91-97`), `Broadcast` nunca se invoca, nadie publica → el WS nunca envía `tracking` ni `tournament`. `UpgradeMiddleware` es código muerto.
9. **Logout stub** (no invalida token); `GetProfile` es alias de `GetUser`.
10. **Seed inservible**: el hash bcrypt de los 3 usuarios no verifica la contraseña documentada `password123` (probado). `AGENTS.md` promete "2 active tracking sessions" pero el seed no inserta en `tracking_sessions`.
11. **`migrations/` vacío**: el DDL es inline en `database.go` y **traga todos los errores** (`log.Printf` y `return nil`) → el servidor arranca "OK" aunque el DDL o el seed fallen.
12. `internal/models/models.go` (189 LOC) no se usa: ningún import.
13. `ConnectRedis` devuelve `nil` si el ping falla y `UpdateLocation` hace `s.redis.Publish` sin nil-check → pánico.
14. Sin tests, sin Dockerfile/CI/`.env`, sin graceful shutdown; `go server.Start()` sin capturar error.
15. **Binario commiteado desactualizado**: `hemingway-server` no contiene la contraseña actual; un build del código sí.
16. Endpoints que el frontend consume y no existen: `/api/v1/users/me` y `/api/v1/users/me/stats`.

### 3.3 Frontend Vue — compila, pero auth rota y 4 rutas muertas

**Existe de verdad**
- `vite build` exit 0 (102 módulos, 3,57 s) sobre copia en scratch; los 19 hashes de `dist/` coinciden con el build nuevo → el `dist` commiteado está al día.
- 7 rutas declaradas (`router/index.js:3-39`) y las 7 vistas existen; navbar, design system (`hemingway.css`), Leaflet imperativo en `TrackingView` y `TournamentDetailView`, store Pinia, composable WS con reconexión.

**Falta / stub / roto**
- **4 enlaces a vistas inexistentes**: `/register` (`LoginView.vue:50`), `/profile/edit` (`ProfileView.vue:153`), `/trips/:id` (`HistoryView.vue:103`), `/tournaments/create` (`TournamentsView.vue:5`). Sin ruta catch-all 404.
- **Auth rota en cada recarga**: `auth.js:54` llama `GET /api/v1/users/me` (inexistente) → cae al catch y ejecuta `logout()` (`auth.js:56-58`), borrando el token válido.
- `register` espera `response.data.token` y `.user` (`auth.js:37-38`) pero `Register` solo devuelve `{success,message}`; `Login` devuelve `userId` pero no `user`.
- Sin guards de navegación (`beforeEach`) ni interceptores axios: rutas protegidas accesibles sin sesión y ningún 401 redirige a `/login`.
- **WebSocket de torneos nunca conecta**: `TournamentDetailView.vue:82-85` no llama `connect()` pero sí `subscribeTournament` → no-op; `onUnmounted` no desuscribe.
- Desajustes de campos con el backend: `tournament.type`/`trip.type`/`kayak.type` vs `tournamenttype`/`triptype`/`kaytype`; `trip.distanceKm` vs `distance`; `kayak.lengthCm` vs `length`; `participantsCount` no existe; `toFixed` sin guarda de nulos.
- **No existe `src/components`** (`AGENTS.md:46` lo afirma); en su lugar hay una carpeta vacía mal llamada `{components,views,stores,composables,assets}`.
- `eslint` declarado en `package.json:10` pero no instalado → `npm run lint` falla. Tailwind/PostCSS instalados sin uso. `favicon.ico` referenciado y ausente. `@import` inválido en `App.vue:34`. Typos `--hemmham-*` en CSS.
- API por rutas relativas (funciona solo bajo `vite dev` por el proxy); el backend no sirve estáticos.

---

## 4. El oro: el trabajo real está en el stash

```
stash@{0}: On new-main: auto-stash-pre-pull-1788026400
stash@{0}^3 (commit 4acdef9, "untracked files on new-main"):
  + app/src/main/java/.../ui/{home,weather,fishing,map,settings}/*Fragment.kt + *ViewModel.kt
  + service/FishingConditionsService.kt, sensors/EmergencyType.kt
  + res/layout/{fragment_home,fragment_weather,fragment_fishing,fragment_settings}.xml
  + res/navigation/nav_graph.xml, res/xml/network_security_config.xml, dimens
  + 6 tests (FishingRecommendation, FishingSpot, Resource, WeatherData, EmergencyNotifier, GeoPoint)
  + gradle.properties (android.useAndroidX=true), AGENTS.md
stash@{0} (modificados tracked):
  ~ app/build.gradle (compileSdk 36, namespace, buildConfigField WEATHER_API_KEY, osmdroid 6.1.18, coroutines-play-services)
  ~ MainActivity (+37 líneas), BillingManager (reescrito), WeatherApiService, manifiesto
```

Es decir: el árbol actual (HEAD `3a90e56`, "Migrate to GitHub - clean") es una versión **despojada**; lo bueno quedó guardado por el auto-stash previo a un pull. **Paso cero del plan: recuperar ese stash (o decidir descartarlo conscientemente).**

---

## 5. Consolidación a "solución estilo Visual Studio"

### 5.1 ¿Qué se puede y qué no?
Un `.sln` de Visual Studio **no puede construir** Gradle/Android, Go ni Node/Vue: como solución de compilación no aporta nada. Lo que sí aporta es la **idea**: una raíz única que agrupe ambos proyectos, con un solo git, un solo issue tracker y un solo README/AGENTS. Eso es exactamente el patrón que ya usas en `POSMaster` (`backend/` + `frontend/`) y `loteria-nacional-v2`.

### 5.2 Opciones

| Opción | Descripción | Veredicto |
|---|---|---|
| **A. Monorepo** | Una sola raíz `hemingway-kayak/` con `android/` y `web/`; un solo repo GitHub; se migra el historial Android y se versiona el web | **Recomendada.** Una fuente de verdad, releases coordinados app↔backend, fin a la confusión actual |
| B. Multi-repo + solution file | Mantener 2-3 repos y agregar un `.sln`/`.code-workspace` que los referencie | No resuelve el problema de fondo (sigue el descontrol de fuentes); solo lo maquilla |
| C. Un repo por capa, sin fusionar | Versionar `hemingway-kayak-web` por separado y dejar Android como está | Aceptable si no quieres tocar el historial Android, pero mantiene dos ciclos de vida |

### 5.3 Estructura propuesta (Opción A)

Construida el 2026-11-15 (con los nombres de carpeta originales, para no romper referencias):

```
HemingwayKayak/              # raiz de la solucion (monorepo, un solo git)
├── hemingway-kayak/         # app Kotlin/Gradle (el repo Android, movido)
│   ├── app/                 # modulo Android (compileSdk 36, minSdk 24)
│   ├── gradle/wrapper/      # wrapper 8.2 (el jar SI se versiona)
│   └── docs/                # manuales tecnico y de usuario
├── hemingway-kayak-web/     # backend Go + frontend Vue
│   ├── cmd/server/          # entrypoint Fiber
│   ├── internal/            # api, models, services, ws
│   ├── pkg/                 # database (MS SQL), auth (JWT)
│   ├── migrations/          # DDL
│   └── web/                 # Vue 3 + Vite + Tailwind
├── docs/                    # analisis + plan + decisiones
├── .gitignore               # binario, node_modules, build, .gradle, local.properties, .env
├── README.md                # vision general + como levantar cada parte
└── HemingwayKayak.sln       # agrupador de carpetas (no construye)
```

Migración de git (lo que se hizo de verdad, 2026-11-15):
1. El `.git` original del repo Android se sacó de la carpeta y se preservó en `~/backups/hk-android-dotgit-20261115`; el bundle en `~/backups/hk-android-20261115.bundle`.
2. El monorepo arrancó con **una instantánea** (`git init` + commit inicial `813c094`), sin reescribir historial.
3. El repo Android de GitHub (`pcaravaca/hemingway-kayak`) queda intacto como archivo histórico.
4. Pendiente de decisión: crear el repo GitHub del monorepo y hacer push (requiere aprobación).

**Limpieza obligatoria antes del primer commit:** fuera `hemingway-server` (14 MB), `node_modules/`, `dist/` (se regenera), `.gradle/`, `local.properties`, y el `.git` viejo de 84 MB si se reescribe historial.

---

## 6. Plan de continuación (fases)

> Regla de seguimiento para este proyecto (personal, no DSB): **sin WP**; sí **GT Kanban + Honcho + Vault push**. Todo en inglés para tarjetas Kanban.

**F0 — Gate de aprobación y copia de seguridad** *(requiere tu OK)*
- Backup del árbol actual + `git bundle` del repo Android + export del stash.
- **Decidir el stash**: recuperar `stash@{0}` + `stash@{0}^3` (recomendado) o descartarlo.
- Rotar secretos: API key OWM, token de GitHub embebido en el remote, credenciales SQL, JWT secret.

**F1 — Consolidación del monorepo**
- Crear estructura A, mover `android/` y `web/`, `.gitignore`, README/AGENTS raíz, solution file.
- Limpiar binario/node_modules/dist/.gradle/local.properties.
- Un solo `docker-compose`/scripts por capa.

**F2 — Backend: que conecte y que autorice** *(crítico)*
- Arreglar `Sprintf` de la conexión; mover credenciales a `.env` (no hardcode).
- Alinear columnas `type` → `kaytype/triptype/tournamenttype/historytype` en todos los servicios.
- Reemplazar `SCOPE_IDENTITY()` por `OUTPUT INSERTED.id` / `NEWID()` + retorno.
- **Middleware JWT** + ownership checks; logout con revocación; `/users/me` y `/users/me/stats`.
- Arreglar historial (conteos de columnas), stubs de torneos, filtros sin argumentos, WS (salas + broadcast real).
- Corregir seed (hash que verifique, `tracking_sessions`); que migraciones/seed no traguen errores; versionar DDL en `migrations/`.
- Tests (al menos auth + CRUD) + Dockerfile + CI + graceful shutdown.
- Nil-check de Redis.

**F3 — Frontend: auth y navegación**
- Guards de router + interceptor 401 → `/login`.
- Implementar `RegisterView`, `TripDetailView`, `TournamentCreateView`, `ProfileEditView` (o retirar los enlaces).
- Alinear `auth.js` con lo que el backend devuelve; ruta catch-all 404.
- Conectar el WS de torneos (`connect()` + suscripción/desuscripción).
- Alias de campos (`distanceKm`→`distance`, `lengthCm`→`length`, `type`→`*type`), guards de `toFixed`.
- Instalar eslint, decidir Tailwind vs CSS, favicon, arreglar `@import` y los `--hemmham-*`.
- Definir despliegue (¿el backend sirve `dist/` o hay reverse proxy?).

**F4 — Android: que compile y que cableada**
- Restaurar el stash; alinear AGP/Gradle/Kotlin (p. ej. AGP 8.x + Gradle 8.x + Kotlin 1.9.x) y compileSdk.
- Añadir dependencias que faltan y recursos rotos; `namespace` + `gradle.properties`.
- Decidir **Google Maps vs OSMDroid** (el stash usa osmdroid).
- Cablear `MainActivity` + fragments/ViewModels con nav_graph; inyectar `FishingSpotDao`.
- Implementar stubs (Places, Sonar, EmergencyNotifier, BootReceiver) o recortarlos del alcance.
- Permisos (`SEND_SMS`, `CALL_PHONE`, `BLUETOOTH`, notificaciones en targetSdk 33+).
- Tests (los 6 del stash) + build de CI.

**F5 — Integración app ↔ web**
- Auth compartida (JWT) en la app; envío de tracking al backend; torneos/capturas desde la app.
- Contrato de API documentado (el `AGENTS.md` ya es la base).

**F6 — QA y release**
- Pruebas de humo por capa, build de release (APK + binario + dist), despliegue.

---

## 7. Decisiones que necesito de ti

1. **Stash:** ¿Recupero `stash@{0}` + `stash@{0}^3` como base de la app, o partimos del árbol actual? (Recomiendo recuperar.)
2. **Consolidación:** ¿Opción A (monorepo) confirmada, o prefieres C (versionar web por separado)?
3. **Nombre:** ¿Dejamos `hemingway-kayak` (como está) o corregimos a `hemmingway-kayak`?
4. **Maps:** ¿Google Maps o OSMDroid? (El stash apunta a OSMDroid; el manifiesto actual a Google Maps.)
5. **¿Creo la tarjeta en GT Kanban y la conclusión en Honcho ahora**, o primero revisas el plan? (Proyecto personal: sin WP; solo GT + Honcho + Vault.)
6. **Secretos:** ¿Roto ya la API key de OWM y el token de GitHub del remote?

---

## 8. Secretos y riesgos de seguridad (resumen)

- `hemingway-kayak/local.properties` (gitignored) contiene la **API key real de OpenWeatherMap**; el código usa un placeholder → la app no obtendría clima. Riesgo: key expuesta en el workspace.
- La URL del remote `origin` del repo Android **incrusta un token de GitHub** (requiere rotación).
- Backend: credenciales SQL (`sa`) y JWT secret hardcodeados; `AGENTS.md` documenta la contraseña.
- Backend sin auth: IDOR total (ver 3.2).
- `allowBackup=true` con backup rules que incluyen sharedpref.
- Binario de 14 MB commiteado (fuera en F1).

---

## 9. Evidencia clave (file:line)

- Rutas registradas: `internal/api/server.go:77-143`; middleware solo recover/logger/cors: `server.go:42-50`.
- Conexión rota: `pkg/database/database.go:23-26,43-46`.
- Columnas renombradas: `database.go:129,146,182,250` vs `user_service.go:231`, `trip_service.go:126`, `tournament_service.go:188`, `history_service.go:22`.
- JWT secret: `internal/services/user_service.go:137`; sin validación en ningún handler.
- Seed hash: `database.go:307,310,313`.
- Historial roto: `history_service.go:87-107,138-158`.
- Stubs: `tournament_service.go:321-325,348-352,489-493`; `hub.go:91-97,101-103`.
- Android no compila: `data/KayakDatabase.kt:22` vs `data/KayakModel.kt:59`; `EmergencyNotifier.kt:24`; `StabilityMonitor.kt:54`; `MapFragment.kt:39`.
- Android esqueleto: `MainActivity.kt:52-55`; `activity_main.xml` (4 botones sin listener); 0 ViewModels en el árbol.
- Frontend auth rota: `web/src/stores/auth.js:54`; rutas muertas: `LoginView.vue:50`, `ProfileView.vue:153`, `HistoryView.vue:103`, `TournamentsView.vue:5`.
- Frontend compila: `vite build` exit 0, hashes `dist/` = build nuevo.
- Stash: `git stash list` → `stash@{0}`; `stash@{0}^3` = commit `4acdef9` con 28 archivos y 3.099 inserciones.

---

## 10. Seguimiento (GT Kanban + Honcho + Vault)

Proyecto **personal** → sin WP. Tablero GT Kanban personal:

| Lista | Card |
|---|---|
| ✅ Personal Done | `🛶 Hemingway Kayak · F0: Recover stash, backup, repo hygiene` — id `MUpYTWNuUXdZRVg4T09Kaw` |
| ✅ Personal Done | `🛶 Hemingway Kayak · F1: Consolidate monorepo` — id `Q05vZ3JsVzZJOG5acnpuag` |
| ✅ Personal Done | `🛶 Hemingway Kayak · F2: Backend — connection, columns, JWT` — `dzlrengteXBnT1E2ZkxwTQ` |
| 📋 Personal Backlog | `🛶 · F3: Frontend — auth guards, missing views, WS` — `eUx2c29hNVN5MVpvN0sycw` |
| 📋 Personal Backlog | `🛶 · F4: Android buildable + wired UI` — `Z2gwZDlvOHRjUFN5VEVmVg` |

Honcho: conclusión guardada (estado del proyecto + decisión de monorepo + regla personal sin WP).

## 11. Ejecución F0 (2026-11-15) — hecho y verificado

- **Backup:** `git bundle` completo del repo Android en `~/backups/hk-android-20261115.bundle`.
- **Recuperación del stash:** en un worktree desechable (`Main Repo/_hk_recover`) se aplicó `stash@{0}` + `stash@{0}^3`. Verificado: `compileSdk 36`, `namespace`, `buildConfigField WEATHER_API_KEY`, `osmdroid 6.1.18`, `coroutines-play-services 1.7.3`, AGP 8.2, Gradle wrapper 8.2, 4 fragments + 4 ViewModels, `nav_graph.xml`, 6 tests, `gradle.properties`.
- **Fix de higiene:** `git config core.fileMode false` (el montaje NTFS hace que git vea 70 cambios de modo fantasmas). El `git status` del repo principal quedó **limpio (0 cambios)**.
- **Blob AGENTS.md del stash:** ya describe la solución consolidada (Android + backend + frontend) y usa **OSMDroid** → confirma la dirección de F1/F4.

**Bloqueantes de build detectados (para F4):**
- `compileSdk 36` pero solo está instalado `platforms/android-34`.
- Wrapper pide Gradle 8.2; solo hay cacheado Gradle 8.7.
- No hay `local.properties` (el repo lo gitignorea); el SDK está en `/home/peter/android-sdk`.
- Go en WSL es el de Windows (`go1.22.5 windows/amd64`) y `go.mod` pide 1.25.0.

**Siguiente paso:** F2 — backend (conexión, columnas, JWT).

---

## 12. Ejecución F1 (2026-11-15) — hecho y verificado

**Estructura construida** en `D:\Workspace\Main Repo\HemingwayKayak\`:

| Elemento | Detalle |
|---|---|
| `hemingway-kayak/` | app Android con el stash aplicado (4 fragments, 4 ViewModels, `nav_graph.xml`, 6 tests, `compileSdk 36`, `osmdroid 6.1.18`) |
| `hemingway-kayak-web/` | backend Go + frontend Vue (31 archivos versionados) |
| `docs/` | este análisis + plan |
| `README.md` | estructura, requisitos y comandos por capa |
| `.gitignore` | `node_modules`, `dist`, `build`, `.gradle`, `gradle-7.5/`, `local.properties`, `.env`, binarios |
| `HemingwayKayak.sln` | agrupador de carpetas (no construye) |

**Verificado:**
- Commit local `813c094` — **119 archivos, 808 KB**, `git status` limpio.
- Fuga de secretos: **no** se versionó `local.properties`, `node_modules`, `dist`, `.gradle`, `gradle-7.5/`, binario `hemingway-server` ni `.env`.
- Backups: `~/backups/hk-android-20261115.bundle`, `~/backups/hk-web-20261115.tar.gz`, `~/backups/hk-android-dotgit-20261115`, `~/backups/hemingway-server-20261115`.
- Se eliminó el worktree desechable `_hk_recover`.

**Pendiente / bloqueado:**
- **Root `AGENTS.md` no se pudo escribir** (el guard de archivos de instrucciones de agentes lo bloqueó y la aprobación expiró). El `AGENTS.md` de cada subproyecto sí está versionado.
- **Push del monorepo:** sin repo GitHub todavía; requiere aprobación.
- **Stash:** aplicado al árbol real, pero el stash sigue existiendo en el `.git` preservado (no se hizo `stash drop`).

---

## 13. Ejecución F2 (2026-11-15) — hecho y verificado

Commit local `68bd291` (no pusheado). **17 tests PASS**, `go build ./...` y `go vet ./...` limpios.

**Config y secretos**
- `pkg/config`: todo por variables de entorno, con `.env` cargado sin dependencias (el entorno real gana sobre el archivo) y `Validate()` que falla rápido si faltan `DB_PASSWORD`/`JWT_SECRET`.
- `hemingway-kayak-web/.env` (gitignored, chmod 600) + `.env.example`.
- Verificado en ejecución: con `.env` el servidor pasa la config y falla sólo en la DB; sin `.env` se detiene con el error de configuración.

**Autenticación y autorización**
- `pkg/auth`: `Generate`/`Parse` (HS256) + lista de revocación.
- `internal/middleware`: `JWTAuth`, `RequireOwnership` (ruta), `RequireResourceOwnership` (dueño en BD).
- `internal/api`: todo protegido salvo `health`, `register`, `login` y las lecturas de torneos/trips/spots; añadidos `GET /users/me` y `GET /users/me/stats`; `logout` revoca el token.
- El backend ya no devuelve al usuario en `login` → ahora sí (`user`), que es lo que espera `web/src/stores/auth.js`.

**Base de datos**
- DDL movida a `migrations/001_init.sql`, embebida con `go:embed` y aplicada una sola vez vía `schema_migrations`; **falla en la primera migración con error** (antes tragaba los errores).
- Seed con hash bcrypt real generado en runtime + `tracking_sessions`; antes el hash era un valor de ejemplo.

**Columnas y queries (lo que rompía el runtime)**
- `kayaks.type` → `kaytype`; `trips.type` → `triptype`; `tournaments.type` → `tournamenttype`; `history_entries.type` → `historytype`.
- `SCOPE_IDENTITY()` → `OUTPUT CAST(INSERTED.id AS VARCHAR(36))` (los PK son `UNIQUEIDENTIFIER`, así que `SCOPE_IDENTITY()` devolvía NULL).
- Conteos de columnas corregidos en `GetTripHistory` (10) y `GetTournamentHistory` (9); `GetUserHistory` ahora pasa los parámetros condicionales.
- `GetTrips` y `GetTournaments` construían `@status`/`@type` pero **no pasaban los argumentos** → corregido.
- Stubs implementados: `GetParticipants`, `GetRanking`, `GetCatches` (devolvían `[]` vacío).
- `nil-check` de Redis antes de publicar.

**WebSocket**
- Se eliminó la carrera `delete` bajo `RLock`; el hub ahora maneja salas (`subscribe`/`unsubscribe` con `room`) en lugar de no-ops.

**Infra**
- `Dockerfile` multi-stage (golang:1.25 → distroless) + `.dockerignore`.
- `main.go`: validación de config y **graceful shutdown** (antes ignoraba el error de `Listen`).

**Bloqueante (no verificable aquí):** `10.10.1.21:1433` no es alcanzable desde WSL (`connect: no route to host`). El backend compila y arranca, pero **no se pudo probar contra la base real**. Tampoco hay Redis local (el código lo tolera con nil).

**Nota para F3:** `web/src/stores/auth.js` aún espera `token` en `register` (el backend no lo emite), mapea errores con `data.message` (el backend usa `error`) y su `logout()` local no llama al backend.

---

## 14. Ejecución F3 (2026-11-15) — hecho y verificado

**Auth y navegación**
- `web/src/router/index.js`: 12 rutas con imports dinámicos + guard `beforeEach` (`requiresAuth` → `/login?redirect=…`; si ya autenticado, `/login` → `/`). Rutas nuevas: `/register`, `/profile/edit`, `/tournaments/create`, `/trips/:id`, catch-all 404.
- `web/src/stores/auth.js`: `login`, `register`, `updateProfile`, `logout`, `fetchUser`; JWT en `localStorage` + header `Authorization` global; errores leídos del campo `error` del backend (con fallback a `message`).
- Vistas nuevas: `RegisterView`, `ProfileEditView`, `TournamentCreateView`, `TripDetailView`, `NotFoundView`.

**WebSocket**
- `TournamentDetailView` ahora llama `connect()` en `onMounted` y `unsubscribeTournament(id)` en `onUnmounted` (antes solo suscribía sin conectar).

**Contrato de campos (alineado con lo que el backend emite de verdad)**
- `trip.distanceKm` → `trip.distance` (`HistoryView`).
- `kayak.lengthCm` → `kayak.length` (`ProfileView`).
- `tournament.participantsCount` ya coincidía; `type` ya coincidía (el backend renombra la **columna** DB, no el JSON).
- Guards nulos en `toFixed` (`HomeView`, `TournamentDetailView`, `TrackingView`).
- `HomeView` ya no filtra por `isActive` (campo inexistente; el backend filtra `is_active = 1`).

**Verificado:** `npm run build` exit 0, 107 módulos.

**Pendiente de F3:** `eslint` declarado pero no instalado (`npm run lint` falla) y el modelo de despliegue (¿el backend sirve `dist/` o hay reverse proxy?).

---

## 15. Ejecución F4 (2026-11-15) — build gate pasado

**Precorrente local**
- `local.properties` (gitignored) apuntaba a un SDK de Windows inexistente → corregido a `sdk.dir=/home/peter/android-sdk`.
- AGP 8.2 auto-instaló `platforms;android-36` (licencias ya aceptadas); `build-tools;34.0.0` siguió sirviendo.

**Verificado (evidencia real)**
- `./gradlew :app:assembleDebug` → **BUILD SUCCESSFUL**; APK `app/build/outputs/apk/debug/app-debug.apk` de 9,6 MB, `com.hemingway.kayak`, `minSdk 24`, `targetSdk 36`, `compileSdk 36`.
- `./gradlew :app:testDebugUnitTest` → **26 tests, 0 fallos** (FishingRecommendation 4, FishingSpot 2, Resource 4, WeatherData 2, EmergencyNotifier 9, MapUtils 5).
- Cableado verificado en fuente, no solo compila: `activity_main.xml` aloja `NavHostFragment` + `nav_graph`; `HomeFragment` cablea los 4 botones a las acciones de navegación; `DataModule` provee `FishingSpotDao`.

**Advertencia abierta:** AGP 8.2.0 está probado hasta `compileSdk 34`; el build funciona con 36 pero advierte. Dos caminos: subir AGP a 8.9+ (implica Gradle 8.11+ en el wrapper) o añadir `android.suppressUnsupportedCompileSdk=36`.

**Alcance abierto:** `PlacesManager`, `SonarManager`, `EmergencyNotifier` y `BootReceiver` siguen siendo stubs (compilan, no aportan). Decidir implementarlos o recortarlos del alcance.

**Seguimiento:** GT Kanban personal — F3 → `✅ Personal Done`, F4 → `▶️ Personal In Progress`. Honcho actualizado.
