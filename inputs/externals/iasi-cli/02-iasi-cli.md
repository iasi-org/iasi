La implementación actual busca el directorio `agentics/` en el filesystem durante la ejecución.

Esto no es correcto para el modelo de distribución de `iasi`.

## Requisito

El ejecutable:

```text
iasi.exe
```

debe ser completamente autocontenido.

Debe incluir dentro del propio binario todos los artefactos metodológicos necesarios para ejecutar:

```bash
iasi install --workspace
```

sin necesitar:

* el repositorio `iasi`;
* un directorio `agentics/` externo;
* archivos auxiliares junto al ejecutable;
* rutas relativas al código fuente.

Debe ser posible copiar únicamente:

```text
iasi.exe
```

a un workspace vacío y ejecutar:

```bash
iasi install --workspace
```

## Implementación

Utilizar `go:embed` / `embed.FS`.

El contenido fuente sigue viviendo en:

```text
iasi/
└── agentics/
    ├── instructions/
    ├── commands/
    ├── skills/
    └── mcp/
```

Durante la compilación debe incorporarse al binario.

La CLI debe leer la metodología desde el filesystem embebido, no desde el filesystem del usuario.

Conceptualmente:

```go
//go:embed ...
var methodology embed.FS
```

## Importante

La ubicación actual del módulo Go es:

```text
src/go
```

y `agentics/` está en la raíz del repositorio.

`go:embed` no permite referenciar archivos fuera del árbol del paquete mediante rutas `../`.

Por tanto, resolver este problema de forma limpia.

Puede utilizarse un paso de build que copie/sincronice los artefactos necesarios a un directorio interno del módulo antes de compilar, por ejemplo:

```text
src/go
└── embedded/
    └── agentics/
```

y posteriormente:

```go
//go:embed embedded/agentics/**
```

Pero:

* `agentics/` en la raíz sigue siendo la fuente de verdad;
* `embedded/agentics/` es únicamente material de build;
* no mantener manualmente dos copias divergentes;
* el proceso de build debe regenerar esa copia automáticamente.

Si existe una solución más simple y limpia respetando estas restricciones, utilizarla.

## Resultado esperado

Tras:

```bash
go build -o iasi.exe ./cmd/iasi
```

debo poder copiar solamente:

```text
iasi.exe
```

a un directorio completamente vacío.

Después:

```bash
iasi install --workspace
```

debe crear:

```text
.iasi/
├── manifest.yml
├── instructions/
├── commands/
├── skills/
└── mcp/
```

sin buscar ningún `agentics/` externo.

## Test de aceptación

Añadir un test que garantice que la instalación utiliza exclusivamente la fuente embebida y no depende de la presencia del repositorio fuente en tiempo de ejecución.

Eliminar el comportamiento actual que produce:

```text
could not locate IASI source directory (agentics)
```

en un ejecutable distribuido.
