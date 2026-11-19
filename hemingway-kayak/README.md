# Hemingway Kayak

## English Description

Hemingway Kayak is an Android application designed to enhance the kayak fishing experience by providing real-time weather conditions and fishing recommendations. The app combines weather data, barometric pressure trends, and moon phases to help anglers make informed decisions about their fishing trips.

### Features

- **Real-time Weather Data**
  - Temperature
  - Barometric pressure
  - Humidity
  - Wind speed and direction
  - Cloud cover
  - Precipitation

- **Fishing Conditions Analysis**
  - Pressure trend monitoring (rising, falling, stable)
  - Moon phase calculation
  - Fishing favorability assessment
  - Detailed recommendations based on current conditions

- **Location Features**
  - GPS tracking
  - Mapa con osmdroid (spots de pesca locales)
  - Current location weather

### Technical Details

- **Platform**: Android (minimum SDK 24)
- **Languages**: Kotlin
- **Architecture**: MVVM
- **Key Dependencies**:
  - osmdroid + Fused Location Provider
  - OpenWeatherMap API
  - Retrofit for API communication
  - Room for local storage
  - Kotlin Coroutines
  - LiveData & ViewModel

### Visual Design

#### Visual Identity
- **Logo**: Minimalist design combining a kayak with technological elements
  - Circular radar with navigation lines
  - Stylized kayak at the center
  - GPS/connectivity signals
  - Stylized water waves

#### Color Palette
- **Neon Green** (`#4DFFA9`): Primary and interactive elements
- **Turquoise** (`#40E0D0`): Secondary elements and details
- **Dark Background** (`#1A1B1E`): Main background color

#### Visual Resources
- **Adaptive Icon**: Support for modern Android with separate background and foreground
- **Splash Screen**: Animated logo with radar effect
- **Different Densities**: Support for all Android screen densities (mdpi to xxxhdpi)

### Setup Requirements

1. Android Studio Arctic Fox or newer
2. OpenWeatherMap API key
3. Android device or emulator (API 24+)

### Possible Improvements

- **Offline Functionality**: Data storage for areas without connection
- **Personalization**: Allow users to configure favorite fishing parameters
- **Community**: Add social elements to share locations and experiences
- **Notifications**: Alerts about optimal fishing conditions
- **Platform Expansion**: iOS version to reach more users

---

## Descripción en Español

Hemingway Kayak es una aplicación Android diseñada para mejorar la experiencia de pesca en kayak proporcionando condiciones climáticas en tiempo real y recomendaciones de pesca. La aplicación combina datos meteorológicos, tendencias de presión barométrica y fases lunares para ayudar a los pescadores a tomar decisiones informadas sobre sus salidas de pesca.

### Características

- **Datos Meteorológicos en Tiempo Real**
  - Temperatura
  - Presión barométrica
  - Humedad
  - Velocidad y dirección del viento
  - Cobertura de nubes
  - Precipitación

- **Análisis de Condiciones de Pesca**
  - Monitoreo de tendencias de presión (subiendo, bajando, estable)
  - Cálculo de fase lunar
  - Evaluación de favorabilidad para pesca
  - Recomendaciones detalladas basadas en condiciones actuales

- **Funciones de Ubicación**
  - Seguimiento GPS
  - Mapa con osmdroid (spots de pesca locales)
  - Clima en ubicación actual

### Detalles Técnicos

- **Plataforma**: Android (SDK mínimo 24)
- **Lenguajes**: Kotlin
- **Arquitectura**: MVVM
- **Dependencias Principales**:
  - osmdroid + Fused Location Provider
  - API de OpenWeatherMap
  - Retrofit para comunicación con API
  - Room para almacenamiento local
  - Corrutinas de Kotlin
  - LiveData y ViewModel

### Diseño Visual

#### Identidad Visual
- **Logo**: Diseño minimalista que combina un kayak con elementos tecnológicos
  - Radar circular con líneas de navegación
  - Kayak estilizado en el centro
  - Señales de GPS/conectividad
  - Ondas de agua estilizadas

#### Paleta de Colores
- **Verde Neón** (`#4DFFA9`): Elementos principales e interactivos
- **Turquesa** (`#40E0D0`): Elementos secundarios y detalles
- **Fondo Oscuro** (`#1A1B1E`): Color de fondo principal

#### Recursos Visuales
- **Icono Adaptativo**: Soporte para Android moderno con fondo y primer plano separados
- **Pantalla de Inicio**: Logo animado con efecto de radar
- **Diferentes Densidades**: Soporte para todas las densidades de pantalla Android (mdpi hasta xxxhdpi)

### Requisitos de Configuración

1. Android Studio Arctic Fox o más reciente
2. Clave de API de OpenWeatherMap
3. Dispositivo Android o emulador (API 24+)

### Posibles Mejoras

- **Funcionalidad offline**: Almacenamiento de datos para áreas sin conexión
- **Personalización**: Permitir a usuarios configurar parámetros de pesca favoritos
- **Comunidad**: Agregar elementos sociales para compartir lugares y experiencias
- **Notificaciones**: Alertas sobre condiciones óptimas de pesca
- **Expansión de plataforma**: Versión iOS para alcanzar más usuarios
