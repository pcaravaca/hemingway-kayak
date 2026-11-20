# Manual de Usuario - Hemingway Kayak / User Manual - Hemingway Kayak

> **Estado de implementación (2026-11-19).** Lo que este manual describe como
> existente está implementado, salvo lo marcado:
>
> - **Voz (sección 3):** implementada en la pantalla de Tracking. Comandos reales:
>   iniciar/detener tracking, SOS, marcar punto de pesca, distancia y velocidad.
>   Los comandos de clima/pesca de la sección 3 responden "todavía no hace nada".
> - **Seguridad (sección 7):** implementada. Botón SOS en Tracking + contacto de
>   emergencia en Ajustes → Seguridad (guarda teléfono y email). Al dispararla pide
>   el permiso de SMS y abre el 911. El monitor de estabilidad usa los sensores y
>   dispara la alerta si detecta vuelco.
> - **Sonar (sección 8):** **no implementado.** Requiere los SDK propietarios de
>   Garmin/Deeper/ReelSonar y cuentas de desarrollador.
> - **Suscripciones (sección 9):** **no implementado.** Requiere Play Billing y una
>   cuenta de developer de Google Play.
> - **Buscar lugares (sección 10):** cubierto por la pantalla **Spots**, que lista y
>   comparte spots desde el backend. No usa la API de Places de Google.

## Español

### 1. Introducción
Hemingway Kayak es tu compañero perfecto para la pesca en kayak, diseñado específicamente para ser usado con bolsas impermeables y control por voz.

### 2. Primeros Pasos

#### 2.1 Instalación
1. Descarga la aplicación desde Google Play Store
2. Acepta los permisos solicitados:
   - Ubicación (para navegación y clima)
   - Micrófono (para comandos de voz)

#### 2.2 Configuración Inicial
1. Selecciona tu idioma preferido
2. Configura tu kayak:
   - Modelo
   - Dimensiones
   - Peso
   - Tipo de propulsión

### 3. Comandos de Voz

#### 3.1 Activación
- La app escucha constantemente
- No necesitas tocar la pantalla
- Habla claramente y en un tono normal

#### 3.2 Comandos Disponibles

##### Navegación
- "Iniciar navegación" / "Empezar ruta"
- "Detener navegación" / "Parar ruta"

##### Clima
- "Ver clima" / "Condiciones actuales"
- "Ver viento" / "Dirección del viento"

##### Pesca
- "Condiciones de pesca" / "Ver pesca"
- "Marcar punto" / "Guardar ubicación"

##### Emergencia
- "Emergencia" / "Ayuda" / "SOS"

##### Seguimiento
- "Iniciar seguimiento" / "Empezar registro"
- "Detener seguimiento" / "Parar registro"

##### Información
- "Velocidad actual" / "Ver velocidad"
- "Distancia recorrida" / "Ver distancia"

### 4. Funciones Principales

#### 4.1 Navegación
- Mapas con rutas optimizadas
- Puntos de interés para pesca
- Tracking en tiempo real

#### 4.2 Clima
- Temperatura actual
- Velocidad del viento
- Presión atmosférica
- Pronóstico

#### 4.3 Pesca
- Análisis de condiciones
- Registro de capturas
- Puntos favoritos
- Estadísticas

### 5. Consejos y Trucos

#### 5.1 Uso con Bolsa Impermeable
- Mantén el micrófono cerca de la apertura
- Habla en dirección al dispositivo
- Evita el ruido excesivo del viento

#### 5.2 Ahorro de Batería
- Reduce el brillo de la pantalla
- Usa el modo oscuro
- Desactiva funciones no utilizadas

#### 5.3 Emergencias
- Configura contactos de emergencia
- Prueba el comando SOS en tierra
- Mantén el teléfono cargado

### 6. Solución de Problemas

#### 6.1 Comandos de Voz
**Problema**: No reconoce comandos
**Solución**:
1. Verifica el permiso del micrófono
2. Habla más cerca del dispositivo
3. Reduce el ruido ambiental

#### 6.2 GPS
**Problema**: Señal débil
**Solución**:
1. Verifica la ubicación está activada
2. Espera en área abierta
3. Reinicia la aplicación

### 7. Características de Seguridad / Safety Features

#### 7.1 Configuración de Contactos / Contact Setup
1. Abre la sección "Contactos de Emergencia"
2. Toca "Agregar Contacto"
3. Selecciona tipo:
   - Personal (familia, amigos)
   - Servicios de emergencia
4. Configura métodos de notificación:
   - SMS
   - Correo electrónico
   - Llamadas
5. Personaliza mensaje de emergencia (opcional)

#### 7.2 Sistema de Estabilidad / Stability System
- **Alertas de Inclinación**:
  - Advertencia: >20°
  - Peligro: >35°
  - Vuelco: >90°

- **Tonos de Alerta**:
  - Advertencia: tono corto
  - Peligro: tono largo
  - Vuelco: tono de emergencia
  - Hundimiento: tono continuo

#### 7.3 Notificaciones de Emergencia / Emergency Notifications
Los contactos recibirán:
- Ubicación exacta del vuelco
- Ubicación actual
- Distancia y dirección de deriva
- Condiciones ambientales
- Enlaces de Google Maps
- Nivel de batería

#### 7.4 Pruebas de Seguridad / Safety Tests
1. **Prueba Mensual**:
   - Verifica contactos
   - Prueba alertas
   - Actualiza información

2. **Antes de Cada Viaje**:
   - Revisa la batería
   - Confirma señal GPS
   - Verifica permisos

#### 7.5 Recomendaciones / Recommendations
- Mantén el dispositivo en bolsa impermeable
- Asegura el dispositivo al kayak
- Mantén >50% de batería
- Lleva batería externa
- Configura múltiples contactos

#### 7.6 En Caso de Emergencia / In Case of Emergency
1. Mantén la calma
2. Escucha las alertas
3. Espera confirmación de ayuda
4. Conserva batería
5. Mantén el dispositivo visible

### 8. Suscripciones / Subscriptions

#### 8.1 Planes Disponibles / Available Plans

##### Plan Básico / Basic Plan (Gratis / Free)
- 2 contactos de emergencia / emergency contacts
- Funciones básicas / Basic features
- Actualizaciones cada 10s / 10s updates

##### Plan Navegante / Navigator Plan ($4.99/mes / month)
- 5 contactos / contacts
- Alertas de clima / Weather alerts
- Modo sin conexión / Offline mode
- Sin anuncios / No ads
- Actualizaciones cada 5s / 5s updates

##### Plan Rescatista / Rescuer Plan ($9.99/mes / month)
- 10 contactos / contacts
- Servicios de emergencia / Emergency services
- Soporte prioritario / Priority support
- Exportación de datos / Data export
- Actualizaciones cada 3s / 3s updates

##### Plan Profesional / Professional Plan ($19.99/mes / month)
- Contactos ilimitados / Unlimited contacts
- Todas las características / All features
- Actualizaciones cada 1s / 1s updates

#### 8.2 Gestión de Suscripción / Subscription Management

1. **Suscribirse / Subscribe**
   - Abre "Planes y Precios" / Open "Plans & Pricing"
   - Selecciona un plan / Select a plan
   - Sigue el proceso de pago / Follow payment process

2. **Cambiar Plan / Change Plan**
   - Abre "Mi Suscripción" / Open "My Subscription"
   - Selecciona "Cambiar Plan" / Select "Change Plan"
   - Elige nuevo plan / Choose new plan

3. **Cancelar Suscripción / Cancel Subscription**
   - Google Play Store > Suscripciones / Subscriptions
   - Encuentra Hemingway / Find Hemingway
   - Selecciona "Cancelar" / Select "Cancel"

4. **Restaurar Compras / Restore Purchases**
   - Abre "Configuración" / Open "Settings"
   - "Restaurar Compras" / "Restore Purchases"
   - Confirma cuenta / Confirm account

#### 8.3 Preguntas Frecuentes / FAQ

**P: ¿Puedo cambiar de plan? / Can I change plans?**
R: Sí, en cualquier momento / Yes, at any time

**P: ¿Hay reembolsos? / Are there refunds?**
R: Según política de Google Play / As per Google Play policy

**P: ¿Prueba gratuita? / Free trial?**
R: 7 días en planes pagos / 7 days on paid plans

---

## English

### 1. Introduction
Hemingway Kayak is your perfect kayak fishing companion, specifically designed to be used with waterproof bags and voice control.

### 2. Getting Started

#### 2.1 Installation
1. Download the app from Google Play Store
2. Accept requested permissions:
   - Location (for navigation and weather)
   - Microphone (for voice commands)

#### 2.2 Initial Setup
1. Select your preferred language
2. Configure your kayak:
   - Model
   - Dimensions
   - Weight
   - Propulsion type

### 3. Voice Commands

#### 3.1 Activation
- App listens continuously
- No need to touch the screen
- Speak clearly and in a normal tone

#### 3.2 Available Commands

##### Navigation
- "Start navigation" / "Begin route"
- "Stop navigation" / "End route"

##### Weather
- "Check weather" / "Current conditions"
- "Check wind" / "Wind direction"

##### Fishing
- "Fishing conditions" / "Check fishing"
- "Mark spot" / "Save location"

##### Emergency
- "Emergency" / "Help" / "SOS"

##### Tracking
- "Start tracking" / "Begin recording"
- "Stop tracking" / "End recording"

##### Information
- "Current speed" / "Check speed"
- "Distance traveled" / "Check distance"

### 4. Main Features

#### 4.1 Navigation
- Maps with optimized routes
- Fishing points of interest
- Real-time tracking

#### 4.2 Weather
- Current temperature
- Wind speed
- Atmospheric pressure
- Forecast

#### 4.3 Fishing
- Condition analysis
- Catch logging
- Favorite spots
- Statistics

### 5. Tips and Tricks

#### 5.1 Using with Waterproof Bag
- Keep microphone near the opening
- Speak towards the device
- Avoid excessive wind noise

#### 5.2 Battery Saving
- Reduce screen brightness
- Use dark mode
- Disable unused features

#### 5.3 Emergencies
- Set up emergency contacts
- Test SOS command on land
- Keep phone charged

### 6. Troubleshooting

#### 6.1 Voice Commands
**Problem**: Commands not recognized
**Solution**:
1. Verify microphone permission
2. Speak closer to device
3. Reduce ambient noise

#### 6.2 GPS
**Problem**: Weak signal
**Solution**:
1. Verify location is enabled
2. Wait in open area
3. Restart application

### 7. Safety Features

#### 7.1 Contact Setup
1. Open "Emergency Contacts" section
2. Tap "Add Contact"
3. Select type:
   - Personal (family, friends)
   - Emergency services
4. Configure notification methods:
   - SMS
   - Email
   - Calls
5. Personalize emergency message (optional)

#### 7.2 Stability System
- **Tilt Alerts**:
  - Warning: >20°
  - Danger: >35°
  - Capsize: >90°

- **Alert Tones**:
  - Warning: short tone
  - Danger: long tone
  - Capsize: emergency tone
  - Sinking: continuous tone

#### 7.3 Emergency Notifications
Contacts will receive:
- Exact location of capsize
- Current location
- Distance and direction of drift
- Environmental conditions
- Google Maps links
- Battery level

#### 7.4 Safety Tests
1. **Monthly Test**:
   - Verify contacts
   - Test alerts
   - Update information

2. **Before Each Trip**:
   - Check battery
   - Confirm GPS signal
   - Verify permissions

#### 7.5 Recommendations
- Keep device in waterproof bag
- Secure device to kayak
- Maintain >50% battery
- Bring external battery
- Configure multiple contacts

#### 7.6 In Case of Emergency
1. Stay calm
2. Listen to alerts
3. Wait for help confirmation
4. Conserve battery
5. Keep device visible

### 8. Suscripciones / Subscriptions

#### 8.1 Planes Disponibles / Available Plans

##### Plan Básico / Basic Plan (Gratis / Free)
- 2 contactos de emergencia / emergency contacts
- Funciones básicas / Basic features
- Actualizaciones cada 10s / 10s updates

##### Plan Navegante / Navigator Plan ($4.99/mes / month)
- 5 contactos / contacts
- Alertas de clima / Weather alerts
- Modo sin conexión / Offline mode
- Sin anuncios / No ads
- Actualizaciones cada 5s / 5s updates

##### Plan Rescatista / Rescuer Plan ($9.99/mes / month)
- 10 contactos / contacts
- Servicios de emergencia / Emergency services
- Soporte prioritario / Priority support
- Exportación de datos / Data export
- Actualizaciones cada 3s / 3s updates

##### Plan Profesional / Professional Plan ($19.99/mes / month)
- Contactos ilimitados / Unlimited contacts
- Todas las características / All features
- Actualizaciones cada 1s / 1s updates

#### 8.2 Gestión de Suscripción / Subscription Management

1. **Suscribirse / Subscribe**
   - Abre "Planes y Precios" / Open "Plans & Pricing"
   - Selecciona un plan / Select a plan
   - Sigue el proceso de pago / Follow payment process

2. **Cambiar Plan / Change Plan**
   - Abre "Mi Suscripción" / Open "My Subscription"
   - Selecciona "Cambiar Plan" / Select "Change Plan"
   - Elige nuevo plan / Choose new plan

3. **Cancelar Suscripción / Cancel Subscription**
   - Google Play Store > Suscripciones / Subscriptions
   - Encuentra Hemingway / Find Hemingway
   - Selecciona "Cancelar" / Select "Cancel"

4. **Restaurar Compras / Restore Purchases**
   - Abre "Configuración" / Open "Settings"
   - "Restaurar Compras" / "Restore Purchases"
   - Confirma cuenta / Confirm account

#### 8.3 Preguntas Frecuentes / FAQ

**P: ¿Puedo cambiar de plan? / Can I change plans?**
R: Sí, en cualquier momento / Yes, at any time

**P: ¿Hay reembolsos? / Are there refunds?**
R: Según política de Google Play / As per Google Play policy

**P: ¿Prueba gratuita? / Free trial?**
R: 7 días en planes pagos / 7 days on paid plans
