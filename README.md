# Hemingway Kayak

Monorepo del proyecto personal **Hemingway Kayak**: app Android de pesca en kayak +
backend de datos/torneos + panel web.

> Una sola solución, tres surfaces. Antes eran dos carpetas sueltas (`hemingway-kayak`
> y `hemingway-kayak-web`) sin versionar juntas; desde el 2026-11-15 viven aquí.

## Estructura

```
HemingwayKayak/
├── hemingway-kayak/          # App Android (Kotlin, Gradle, MVVM)
│   ├── app/                  # módulo Android (compileSdk 36, minSdk 24)
│   ├── gradle/wrapper/       # wrapper 8.2 (el jar SÍ va versionado)
│   ├── docs/                 # manuales técnico y de usuario
│   └── README.md
├── hemingway-kayak-web/      # Backend + frontend
│   ├── cmd/server/           # entrypoint Go (Fiber)
│   ├── internal/             # api, models, services, ws
│   ├── pkg/                  # database (PostgreSQL), auth (JWT)
│   ├── migrations/           # SQL de migración
│   └── web/                  # frontend Vue 3 + Vite + Tailwind
├── docs/                     # análisis, plan de desarrollo, decisiones
├── HemingwayKayak.sln        # agrupador estilo Visual Studio (solo carpetas)
├── AGENTS.md                 # reglas por capa
└── README.md
```

## Requisitos

| Capa | Herramienta | Versión |
|---|---|---|
| Android | JDK / Android SDK | compileSdk 36, minSdk 24 |
| Android | Gradle | 8.2 (via wrapper) |
| Backend | Go | 1.25.0 (ver `go.mod`) |
| Frontend | Node | 18+ |
| BD | PostgreSQL | 18 (servidor Ubuntu `100.89.167.67:5499`, ver `migrations/`) |

## Cómo ejecutar

### Android
```bash
cd hemingway-kayak
# requiere local.properties con sdk.dir, MAPS_API_KEY y WEATHER_API_KEY (gitignored)
./gradlew :app:assembleDebug
```

### Backend
```bash
cd hemingway-kayak-web
go run ./cmd/server
```

### Frontend
```bash
cd hemingway-kayak-web/web
npm install
npm run dev        # Vite
npm run build      # salida en web/dist/
```

## Estado del proyecto (2026-11-16)

Resumen del análisis completo en [`docs/Hemingway-Kayak-Analisis-2026-11-15.md`](docs/Hemingway-Kayak-Analisis-2026-11-15.md).

- **Android**: construye (`assembleDebug` → APK 9,6 MB), 26 tests verdes, `lintDebug` 0 errores.
- **Backend**: migrado de MS SQL Server a **PostgreSQL 18** (servidor Ubuntu). JWT + ownership,
  migraciones y seed ampliado; 37/37 llamadas probadas end-to-end contra la DB real.
- **Frontend**: auth de punta a punta, rutas protegidas, 5 vistas nuevas, build limpio.
- **Datos de demo**: 6 usuarios, 7 kayaks, 6 viajes con track, 3 torneos, 7 participantes,
  8 capturas, 36 entradas de historial. Login: `peter@hemingway.com` / `HemingwayKayak2026!`

## Historial

El repositorio Android original (`pcaravaca/hemingway-kayak` en GitHub) queda como archivo
histórico. Este monorepo arranca con una instantánea; el `.git` original se preservó en
`~/backups/hk-android-dotgit-20261115` y el bundle en `~/backups/hk-android-20261115.bundle`.
