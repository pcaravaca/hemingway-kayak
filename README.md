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
│   ├── pkg/                  # database (MS SQL), auth (JWT)
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
| BD | MS SQL Server | ver `migrations/` |

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

## Estado del proyecto (2026-11-15)

Resumen del análisis completo en [`docs/Hemingway-Kayak-Analisis-2026-11-15.md`](docs/Hemingway-Kayak-Analisis-2026-11-15.md).

- **Android**: el trabajo real (fragments, ViewModels, `nav_graph`, 6 tests, `compileSdk 36`,
  `osmdroid`) se recuperó del `git stash`; el árbol anterior estaba despojado.
- **Backend**: compila, pero la conexión a MS SQL está rota (`Sprintf` con 4 verbos y 5 args)
  y usa columnas `type` donde el esquema tiene `kaytype`/`triptype`/`tournamenttype`/`historytype`.
  Sin middleware JWT.
- **Frontend**: login degradado (`/auth/me` inexistente → logout en cascada), `localhost:8080`
  hardcodeado, socket de tracking sin reconexión.

## Historial

El repositorio Android original (`pcaravaca/hemingway-kayak` en GitHub) queda como archivo
histórico. Este monorepo arranca con una instantánea; el `.git` original se preservó en
`~/backups/hk-android-dotgit-20261115` y el bundle en `~/backups/hk-android-20261115.bundle`.
