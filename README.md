# R2k Service Family Manager — Guía de Compilación, Pruebas y Ejecución

## Tabla de Contenidos

1. [Prerrequisitos](#1-prerrequisitos)
2. [Configuración del Repositorio y Entorno](#2-configuración-del-repositorio-y-entorno)
3. [Estructura del Proyecto](#3-estructura-del-proyecto)
4. [Proceso de Compilación](#4-proceso-de-compilación)
5. [Ejecución del Instalador](#5-ejecución-del-instalador)
6. [Escenarios de Prueba](#6-escenarios-de-prueba)
7. [Solución de Problemas](#7-solución-de-problemas)

---

## 1. Prerrequisitos

### Software Requerido

| Herramienta       | Versión | Propósito                         | Instalación                                                     |
|-------------------|---------|-----------------------------------|-----------------------------------------------------------------|
| **Windows 10/11** | Cualq.  | SO Objetivo (servicios usan `sc`) | —                                                               |
| **Go**            | 1.25+   | Compilador                        | [go.dev/dl](https://go.dev/dl/)                                 |
| **Task**          | 3.x     | Orquestación de compilación       | [taskfile.dev/installation](https://taskfile.dev/installation/) |
| **Git**           | Cualq.  | Clonar repositorios               | [git-scm.com](https://git-scm.com/)                             |
| **PowerShell**    | 5.1+    | Scripts de compilación (nativo)   | Integrado                                                       |

### Verificar Instalación

Abra **PowerShell como Administrador** y ejecute:

```powershell
go version          # Esperado: go1.25+
task --version      # Esperado: Task version 3.x
git --version       # Esperado: git version 2.x+

```

---

## 2. Configuración del Repositorio y Entorno

Los tres repositorios deben ser **directorios hermanos** (al mismo nivel) — el Taskfile hace referencia a
`../scale-daemon` y `../ticket-daemon` en relación con `poster-tuis`.

```powershell
# Crear directorio de trabajo
mkdir C:\dev\r2k
cd C:\dev\r2k

# Clonar los tres repositorios uno al lado del otro
git clone https://github.com/adcondev/poster-tuis.git
git clone https://github.com/adcondev/scale-daemon.git
git clone https://github.com/adcondev/ticket-daemon.git

```

Después de clonar, verifique esta estructura de directorios:

```
C:\dev\r2k\
├── poster-tuis\       ← Instalador TUI (este proyecto)
├── scale-daemon\      ← Demonio del servicio Scale/Báscula
└── ticket-daemon\     ← Demonio del servicio de impresora de Tickets

```

### Configurar Variables de Entorno (.env)

El `Taskfile.yml` requiere un archivo `.env` en la raíz de `poster-tuis` para hashear e inyectar las credenciales por
seguridad.

Cree un archivo `.env` en `C:\dev\r2k\poster-tuis\.env` con el siguiente contenido (ajuste los valores según su
entorno):

```env
SCALE_DASHBOARD_PASSWORD=scale
SCALE_AUTH_TOKEN=mi_token_secreto
SCALE_PORT=8765
TICKET_DASHBOARD_PASSWORD=ticket
TICKET_AUTH_TOKEN=mi_token_secreto
TICKET_PORT=8766

```

### Inicializar Módulos de Go

```powershell
cd C:\dev\r2k\poster-tuis
go mod tidy

cd C:\dev\r2k\scale-daemon
go mod tidy

cd C:\dev\r2k\ticket-daemon
go mod tidy

```

---

## 3. Estructura del Proyecto

```
poster-tuis/
├── main.go                          # Punto de entrada (check de admin + tea.NewProgram)
├── go.mod                           # Definición del módulo
├── go.sum                           # Checksums de dependencias
├── Taskfile.yml                     # Orquestación de compilación
├── .env                             # Variables de seguridad (IGNORADO POR GIT)
├── .gitignore                       # Excluye bin/, dist/, internal/assets/bin/
│
├── internal/
│   ├── assets/
│   │   ├── assets.go                # Directivas //go:embed para los 4 binarios de servicio
│   │   └── bin/                     # ← Generado por Taskfile (ignorado por git)
│   │       ├── R2k_BasculaServicio_Local.exe
│   │       ├── R2k_BasculaServicio_Remote.exe
│   │       ├── R2k_TicketServicio_Local.exe
│   │       └── R2k_TicketServicio_Remote.exe
│   │
│   ├── config/
│   │   └── config.go                # BuildDate/Time, GenerateServiceNames(), GetBanner()
│   │
│   ├── service/
│   │   ├── registry.go              # Structs de variantes de servicio + GetServiceRegistry()
│   │   ├── manager.go               # Instalar/Desinstalar/Iniciar/Detener/Reiniciar
│   │   ├── status.go                # Enum de estatus + FamilyStatus (exclusividad mutua)
│   │   └── logs.go                  # Resolución de ruta de logs + abrir en Notepad/Explorer
│   │
│   └── ui/
│       ├── screens.go               # Enum de estado de pantalla (6 pantallas)
│       ├── builders.go              # Construcción de ítems de menú (impone exclusividad mutua)
│       ├── model.go                 # Struct del modelo + InitialModel() + ayudas de navegación
│       ├── update.go                # Update() + manejadores de teclas por pantalla + ayudas de acción
│       ├── view.go                  # View() + renderizadores de las 6 pantallas
│       └── styles.go                # Colores + estilos lipgloss
│
└── dist/                            # ← Generado por Taskfile (ignorado por git)
    └── R2k_POS_Instalador.exe       # Salida final (~15-20MB)

```

### Flujo de Dependencias Clave (Sin Importaciones Circulares)

```
main.go → internal/ui
internal/ui → internal/service, internal/config
internal/service → internal/assets, internal/config
internal/config → (solo stdlib)
internal/assets → (solo stdlib embed)

```

---

## 4. Proceso de Compilación

Todos los comandos se ejecutan desde el directorio `poster-tuis`.

```powershell
cd C:\dev\r2k\poster-tuis

```

### Paso 1: Limpiar Todo

```powershell
task clean:all

```

**Salida esperada:**

```
✅ Artefactos del instalador limpiados
✅ Todos los proyectos limpiados

```

### Paso 2: Compilar Binarios de Servicio para Incrustar

Esto compila los 4 archivos `.exe` de servicio desde los repositorios hermanos y los coloca en `internal/assets/bin/`:

```powershell
task installer:build:services

```

**Salida esperada:**

```
✅ Compilado R2k_BasculaServicio_Local.exe para el instalador
✅ Compilado R2k_BasculaServicio_Remote.exe para el instalador
✅ Compilado R2k_TicketServicio_Local.exe para el instalador
✅ Compilado R2k_TicketServicio_Remote.exe para el instalador

══════════════════════════════════════════════════════════════
  ✅ Todos los binarios de servicios están listos para integrarse
══════════════════════════════════════════════════════════════

```

**Verificar:**

```powershell
Get-ChildItem .\internal\assets\bin\

# Esperado: 4 archivos .exe
# R2k_BasculaServicio_Local.exe
# R2k_BasculaServicio_Remote.exe
# R2k_TicketServicio_Local.exe
# R2k_TicketServicio_Remote.exe

```

### Paso 3: Compilar el Instalador TUI

Esto compila la TUI en Go, incrustando los 4 binarios de servicio dentro del ejecutable final:

```powershell
task installer:build

```

**Salida esperada:**

```
══════════════════════════════════════════════════════════════
  ✅ COMPILACIÓN DEL INSTALADOR FINALIZADA
  📁 Salida: C:\dev\r2k\poster-tuis\dist
══════════════════════════════════════════════════════════════

```

**Verificar tamaño del archivo** (debe ser de ~15-20MB porque contiene 4 binarios de servicio incrustados):

```powershell
(Get-Item .\dist\R2k_POS_Instalador.exe).Length / 1MB
# Esperado: aproximadamente 15-20

```

### Compilación en un Comando (Pasos Combinados)

```powershell
task installer:rebuild

```

Esto ejecuta la limpieza de artefactos y luego reconstruye los binarios y el instalador final de forma secuencial.

---

## 5. Ejecución del Instalador

### ⚠️ CRÍTICO: Debe ejecutarse como Administrador

La TUI requiere privilegios de administrador para operaciones con `sc.exe` (control de servicios de Windows). Si se
lanza sin permisos de admin, muestra un mensaje de error y sale.

### Opciones de Lanzamiento

**Opción A: Método de clic derecho**

1. Navegue a la carpeta `dist\` en el Explorador de Archivos.
2. Haga clic derecho en `R2k_POS_Instalador.exe`.
3. Seleccione **"Ejecutar como administrador"**.

**Opción B: PowerShell como Admin**

```powershell
# Abra PowerShell como Administrador, luego:
cd C:\dev\r2k\poster-tuis
.\dist\R2k_POS_Instalador.exe

```

**Opción C: CMD como Admin**

```cmd
:: Abra CMD como Administrador, luego:
cd C:\dev\r2k\poster-tuis
dist\R2k_POS_Instalador.exe

```

### Lo que debería ver

Al iniciarse, la TUI muestra:

```
╔═════════════════════════════════════════════╗
║        SERVICE FAMILY MANAGER v2.0          ║
║       Build: 2026-02-05 14:30:22            ║
║                                             ║
║     Gestión de Servicios Red2000            ║
║     - Scale Service (Local/Remote)          ║
║     - Ticket Service (Local/Remote)         ║
║                                             ║
╚═════════════════════════════════════════════╝

SELECCIONE UNA FAMILIA DE SERVICIOS

  [1] Scale Service    No instalado
  [2] Ticket Service   No instalado
  [Q] Salir

scale: [-] | ticket: [-]

```

### Controles de Teclado

| Tecla       | Acción                                                      |
|-------------|-------------------------------------------------------------|
| `↑` / `k`   | Mover arriba                                                |
| `↓` / `j`   | Mover abajo                                                 |
| `Enter`     | Seleccionar elemento                                        |
| `r`         | Reinicio rápido (pantalla de familia, servicio debe correr) |
| `ESC` / `q` | Volver atrás / Salir                                        |
| `?`         | Alternar ayuda                                              |
| `Ctrl+C`    | Forzar salida (pantalla de procesamiento)                   |

### Flujo de Navegación

```
Dashboard (selector de familia)
  ├── Scale Service
  │     ├── [Sin servicio instalado] → Instalar Local / Instalar Remote
  │     └── [Servicio instalado] → Iniciar / Detener / Reiniciar / Logs / Desinstalar
  │           └── Logs → Abrir Archivo / Abrir Carpeta
  └── Ticket Service
        └── (misma estructura que Scale)

```

---

## 6. Escenarios de Prueba

Abra **Services.msc** (`Win+R` → `services.msc`) junto a la TUI para verificar.

### Prueba 1: Instalación Limpia — Scale Local

| Paso | Acción                           | Esperado                                                                       |
|------|----------------------------------|--------------------------------------------------------------------------------|
| 1    | Lanzar instalador                | Dashboard muestra ambas familias como "No instalado"                           |
| 2    | Seleccionar "Scale Service"      | Menú de familia muestra: "Instalar LOCAL", "Instalar REMOTE", "Volver"         |
| 3    | Selecc. "Instalar Versión LOCAL" | Diálogo de confirmación: "¿Instalar versión Local de Scale?"                   |
| 4    | Presionar `S` para confirmar     | Pantalla de procesamiento con spinner + barra de progreso                      |
| 5    | Esperar a que complete           | Resultado: "[+] Local instalado e iniciado correctamente"                      |
| 6    | Presionar `Enter` para seguir    | Menú de familia ahora muestra: Detener, Reiniciar, Ver Logs, Desinstalar, etc. |
| 7    | Verificar en Services.msc        | Servicio `R2k_BasculaServicio_Local` existe, estado RUNNING (En ejecución)     |

### Prueba 2: Exclusividad Mutua

| Paso | Acción                            | Esperado                                                   |
|------|-----------------------------------|------------------------------------------------------------|
| 1    | Con Scale Local instalado         | Menú de familia NO muestra "Instalar REMOTE"               |
| 2    | Selecc. "Desinstalar Local" → `S` | Procesando → "[-] Local desinstalado correctamente"        |
| 3    | Presionar `Enter`                 | Menú muestra "Instalar LOCAL" e "Instalar REMOTE" de nuevo |
| 4    | Verificar en Services.msc         | Servicio `R2k_BasculaServicio_Local` ya no existe          |

### Prueba 3: Ambas Familias Instaladas

| Paso | Acción                 | Esperado                                                                |
|------|------------------------|-------------------------------------------------------------------------|
| 1    | Instalar Scale Local   | Dashboard muestra "Scale Service: Local - [+] EN EJECUCION"             |
| 2    | Volver, selecc. Ticket | Menú muestra opciones de instalación                                    |
| 3    | Instalar Ticket Remote | Barra de estado del Dashboard: `scale: [+] Local                        |
| 4    | Verificar Services.msc | Ambos `R2k_BasculaServicio_Local` y `R2k_TicketServicio_Remote` existen |

### Prueba 4: Operaciones de Servicio (Detener / Iniciar / Reiniciar)

| Paso | Acción                     | Esperado                                                     |
|------|----------------------------|--------------------------------------------------------------|
| 1    | Con un servicio corriendo  | Menú muestra "Detener" y "Reiniciar"                         |
| 2    | Selecc. "Detener Servicio" | Resultado: "[OK] Detener Servicio completado"                |
| 3    | Volver al menú             | Menú ahora muestra "Iniciar" (Detener/Reiniciar desaparecen) |
| 4    | Selecc. "Iniciar Servicio" | Resultado: "[OK] Iniciar Servicio completado"                |
| 5    | Presionar atajo `r`        | Reinicio se ejecuta directamente (sin selección de menú)     |

### Prueba 5: Gestión de Logs

| Paso | Acción                          | Esperado                                           |
|------|---------------------------------|----------------------------------------------------|
| 1    | Con un servicio instalado       | Seleccionar "Ver Logs" desde menú de familia       |
| 2    | Selecc. "Abrir Archivo de Logs" | Se abre el Bloc de notas con el archivo `.log`     |
| 3    | Selecc. "Abrir Carpeta de Logs" | Se abre Explorer en `%PROGRAMDATA%\{ServiceName}\` |
| 4    | Presionar `ESC`                 | Regresa al menú de familia                         |

### Prueba 6: Sondeo de Estado en Segundo Plano

| Paso | Acción                                         | Esperado                                                                |
|------|------------------------------------------------|-------------------------------------------------------------------------|
| 1    | Dejar TUI abierta en pantalla de familia       | El servicio se muestra como "RUNNING"                                   |
| 2    | En Services.msc, clic derecho servicio → Stop  | —                                                                       |
| 3    | Esperar ~5 segundos                            | Menú TUI se auto-actualiza: "DETENIDO", muestra opción "Iniciar"        |
| 4    | En Services.msc, clic derecho servicio → Start | —                                                                       |
| 5    | Esperar ~5 segundos                            | Menú TUI se auto-actualiza: "EN EJECUCION", muestra "Detener/Reiniciar" |

### Prueba 7: Chequeo de Admin (Prueba Negativa)

| Paso | Acción                                             | Esperado                                           |
|------|----------------------------------------------------|----------------------------------------------------|
| 1    | Doble clic en `R2k_POS_Instalador.exe` (sin admin) | Muestra "[!] Permisos de Administrador Requeridos" |
| 2    | Presionar Enter                                    | La ventana se cierra                               |

---

## 7. Solución de Problemas

### Errores de Compilación

| Error                                                   | Causa                             | Solución                                                                       |
|---------------------------------------------------------|-----------------------------------|--------------------------------------------------------------------------------|
| `pattern bin/*.exe: no matching files found`            | `internal/assets/bin/` está vacío | Ejecute `task installer:build:services` primero                                |
| `cannot find package ".../internal/assets"`             | Módulo no resuelto                | Ejecute `go mod tidy` en `poster-tuis/`                                        |
| `undefined: assets.BasculaLocalBinary`                  | Nombre de paquete erróneo         | Verifique `package assets` (no `package main`) en `assets.go`                  |
| `config.ScaleBaseName undefined`                        | ldflags no se están inyectando    | Revise que `LDFLAGS_INSTALLER` en Taskfile incluya `-X '...ScaleBaseName=...'` |
| `cannot find module providing package .../scale-daemon` | Ruta hermana incorrecta           | Verifique que exista `../scale-daemon/` relativo a `poster-tuis/`              |
| La compilación tiene éxito pero el `.exe` pesa ~5MB     | Falló la incrustación             | Verifique que existan 4 archivos en `internal/assets/bin/` antes de compilar   |
| `Error cargando .env` o faltan variables                | Archivo `.env` ausente            | Asegúrese de haber creado `.env` en el root del proyecto como se indicó        |

### Errores en Tiempo de Ejecución

| Síntoma                                   | Causa                                     | Solución                                                     |
|-------------------------------------------|-------------------------------------------|--------------------------------------------------------------|
| TUI muestra error de admin y sale         | No se ejecuta como administrador          | Clic derecho → Ejecutar como administrador                   |
| Instala con éxito pero servicio no inicia | Binario de servicio necesita dependencias | Revise el log del servicio en `%PROGRAMDATA%\{ServiceName}\` |
| El estado muestra "DESCONOCIDO"           | Servicio está en estado de transición     | Espere unos segundos, el estado debería resolverse           |
| Notepad abre pero el archivo está vacío   | Servicio aún no ha escrito logs           | Inicie el servicio y espere actividad                        |

### Comandos de Limpieza

```powershell
# Eliminar todos los servicios instalados por la TUI (ejecutar como admin)
sc stop R2k_BasculaServicio_Local 2>$null
sc delete R2k_BasculaServicio_Local 2>$null
sc stop R2k_BasculaServicio_Remote 2>$null
sc delete R2k_BasculaServicio_Remote 2>$null
sc stop R2k_TicketServicio_Local 2>$null
sc delete R2k_TicketServicio_Local 2>$null
sc stop R2k_TicketServicio_Remote 2>$null
sc delete R2k_TicketServicio_Remote 2>$null

# Eliminar binarios instalados
Remove-Item -Recurse -Force "$env:ProgramFiles\R2k_*" -ErrorAction SilentlyContinue

# Limpiar artefactos de compilación
cd C:\dev\r2k\poster-tuis
task clean:all

```

### Verificar Estado Limpio

```powershell
# No deberían existir servicios R2k
sc query state= all | findstr "R2k_"
# Esperado: sin salida

# No deberían haber directorios R2k en Archivos de Programa
Get-ChildItem "$env:ProgramFiles\R2k_*"
# Esperado: sin resultados

# Sin artefactos de compilación
Test-Path .\dist\R2k_POS_Instalador.exe       # Esperado: False
Test-Path .\internal\assets\bin\               # Esperado: False

```