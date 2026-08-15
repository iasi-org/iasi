# Objetivo

Implementar la primera versión de la CLI `iasi` en Go.

El repositorio `iasi` contiene la metodología IASI de forma declarativa. La CLI no debe contener ni codificar la metodología: debe operar sobre los artefactos existentes en el repositorio.

En esta primera iteración queremos únicamente dos comandos:

```bash
iasi install --workspace
iasi status
```

No implementar todavía:

* `validate`
* `resolve`
* adapters para Codex, Copilot, Claude, etc.
* instalación de proyecto
* actualización de instalaciones
* descarga remota
* resolución de dependencias
* mecanismos de plugins

Mantener esta primera versión deliberadamente simple.

# Contexto

La estructura conceptual actual del repositorio es:

```text
iasi/
├── README.md
├── README_en.md
└── agentics/
    ├── instructions/
    │   ├── README.md
    │   ├── schema/
    │   │   └── instructions.md
    │   ├── general/
    │   ├── documentation/
    │   ├── code/
    │   └── diagrams/
    ├── commands/
    ├── skills/
    └── mcp/
```

Algunas de las carpetas de `agentics/` pueden todavía no existir o estar vacías.

Esto no debe considerarse un error.

# Principio arquitectónico

Separar estrictamente:

```text
metodología declarativa
        ↓
     CLI iasi
        ↓
 instalación .iasi
```

La metodología vive en archivos bajo `agentics/`.

La CLI es únicamente el motor que instala, inspecciona y, en futuras iteraciones, compondrá y validará esos artefactos.

No introducir reglas metodológicas dentro del código Go.

# Tecnología

Implementar la CLI en Go.

La aplicación debe generar un único binario:

```text
iasi
```

o en Windows:

```text
iasi.exe
```

Evitar dependencias externas salvo que aporten un beneficio claro.

Para esta primera versión se prefiere utilizar la biblioteca estándar de Go siempre que sea razonable.

# Comando `iasi install --workspace`

Debe instalar IASI en el directorio de trabajo actual.

Ejemplo:

```bash
cd C:/workspace
iasi install --workspace
```

Debe producir:

```text
C:/workspace/
└── .iasi/
    ├── manifest.yml
    ├── instructions/
    ├── commands/
    ├── skills/
    └── mcp/
```

## Comportamiento

La instalación de tipo `workspace` debe incorporar todo el contenido disponible de:

```text
agentics/instructions/
agentics/commands/
agentics/skills/
agentics/mcp/
```

La estructura interna debe conservarse.

Ejemplo:

```text
agentics/instructions/general/behavior.md
```

debe terminar como:

```text
.iasi/instructions/general/behavior.md
```

Si una categoría como `commands`, `skills` o `mcp` todavía no existe en el repositorio fuente, crear igualmente su directorio vacío en `.iasi/`.

## Manifest

Crear:

```text
.iasi/manifest.yml
```

Con una estructura inicial sencilla, por ejemplo:

```yaml
version: 0.1.0
profile: workspace

installed:
  instructions: all
  commands: all
  skills: all
  mcp: all
```

No diseñar todavía un schema complejo para el manifest.

La versión puede mantenerse inicialmente como una constante de la CLI:

```text
0.1.0
```

hasta que definamos posteriormente el mecanismo definitivo de versionado.

# Reinstalación

No sobrescribir silenciosamente una instalación existente.

Si:

```text
.iasi/
```

ya existe, el comando debe detenerse con un mensaje claro.

Ejemplo conceptual:

```text
IASI is already installed in this workspace:
C:/workspace/.iasi
```

No implementar todavía `--force`, `update` ni mecanismos de merge.

# Comando `iasi status`

Debe localizar una instalación `.iasi` aplicable al directorio actual.

Inicialmente basta con buscar:

1. `.iasi` en el directorio actual.
2. Si no existe, ascender por los directorios padre hasta encontrar una.

Esto permitirá ejecutar:

```bash
C:/workspace/project-a/src> iasi status
```

y detectar:

```text
C:/workspace/.iasi
```

## Salida

Mostrar una salida humana, sencilla y estable.

Ejemplo:

```text
IASI

Type    : workspace
Path    : C:/workspace/.iasi
Version : 0.1.0

Instructions : 14
Commands     : 0
Skills       : 0
MCP          : 0
```

Los contadores deben calcularse a partir de los archivos realmente presentes en cada categoría.

No hace falta diferenciar todavía tipos de archivo ni interpretar el contenido.

Si no existe ninguna instalación:

```text
IASI is not installed for this location.
```

y devolver un código de salida distinto de cero.

# Localización de la metodología fuente

La CLI necesita acceder al contenido de `agentics/` para realizar la instalación.

En esta primera versión, diseñar esta parte de forma sencilla pero aislada para que posteriormente podamos cambiar el mecanismo de distribución.

No dispersar por el código referencias directas a rutas de `agentics`.

Debe existir un único componente responsable de localizar la fuente de la metodología.

Para desarrollo, puede resolverse inicialmente desde la estructura del propio repositorio.

Si para ejecutar un binario completamente independiente fuese necesario incrustar los archivos mediante `embed`, puede utilizarse `go:embed`.

En ese caso:

* mantener intacta la estructura de archivos;
* considerar los archivos embebidos como datos;
* no convertir su contenido en código Go.

# Estructura interna

No sobrearquitecturar.

Una estructura razonable sería:

```text
cmd/
└── iasi/
    └── main.go

internal/
├── install/
├── status/
├── manifest/
└── source/
```

Puede simplificarse si alguna de estas capas resulta artificial en esta primera versión.

Priorizar claridad sobre número de paquetes.

# CLI

La sintaxis debe ser:

```bash
iasi install --workspace
iasi status
```

Para esta versión no necesitamos un framework complejo de CLI.

Puede utilizarse `flag` o parsing sencillo con biblioteca estándar.

Los errores deben:

* ser claros para el usuario;
* escribirse en stderr;
* devolver código de salida distinto de cero.

# Compatibilidad

Debe funcionar al menos en:

* Windows
* Linux

No asumir separadores `/` o `\`.

Utilizar las funciones estándar de manejo de paths de Go.

# Tests

Añadir tests desde esta primera versión.

Como mínimo:

## Install

Validar que:

* crea `.iasi`;
* crea `manifest.yml`;
* copia `instructions`;
* conserva subdirectorios;
* crea `commands`, `skills` y `mcp` aunque estén vacíos;
* falla si `.iasi` ya existe.

## Status

Validar que:

* encuentra `.iasi` en el directorio actual;
* encuentra `.iasi` ascendiendo desde un subdirectorio;
* lee la versión del manifest;
* cuenta los archivos de cada categoría;
* falla correctamente cuando no existe instalación.

Los tests deben utilizar directorios temporales.

No deben depender del workspace real del desarrollador.

# Criterios de aceptación

La iteración se considera terminada cuando podamos ejecutar:

```bash
go test ./...
```

y todos los tests pasen.

Y podamos construir:

```bash
go build ./cmd/iasi
```

Después, desde un workspace vacío:

```bash
iasi install --workspace
```

debe crear correctamente:

```text
.iasi/
```

con el contenido metodológico.

A continuación:

```bash
iasi status
```

debe identificar esa instalación y mostrar su información.

Desde un subdirectorio del mismo workspace:

```bash
iasi status
```

debe encontrar la misma instalación ascendiendo por el árbol de directorios.

# Restricción principal

No implementar nada más allá de este alcance salvo que sea estrictamente necesario para que estos dos comandos funcionen correctamente.

Queremos validar primero el modelo:

```text
iasi
  → install
  → .iasi
  → status
```

antes de añadir más capacidades.
