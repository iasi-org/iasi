La CLI ya funciona sin depender del repositorio externo, pero:

```bash
iasi install --workspace
```

crea las carpetas de `.iasi/` vacías.

Esto no es correcto.

Actualmente existen archivos reales en:

```text
agentics/instructions/
```

y deben quedar instalados en:

```text
.iasi/instructions/
```

conservando exactamente su estructura.

Por ejemplo:

```text
agentics/instructions/general/behavior.md
```

debe producir:

```text
.iasi/instructions/general/behavior.md
```

Lo mismo para:

```text
schema/
general/
documentation/
code/
diagrams/
```

`commands/`, `skills/` y `mcp/` pueden quedar vacíos mientras no tengan contenido fuente.

## Revisar

Comprobar que el proceso de build:

1. sincroniza realmente `agentics/` con el directorio utilizado por `go:embed`;
2. incluye los archivos, no solo la estructura conceptual;
3. utiliza un patrón `go:embed` que incluya recursivamente todos los archivos;
4. durante `install`, recorre el `embed.FS` y escribe cada archivo en `.iasi/`;
5. conserva rutas relativas y subdirectorios.

## Importante

`go:embed` no representa directorios vacíos.

Por tanto:

* los directorios con archivos aparecerán al extraer esos archivos;
* `commands`, `skills` y `mcp` deben crearse explícitamente si están vacíos.

## Criterio de aceptación

Después de compilar `iasi.exe`, copiar únicamente el ejecutable a un workspace vacío y ejecutar:

```bash
iasi install --workspace
```

Debe generar, como mínimo:

```text
.iasi/
├── manifest.yml
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

Los archivos de `instructions/` deben contener exactamente el contenido incluido en la metodología fuente.

Añadir un test que verifique al menos un archivo representativo, por ejemplo:

```text
.iasi/instructions/general/behavior.md
```

y compruebe que no está vacío.
