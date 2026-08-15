# Objetivo

Añadir versionado explícito a IASI y hacer que esa versión forme parte tanto del ejecutable como de las instalaciones `.iasi`.

Actualmente disponemos de una CLI Go autocontenida que incluye la metodología IASI embebida y permite:

```bash
iasi install --workspace
iasi status
```

Queremos que cada binario represente una **distribución concreta y reproducible de IASI**.

La versión inicial será:

```text
0.1.0
```

No implementar todavía mecanismos de actualización ni múltiples versiones de los componentes internos.

---

# Principio

Una versión de IASI identifica conjuntamente:

* la metodología contenida en `agentics/`;
* los schemas correspondientes;
* los adaptadores que existan en esa versión;
* y las herramientas necesarias para instalar y utilizar esa metodología.

El ejecutable `iasi` debe conocer qué versión contiene.

Conceptualmente:

```text
IASI repository
      │
      ├── VERSION
      │
      ├── agentics/
      │
      └── src/go/
               │
               ▼
            iasi.exe
            IASI 0.1.0
               │
               ▼
        iasi install --workspace
               │
               ▼
            .iasi/
          version 0.1.0
```

---

# Archivo `VERSION`

Crear en la raíz del repositorio:

```text
iasi/
├── VERSION
├── README.md
├── README_en.md
├── agentics/
└── src/
```

Contenido inicial:

```text
0.1.0
```

`VERSION` debe ser la **única fuente de verdad de la versión de distribución de IASI**.

No duplicar manualmente `0.1.0` como constante independiente en distintos lugares del código.

---

# Versión embebida

El ejecutable debe contener la versión durante la compilación, igual que contiene actualmente los artefactos metodológicos.

La ejecución del binario no debe requerir acceso al archivo `VERSION` original del repositorio.

Debe seguir siendo posible distribuir exclusivamente:

```text
iasi.exe
```

y disponer tanto de:

* la metodología;
* como de la versión correspondiente.

Si el proceso de build actual ya prepara contenido embebido antes de compilar, incluir también `VERSION` en ese proceso.

La raíz del repositorio sigue siendo la fuente de verdad.

---

# Nuevo comando `iasi version`

Implementar:

```bash
iasi version
```

Salida:

```text
IASI 0.1.0
```

La versión mostrada debe proceder de la información embebida en el ejecutable.

No leer `VERSION` desde el filesystem del usuario en tiempo de ejecución.

---

# `iasi install --workspace`

Modificar la generación de:

```text
.iasi/manifest.yml
```

para que utilice la versión real contenida en el ejecutable.

Ejemplo:

```yaml
version: 0.1.0
profile: workspace

installed:
  instructions: all
  commands: all
  skills: all
  mcp: all
```

No mantener una constante de versión separada para el manifest.

La versión del manifest debe ser exactamente la versión de IASI contenida en el ejecutable que realizó la instalación.

---

# `iasi status`

Mantener la información actual y añadir, si no existe ya, la distinción entre:

* versión instalada en `.iasi`;
* versión del ejecutable que está ejecutando `status`.

Ejemplo cuando coinciden:

```text
IASI

Type      : workspace
Path      : C:/workspace/.iasi
Installed : 0.1.0
Binary    : 0.1.0

Instructions : 14
Commands     : 0
Skills       : 0
MCP          : 0
```

Ejemplo cuando son diferentes:

```text
IASI

Type      : workspace
Path      : C:/workspace/.iasi
Installed : 0.1.0
Binary    : 0.2.0

Instructions : 14
Commands     : 0
Skills       : 0
MCP          : 0
```

Una diferencia de versiones **no debe considerarse todavía un error**.

Solo debe mostrarse.

No implementar aún:

```bash
iasi update
```

ni recomendaciones automáticas de actualización.

---

# Versionado

IASI utilizará Semantic Versioning.

Durante la fase actual permaneceremos en versiones `0.x.y`.

Orientación inicial:

```text
0.1.0 → 0.1.1
```

Correcciones compatibles, ajustes menores y defectos que no cambien de forma relevante el modelo metodológico.

```text
0.1.0 → 0.2.0
```

Nuevas capacidades, instrucciones, comandos, skills, adaptadores o cambios metodológicos compatibles.

```text
1.x.x → 2.0.0
```

Cambios incompatibles en schemas, estructura o semántica una vez que IASI alcance estabilidad.

No implementar lógica SemVer en la CLI en esta iteración.

El archivo `VERSION` es suficiente.

---

# Versiones internas de instructions

Actualmente algunas instrucciones pueden contener metadata como:

```yaml
version: 0.1.0
```

No desarrollar todavía un sistema independiente de versionado para cada instruction.

La unidad de distribución, instalación y compatibilidad será:

```text
IASI <version>
```

No intentar resolver matrices de versiones entre instructions, skills, commands, adapters y CLI.

Ese problema se abordará únicamente si aparece una necesidad real.

---

# Build autocontenido

El requisito existente se mantiene:

```text
iasi.exe
```

debe ser autocontenido.

Después de compilarlo debemos poder copiar únicamente el ejecutable a un directorio vacío y ejecutar:

```bash
iasi version
```

obteniendo:

```text
IASI 0.1.0
```

y posteriormente:

```bash
iasi install --workspace
```

sin acceso al repositorio fuente.

La instalación debe contener la metodología correspondiente exactamente a esa versión.

---

# Tests

Añadir tests para verificar como mínimo:

## Version

```bash
iasi version
```

utiliza la versión embebida.

## Manifest

Una instalación nueva escribe en:

```text
.iasi/manifest.yml
```

la misma versión contenida en el ejecutable.

## Status

`status` puede mostrar de forma independiente:

```text
Installed
Binary
```

y funciona tanto cuando ambas versiones coinciden como cuando son diferentes.

## Independencia del repositorio

Las operaciones de runtime no dependen de:

```text
VERSION
agentics/
```

externos al ejecutable.

---

# Criterios de aceptación

Debe seguir funcionando:

```bash
go test ./...
```

y:

```bash
go build -o iasi.exe ./cmd/iasi
```

Después, copiando únicamente:

```text
iasi.exe
```

a un workspace vacío:

```bash
iasi version
```

debe devolver:

```text
IASI 0.1.0
```

Después:

```bash
iasi install --workspace
```

debe crear una instalación cuyo:

```text
.iasi/manifest.yml
```

contenga:

```yaml
version: 0.1.0
```

Y:

```bash
iasi status
```

debe mostrar:

```text
Installed : 0.1.0
Binary    : 0.1.0
```

junto con la información existente de la instalación.

---

# Fuera de alcance

No implementar en esta iteración:

* `iasi update`;
* descarga de releases;
* comparación SemVer;
* actualización automática;
* rollback;
* migraciones;
* versionado independiente de componentes;
* adaptadores;
* `resolve`;
* `validate`.

El objetivo de esta iteración es únicamente establecer:

```text
VERSION
   ↓
iasi.exe
   ↓
.iasi/manifest.yml
```

como cadena única y trazable de versionado.
