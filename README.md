**Castellano** | [🇬🇧 English](README_en.md)

# iasi

**Metodología de Ingeniería Asistida por Sistemas Inteligentes**

`iasi` es el repositorio donde se define, formaliza y evoluciona la metodología IASI.

El ecosistema IASI contiene libros, productos, herramientas, documentación, infraestructura y memoria de ingeniería. Este repositorio cumple una función diferente: contiene las reglas, estructuras y artefactos que describen **cómo trabajamos con Sistemas Inteligentes para hacer ingeniería**.

No está ligado a un modelo, proveedor o herramienta concreta. La metodología debe poder expresarse de forma independiente y, cuando sea necesario, adaptarse después a entornos como Codex, Copilot, Claude u otros sistemas.

---

# Propósito

El objetivo de `iasi` es convertir una forma de trabajo descubierta mediante la práctica en una metodología explícita, reproducible y cuestionable.

Aquí se formalizan aspectos como:

- cómo debe comportarse un agente;
- qué instrucciones debe respetar;
- qué capacidades pueden reutilizarse;
- qué acciones pueden exponerse como comandos;
- cómo se integran herramientas externas;
- cómo se valida el trabajo realizado;
- qué decisiones requieren intervención humana;
- y cómo se conserva el conocimiento necesario para que otro humano o agente pueda continuar el trabajo.

La metodología no pretende eliminar el juicio humano ni sustituir la ingeniería.

Pretende crear un marco en el que humanos y Sistemas Inteligentes puedan colaborar de forma controlada, trazable y reproducible.

---

# Nuestra aproximación

IASI no parte de una metodología cerrada para después intentar encajar los problemas dentro de ella.

Partimos de los problemas, observamos cómo los resolvemos y formalizamos únicamente aquello que demuestra ser útil.

> **No adoptamos metodologías para resolver problemas. Partimos de los problemas para descubrir la metodología.**

Por tanto, `iasi` no debe convertirse en una colección de reglas teóricas desconectadas de la práctica.

Cada elemento incorporado a la metodología debe responder a una necesidad real, poder explicarse y estar sujeto a revisión.

La metodología evoluciona con la experiencia.

---

# Estructura IASI

Los artefactos canónicos de IASI viven bajo `iasi/`. Definen cómo participan
los agentes de IA en el proceso de ingeniería.

```text
iasi/
├── instructions/
├── commands/
├── skills/
├── mcp/
└── adapters/
```

Cada pieza cumple una función distinta.

| Área | Propósito |
|------|-----------|
| `instructions/` | Define cómo debe comportarse un agente y qué reglas debe respetar al producir trabajo. |
| `commands/` | Define acciones explícitas que pueden solicitarse a los agentes. |
| `skills/` | Define capacidades reutilizables que pueden aplicarse en distintos contextos. |
| `mcp/` | Define la integración metodológica con herramientas y servicios expuestos mediante Model Context Protocol. |

Estas definiciones son conceptuales y deben permanecer independientes de la implementación concreta de cada plataforma.

---

# Instrucciones

`iasi/instructions/` contiene reglas persistentes que gobiernan el comportamiento de los agentes.

Las instrucciones pueden definir, entre otras cosas:

- comportamiento general;
- control humano;
- tratamiento de incertidumbre;
- validación;
- uso de herramientas;
- precedencia entre reglas;
- estilo de código;
- pruebas;
- estructura y estilo documental;
- tratamiento de fuentes;
- convenciones para diagramas.

Las instrucciones se diseñan para ser:

- **atómicas**, centradas en un único aspecto;
- **declarativas**, sin depender de un proveedor;
- **componibles**, de modo que puedan combinarse según el trabajo;
- **observables**, para que su cumplimiento pueda revisarse;
- **versionables**, porque la metodología también evoluciona.

Su estructura común se define en:

```text
iasi/instructions/schema/instructions.md
```

`iasi/commands/validate.md`, `iasi/commands/archive.md` y
`iasi/commands/plan.md` definen los comandos agénticos canónicos `/validate`,
`/archive` y `/plan`. Los adapters pueden proyectarlos a mecanismos nativos,
pero no redefinirlos.

## CLI

```text
iasi install
iasi reinstall
iasi status
iasi version
iasi adapt copilot
```

Una instalación local válida se identifica por `.iasi/manifest.yml`. Las capas
instaladas en directorios padre se combinan con las locales; `validation.json`
es estado local de workflow, no una capa instalada.

---

# Independencia de plataforma

IASI distingue entre **la metodología** y **la forma concreta en que una herramienta la consume**.

Por ejemplo, una misma instrucción puede terminar expresándose como:

- instrucciones de repositorio;
- archivos de configuración;
- prompts;
- skills;
- agentes especializados;
- reglas de un entorno de desarrollo;
- o cualquier otro mecanismo que soporte una plataforma.

La fuente de verdad metodológica permanece en `iasi`.

Las adaptaciones a herramientas concretas no deben determinar el diseño de la metodología.

---

# Humano y Sistema Inteligente

La disponibilidad de un Sistema Inteligente no implica disponibilidad ilimitada del humano.

IASI considera que la velocidad sostenible de un sistema de ingeniería no viene determinada únicamente por la velocidad de ejecución de los agentes, sino por la capacidad humana de:

- comprender;
- decidir;
- revisar;
- validar;
- y asimilar los resultados.

La automatización debe reducir trabajo mecánico, no convertir cada pausa eliminada en una nueva decisión para el humano.

El objetivo no es maximizar actividad.

El objetivo es mejorar la capacidad de hacer buena ingeniería.

---

# Validación

En IASI, ejecutar una tarea no equivale a completarla.

La metodología distingue entre:

```text
implemented
validated
accepted
```

Una solución puede estar implementada sin haber sido validada.

Puede estar validada técnicamente sin haber sido aceptada por el humano.

Los criterios de validación deben apoyarse siempre que sea posible en resultados observables, pruebas y criterios de aceptación.

---

# Principios

- La ingeniería es el objetivo.
- Los Sistemas Inteligentes son asistentes de ingeniería.
- El humano conserva el control sobre objetivos, prioridades y aceptación.
- La metodología debe ser independiente de modelos y proveedores.
- El conocimiento debe persistirse de forma explícita.
- Las instrucciones deben ser componibles y verificables.
- La automatización debe reducir complejidad para el usuario, no trasladársela.
- La implementación no sustituye a la validación.
- Las decisiones deben poder cuestionarse.
- La metodología también evoluciona.
- Si la implementación cuesta cada vez menos, el pensamiento vale cada vez más.

---

# Estado del repositorio

`iasi` se construye de forma incremental.

No intentamos diseñar toda la metodología por adelantado.

Primero formalizamos las piezas que ya han aparecido como necesarias en proyectos reales. Después las utilizamos, validamos y modificamos.

Una estructura se considera útil cuando resiste su aplicación en distintos casos sin obligarnos a introducir excepciones artificiales.

Por eso, el contenido de este repositorio debe entenderse como una metodología viva.

---

# Relación con el ecosistema IASI

`iasi` define la metodología.

Los demás repositorios del ecosistema proporcionan los lugares donde esa metodología se utiliza, se prueba, se documenta o se convierte en producto.

Los libros explican las ideas.

Los productos materializan soluciones.

La memoria de ingeniería conserva el camino recorrido.

`iasi` intenta convertir lo aprendido durante ese camino en una forma de trabajo explícita y reutilizable.

---

# Un proyecto abierto

IASI nace de la práctica y debe seguir pudiendo ser cuestionado desde la práctica.

Las reglas de este repositorio no son dogmas. Son decisiones de ingeniería que deben poder justificarse, probarse, refinarse o descartarse cuando exista una alternativa mejor.

**Toda aportación fundamentada será bienvenida.**

---

> *"La mejor metodología no es la que seguimos, sino la que descubrimos mientras resolvemos el problema adecuado."*

---

**IASI**

*Ingeniería Asistida por Sistemas Inteligentes*
