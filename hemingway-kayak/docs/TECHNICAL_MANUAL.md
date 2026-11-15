# Manual Técnico - Hemingway Kayak / Technical Manual - Hemingway Kayak

## Español

### 1. Arquitectura del Sistema

#### 1.1 Componentes Principales
- **Actividad Principal**: `MainActivity.kt`
  - Gestión de UI
  - Inicialización de servicios
  - Manejo de permisos

- **Base de Datos**: `KayakDatabase.kt`
  - Motor: Room
  - Entidades: KayakModel
  - DAOs: KayakDao

- **Servicios**:
  - Clima: OpenWeatherMap API
  - Mapas: Google Maps API
  - Voz: Android Speech Recognition API

#### 1.2 Estructura de Datos
```kotlin
data class KayakModel(
    id: Int,
    name: String,
    type: KayakType,
    length: Float,
    width: Float,
    weight: Float,
    maxCapacity: Float,
    material: KayakMaterial,
    propulsionType: PropulsionType
)
```

### 2. Integración de APIs

#### 2.1 OpenWeatherMap
- Endpoint: `api.openweathermap.org/data/2.5/weather`
- Parámetros requeridos:
  - `lat`: Latitud
  - `lon`: Longitud
  - `appid`: API Key
  - `units`: metric

#### 2.2 Google Maps
- Servicios utilizados:
  - Maps SDK para Android
  - Location Services
  - Geocoding API

### 3. Sistema de Comandos de Voz

#### 3.1 Implementación
```kotlin
class VoiceCommandManager(activity: Activity) {
    // Inicialización
    fun initialize()
    // Escucha
    fun startListening()
    // Detener
    fun stopListening()
}
```

#### 3.2 Comandos Disponibles
- Navegación: START_NAVIGATION, STOP_NAVIGATION
- Clima: CHECK_WEATHER, CHECK_WIND
- Pesca: CHECK_FISHING_CONDITIONS, MARK_FISHING_SPOT
- Emergencia: EMERGENCY_SOS

### 4. Cálculos y Algoritmos

#### 4.1 Condiciones de Pesca
```kotlin
fun analyzeFishingConditions(
    pressure: Float,
    windSpeed: Float,
    moonPhase: Float
): FishingConditions
```

#### 4.2 Rendimiento del Kayak
```kotlin
fun calculateMaxSpeed(
    kayak: KayakModel,
    paddlerWeight: Float
): Float
```

### 5. Requisitos del Sistema

#### 5.1 Desarrollo
- Android Studio Arctic Fox o superior
- Kotlin 1.8.0
- Gradle 7.0+
- JDK 11

#### 5.2 Ejecución
- Android API 24+
- Google Play Services
- Permisos:
  - INTERNET
  - ACCESS_FINE_LOCATION
  - RECORD_AUDIO

### 6. Guía de Deployment

#### 6.1 Configuración de API Keys
1. Crear archivo `local.properties`
2. Agregar:
   ```
   MAPS_API_KEY=your_key_here
   WEATHER_API_KEY=your_key_here
   ```

#### 6.2 Generación de APK
```bash
./gradlew assembleRelease
```

### 7. Sistemas de Seguridad / Safety Systems

#### 7.1 Monitor de Estabilidad / Stability Monitor
```kotlin
class StabilityMonitor(context: Context) {
    // Umbrales / Thresholds
    WARNING_ANGLE = 20f
    DANGER_ANGLE = 35f
    CAPSIZED_ANGLE = 90f
    SUBMERGED_TIMEOUT = 30_000L
}
```

#### 7.2 Rastreador de Ubicación / Location Tracker
```kotlin
class LocationTracker(context: Context) {
    // Configuración / Configuration
    priority = LocationRequest.PRIORITY_HIGH_ACCURACY
    interval = 5000        // 5 segundos / 5 seconds
    displacement = 1f      // 1 metro / 1 meter
}
```

#### 7.3 Sistema de Notificaciones / Notification System
- **Tipos de Contacto / Contact Types**:
  - Personal
  - Servicios de Emergencia / Emergency Services
  - Guardia Costera / Coast Guard
  - Rescate Marítimo / Marine Rescue

- **Métodos de Notificación / Notification Methods**:
  - SMS
  - Email
  - Llamadas / Phone Calls
  - APIs de Emergencia / Emergency APIs

#### 7.4 Formato de Mensajes / Message Format
```
ALERTA DE EMERGENCIA / EMERGENCY ALERT
[Mensaje personalizado / Custom message]

Tipo / Type: [Tipo de emergencia / Emergency type]
Hora / Time: [Timestamp]

Ubicación del vuelco / Capsize location:
  - [Google Maps URL]
  - Precisión / Accuracy: ±[X] metros

Ubicación actual / Current location:
  - [Google Maps URL]

Deriva / Drift:
  - Distancia / Distance: [X] metros
  - Dirección / Direction: [X]° (N/NE/E/SE/S/SW/W/NW)
  - Velocidad / Speed: [X] m/s ([X] km/h)

Información adicional / Additional information:
  - Temperatura del agua / Water temperature: [X]°C
  - Viento / Wind: [X] m/s a / at [dirección/direction]
  - Batería / Battery: [X]%
```

#### 7.5 Flujo de Datos / Data Flow
1. Sensores detectan vuelco / Sensors detect capsize
2. Se crea evento de vuelco / Capsize event is created
3. Inicia seguimiento de ubicación / Location tracking starts
4. Se envían notificaciones / Notifications are sent
5. Actualizaciones periódicas / Periodic updates

#### 7.6 Consideraciones de Seguridad / Safety Considerations
- Redundancia en notificaciones / Notification redundancy
- Múltiples métodos de contacto / Multiple contact methods
- Seguimiento continuo / Continuous tracking
- Mensajes bilingües / Bilingual messages
- Confirmación de envío / Delivery confirmation

### 8. Integración con Sonares / Sonar Integration

#### 8.1 Dispositivos Soportados / Supported Devices

##### Deeper Pro+
- Profundidad máxima / Max depth: 80m
- Frecuencia dual / Dual frequency: 290kHz/90kHz
- GPS integrado / Integrated GPS
- Batería / Battery: 6 horas / 6 hours
- Bluetooth: 100m alcance / 100m range

##### Garmin Striker 4
- CHIRP tradicional / Traditional CHIRP
- GPS de alta sensibilidad / High-sensitivity GPS
- Pantalla / Screen: 3.5" color
- Frecuencia / Frequency: 50/77/200 kHz
- Potencia / Power: 200W RMS

##### iBobber
- Profundidad máxima / Max depth: 40m
- Bluetooth: 30m alcance / 30m range
- Batería / Battery: 10 horas / 10 hours
- LED incorporado / Built-in LED

#### 8.2 Gestión de Datos / Data Management

El `SonarManager` maneja los siguientes tipos de datos / The `SonarManager` handles the following data types:

1. **Datos Básicos / Basic Data**
   - Profundidad / Depth
   - Temperatura / Temperature
   - Posición GPS / GPS Position

2. **Detección de Peces / Fish Detection**
   - Tamaño / Size
   - Profundidad / Depth
   - Movimiento / Movement
   - Confianza / Confidence

3. **Estructuras / Structures**
   - Rocas / Rocks
   - Vegetación / Vegetation
   - Naufragios / Wrecks
   - Artificiales / Artificial

4. **Puntos de Interés / Points of Interest**
   - Puntos de pesca / Fishing spots
   - Estructuras / Structures
   - Peligros / Dangers
   - Waypoints

#### 8.3 Integración / Integration

```kotlin
// Inicializar el gestor / Initialize manager
val sonarManager = SonarManager(context)

// Conectar dispositivo / Connect device
sonarManager.connectDeeper(deviceId)
sonarManager.connectGarmin(deviceId)
sonarManager.connectIBobber(deviceId)

// Iniciar escaneo / Start scanning
sonarManager.startScanning()

// Guardar punto / Save point
sonarManager.savePoint(
    name = "Punto de pesca / Fishing spot",
    description = "Buenos resultados / Good results",
    type = PointType.FISH_SPOT
)
```

#### 8.4 Base de Datos / Database

Los datos se almacenan localmente usando Room / Data is stored locally using Room:

- Historial de escaneos / Scan history
- Puntos guardados / Saved points
- Estadísticas / Statistics
- Configuración / Settings

#### 8.5 Permisos Requeridos / Required Permissions

```xml
<!-- Bluetooth -->
<uses-permission android:name="android.permission.BLUETOOTH" />
<uses-permission android:name="android.permission.BLUETOOTH_ADMIN" />

<!-- Ubicación / Location -->
<uses-permission android:name="android.permission.ACCESS_FINE_LOCATION" />
<uses-permission android:name="android.permission.ACCESS_COARSE_LOCATION" />

<!-- Almacenamiento / Storage -->
<uses-permission android:name="android.permission.WRITE_EXTERNAL_STORAGE" />
<uses-permission android:name="android.permission.READ_EXTERNAL_STORAGE" />

### 9. Sistema de Suscripción / Subscription System

#### 9.1 Google Play Billing
```kotlin
class BillingManager(context: Context) {
    // Configuración / Configuration
    private val billingClient = BillingClient.newBuilder(context)
        .setListener(this)
        .enablePendingPurchases()
        .build()
}
```

#### 9.2 Niveles de Suscripción / Subscription Levels
```kotlin
enum class SubscriptionLevel {
    FREE,           // Básico / Basic
    NAVIGATOR,      // Navegante / Navigator
    RESCUER,        // Rescatista / Rescuer
    PROFESSIONAL    // Profesional / Professional
}
```

#### 9.3 Características por Nivel / Level Features
- **Básico / Basic (Gratis / Free)**
  - 2 contactos / contacts
  - 10s intervalo / interval
  - Funciones básicas / Basic features

- **Navegante / Navigator ($4.99/mes / month)**
  - 5 contactos / contacts
  - 5s intervalo / interval
  - Clima / Weather
  - Modo sin conexión / Offline mode
  - Sin anuncios / No ads

- **Rescatista / Rescuer ($9.99/mes / month)**
  - 10 contactos / contacts
  - 3s intervalo / interval
  - Servicios de emergencia / Emergency services
  - Soporte prioritario / Priority support
  - Exportación / Export

- **Profesional / Professional ($19.99/mes / month)**
  - Contactos ilimitados / Unlimited contacts
  - 1s intervalo / interval
  - Todas las características / All features

#### 9.4 Gestión de Pagos / Payment Management
1. Verificación de compras / Purchase verification
2. Restauración de suscripciones / Subscription restoration
3. Gestión de actualizaciones / Upgrade management
4. Manejo de cancelaciones / Cancellation handling

#### 9.5 UI/UX
- Material Design 3
- Animaciones / Animations
- Temas claro/oscuro / Light/dark themes
- Interfaz bilingüe / Bilingual interface

#### 9.6 Seguridad / Security
- Verificación de recibos / Receipt verification
- Protección anti-fraude / Anti-fraud protection
- Encriptación de datos / Data encryption
- Tokens seguros / Secure tokens

### 10. Sistema de Búsqueda de Lugares / Places Search System

#### 10.1 Tipos de Lugares / Place Types

El sistema soporta los siguientes tipos de lugares / The system supports the following place types:

- Lagos / Lakes
- Ríos / Rivers
- Océanos / Oceans
- Marinas / Marinas
- Muelles / Docks
- Rampas / Boat ramps
- Puntos de pesca / Fishing spots

#### 10.2 Información de Lugares / Place Information

Cada lugar incluye / Each place includes:

1. **Datos Básicos / Basic Data**
   - Nombre / Name
   - Ubicación / Location
   - Tipo / Type
   - Calificación / Rating
   - Fotos / Photos

2. **Condiciones del Agua / Water Conditions**
   - Temperatura / Temperature
   - Profundidad / Depth
   - Velocidad de corriente / Current speed
   - Altura de olas / Wave height
   - Visibilidad / Visibility
   - Horario de mareas / Tide schedule

3. **Instalaciones / Facilities**
   - Estacionamiento / Parking
   - Baños / Restrooms
   - Camping / Camping
   - Alquiler de equipos / Equipment rentals
   - Servicios de comida / Food services
   - Servicios de emergencia / Emergency services

4. **Información Adicional / Additional Information**
   - Restricciones / Restrictions
   - Mejor momento para visitar / Best time to visit
   - Nivel de dificultad / Difficulty level
   - Reseñas de usuarios / User reviews

#### 10.3 Búsqueda de Lugares / Place Search

```kotlin
// Inicializar gestor / Initialize manager
val placesManager = PlacesManager(context)

// Buscar lugares cercanos / Search nearby places
placesManager.searchNearby(
    location = LatLng(latitude, longitude),
    radius = 5000,  // metros / meters
    types = listOf(
        KayakPlaceType.LAKE,
        KayakPlaceType.RIVER,
        KayakPlaceType.FISHING_SPOT
    )
)

// Buscar por ubicación específica / Search by specific location
placesManager.searchByLocation(
    query = "Lago de Chapala, México",
    types = listOf(KayakPlaceType.LAKE)
)

// Obtener detalles / Get details
val placeDetails = placesManager.getPlaceDetails(placeId)

// Guardar favorito / Save favorite
placesManager.saveFavorite(place)
```

#### 10.4 Base de Datos / Database

Los lugares se almacenan usando Room / Places are stored using Room:

- Lugares favoritos / Favorite places
- Historial de visitas / Visit history
- Reseñas de usuario / User reviews
- Condiciones registradas / Recorded conditions

#### 10.5 Permisos Requeridos / Required Permissions

```xml
<!-- Ubicación / Location -->
<uses-permission android:name="android.permission.ACCESS_FINE_LOCATION" />
<uses-permission android:name="android.permission.ACCESS_COARSE_LOCATION" />

<!-- Internet -->
<uses-permission android:name="android.permission.INTERNET" />
<uses-permission android:name="android.permission.ACCESS_NETWORK_STATE" />
```

#### 10.6 APIs Utilizadas / Used APIs

1. **Google Places API**
   - Búsqueda de lugares / Place search
   - Detalles de lugares / Place details
   - Fotos / Photos
   - Reseñas / Reviews

2. **OpenWeatherMap API**
   - Condiciones del agua / Water conditions
   - Temperatura / Temperature
   - Mareas / Tides

3. **NOAA API**
   - Datos marinos / Marine data
   - Corrientes / Currents
   - Pronósticos / Forecasts

---

## English

### 1. System Architecture

#### 1.1 Main Components
- **Main Activity**: `MainActivity.kt`
  - UI Management
  - Service Initialization
  - Permission Handling

- **Database**: `KayakDatabase.kt`
  - Engine: Room
  - Entities: KayakModel
  - DAOs: KayakDao

- **Services**:
  - Weather: OpenWeatherMap API
  - Maps: Google Maps API
  - Voice: Android Speech Recognition API

#### 1.2 Data Structures
```kotlin
data class KayakModel(
    id: Int,
    name: String,
    type: KayakType,
    length: Float,
    width: Float,
    weight: Float,
    maxCapacity: Float,
    material: KayakMaterial,
    propulsionType: PropulsionType
)
```

### 2. API Integration

#### 2.1 OpenWeatherMap
- Endpoint: `api.openweathermap.org/data/2.5/weather`
- Required parameters:
  - `lat`: Latitude
  - `lon`: Longitude
  - `appid`: API Key
  - `units`: metric

#### 2.2 Google Maps
- Services used:
  - Maps SDK for Android
  - Location Services
  - Geocoding API

### 3. Voice Command System

#### 3.1 Implementation
```kotlin
class VoiceCommandManager(activity: Activity) {
    // Initialization
    fun initialize()
    // Listening
    fun startListening()
    // Stop
    fun stopListening()
}
```

#### 3.2 Available Commands
- Navigation: START_NAVIGATION, STOP_NAVIGATION
- Weather: CHECK_WEATHER, CHECK_WIND
- Fishing: CHECK_FISHING_CONDITIONS, MARK_FISHING_SPOT
- Emergency: EMERGENCY_SOS

### 4. Calculations and Algorithms

#### 4.1 Fishing Conditions
```kotlin
fun analyzeFishingConditions(
    pressure: Float,
    windSpeed: Float,
    moonPhase: Float
): FishingConditions
```

#### 4.2 Kayak Performance
```kotlin
fun calculateMaxSpeed(
    kayak: KayakModel,
    paddlerWeight: Float
): Float
```

### 5. System Requirements

#### 5.1 Development
- Android Studio Arctic Fox or higher
- Kotlin 1.8.0
- Gradle 7.0+
- JDK 11

#### 5.2 Runtime
- Android API 24+
- Google Play Services
- Permissions:
  - INTERNET
  - ACCESS_FINE_LOCATION
  - RECORD_AUDIO

### 6. Deployment Guide

#### 6.1 API Keys Configuration
1. Create `local.properties` file
2. Add:
   ```
   MAPS_API_KEY=your_key_here
   WEATHER_API_KEY=your_key_here
   ```

#### 6.2 APK Generation
```bash
./gradlew assembleRelease
```

### 7. Safety Systems

#### 7.1 Stability Monitor
```kotlin
class StabilityMonitor(context: Context) {
    // Thresholds
    WARNING_ANGLE = 20f
    DANGER_ANGLE = 35f
    CAPSIZED_ANGLE = 90f
    SUBMERGED_TIMEOUT = 30_000L
}
```

#### 7.2 Location Tracker
```kotlin
class LocationTracker(context: Context) {
    // Configuration
    priority = LocationRequest.PRIORITY_HIGH_ACCURACY
    interval = 5000        // 5 seconds
    displacement = 1f      // 1 meter
}
```

#### 7.3 Notification System
- **Contact Types**:
  - Personal
  - Emergency Services
  - Coast Guard
  - Marine Rescue

- **Notification Methods**:
  - SMS
  - Email
  - Phone Calls
  - Emergency APIs

#### 7.4 Message Format
```
EMERGENCY ALERT
[Custom message]

Type: [Emergency type]
Time: [Timestamp]

Capsize location:
  - [Google Maps URL]
  - Accuracy: ±[X] meters

Current location:
  - [Google Maps URL]

Drift:
  - Distance: [X] meters
  - Direction: [X]° (N/NE/E/SE/S/SW/W/NW)
  - Speed: [X] m/s ([X] km/h)

Additional information:
  - Water temperature: [X]°C
  - Wind: [X] m/s at [direction]
  - Battery: [X]%
```

#### 7.5 Data Flow
1. Sensors detect capsize
2. Capsize event is created
3. Location tracking starts
4. Notifications are sent
5. Periodic updates

#### 7.6 Safety Considerations
- Notification redundancy
- Multiple contact methods
- Continuous tracking
- Bilingual messages
- Delivery confirmation

### 8. Sonar Integration

#### 8.1 Supported Devices

##### Deeper Pro+
- Max depth: 80m
- Dual frequency: 290kHz/90kHz
- Integrated GPS
- Battery: 6 hours
- Bluetooth: 100m range

##### Garmin Striker 4
- Traditional CHIRP
- High-sensitivity GPS
- Screen: 3.5" color
- Frequency: 50/77/200 kHz
- Power: 200W RMS

##### iBobber
- Max depth: 40m
- Bluetooth: 30m range
- Battery: 10 hours
- Built-in LED

#### 8.2 Data Management

The `SonarManager` handles the following data types:

1. **Basic Data**
   - Depth
   - Temperature
   - GPS Position

2. **Fish Detection**
   - Size
   - Depth
   - Movement
   - Confidence

3. **Structures**
   - Rocks
   - Vegetation
   - Wrecks
   - Artificial

4. **Points of Interest**
   - Fishing spots
   - Structures
   - Dangers
   - Waypoints

#### 8.3 Integration

```kotlin
// Initialize manager
val sonarManager = SonarManager(context)

// Connect device
sonarManager.connectDeeper(deviceId)
sonarManager.connectGarmin(deviceId)
sonarManager.connectIBobber(deviceId)

// Start scanning
sonarManager.startScanning()

// Save point
sonarManager.savePoint(
    name = "Fishing spot",
    description = "Good results",
    type = PointType.FISH_SPOT
)
```

#### 8.4 Database

Data is stored locally using Room:

- Scan history
- Saved points
- Statistics
- Settings

#### 8.5 Required Permissions

```xml
<!-- Bluetooth -->
<uses-permission android:name="android.permission.BLUETOOTH" />
<uses-permission android:name="android.permission.BLUETOOTH_ADMIN" />

<!-- Location -->
<uses-permission android:name="android.permission.ACCESS_FINE_LOCATION" />
<uses-permission android:name="android.permission.ACCESS_COARSE_LOCATION" />

<!-- Storage -->
<uses-permission android:name="android.permission.WRITE_EXTERNAL_STORAGE" />
<uses-permission android:name="android.permission.READ_EXTERNAL_STORAGE" />

### 9. Subscription System

#### 9.1 Google Play Billing
```kotlin
class BillingManager(context: Context) {
    // Configuration
    private val billingClient = BillingClient.newBuilder(context)
        .setListener(this)
        .enablePendingPurchases()
        .build()
}
```

#### 9.2 Subscription Levels
```kotlin
enum class SubscriptionLevel {
    FREE,           // Basic
    NAVIGATOR,      // Navigator
    RESCUER,        // Rescuer
    PROFESSIONAL    // Professional
}
```

#### 9.3 Level Features
- **Basic (Free)**
  - 2 contacts
  - 10s interval
  - Basic features

- **Navigator ($4.99/month)**
  - 5 contacts
  - 5s interval
  - Weather
  - Offline mode
  - No ads

- **Rescuer ($9.99/month)**
  - 10 contacts
  - 3s interval
  - Emergency services
  - Priority support
  - Export

- **Professional ($19.99/month)**
  - Unlimited contacts
  - 1s interval
  - All features

#### 9.4 Payment Management
1. Purchase verification
2. Subscription restoration
3. Upgrade management
4. Cancellation handling

#### 9.5 UI/UX
- Material Design 3
- Animations
- Light/dark themes
- Bilingual interface

#### 9.6 Security
- Receipt verification
- Anti-fraud protection
- Data encryption
- Secure tokens

### 10. Places Search System

#### 10.1 Place Types

The system supports the following place types:

- Lakes
- Rivers
- Oceans
- Marinas
- Docks
- Boat ramps
- Fishing spots

#### 10.2 Place Information

Each place includes:

1. **Basic Data**
   - Name
   - Location
   - Type
   - Rating
   - Photos

2. **Water Conditions**
   - Temperature
   - Depth
   - Current speed
   - Wave height
   - Visibility
   - Tide schedule

3. **Facilities**
   - Parking
   - Restrooms
   - Camping
   - Equipment rentals
   - Food services
   - Emergency services

4. **Additional Information**
   - Restrictions
   - Best time to visit
   - Difficulty level
   - User reviews

#### 10.3 Place Search

```kotlin
// Initialize manager
val placesManager = PlacesManager(context)

// Search nearby places
placesManager.searchNearby(
    location = LatLng(latitude, longitude),
    radius = 5000,  // meters
    types = listOf(
        KayakPlaceType.LAKE,
        KayakPlaceType.RIVER,
        KayakPlaceType.FISHING_SPOT
    )
)

// Search by specific location
placesManager.searchByLocation(
    query = "Lake Chapala, Mexico",
    types = listOf(KayakPlaceType.LAKE)
)

// Get details
val placeDetails = placesManager.getPlaceDetails(placeId)

// Save favorite
placesManager.saveFavorite(place)
```

#### 10.4 Database

Places are stored using Room:

- Favorite places
- Visit history
- User reviews
- Recorded conditions

#### 10.5 Required Permissions

```xml
<!-- Location -->
<uses-permission android:name="android.permission.ACCESS_FINE_LOCATION" />
<uses-permission android:name="android.permission.ACCESS_COARSE_LOCATION" />

<!-- Internet -->
<uses-permission android:name="android.permission.INTERNET" />
<uses-permission android:name="android.permission.ACCESS_NETWORK_STATE" />
```

#### 10.6 Used APIs

1. **Google Places API**
   - Place search
   - Place details
   - Photos
   - Reviews

2. **OpenWeatherMap API**
   - Water conditions
   - Temperature
   - Tides

3. **NOAA API**
   - Marine data
   - Currents
   - Forecasts
