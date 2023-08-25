# Tastes


### Install required dependencies and tools

```
go mod tidy
brew install golangci-lint
```

### Run app:
You can start vscode debugger or by terminal:
```
go run ./cmd/api/main.go
```

### Run linters:

```
golangci-lint run --max-issues-per-linter=0 --max-same-issues=0 --config=./.golangci.yml

revive -config revive.toml -formatter friendly ./...
```

### Generate swagger documentation:

```
swag init -g routes.go --dir internal/handlers --parseDependency
```

### Check go vulnerabilities
```
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
```


# SolucionDynamo4: Diseño de Tabla de DynamoDB para Preferencias Musicales y Cinematográficas

## Descripción:
Este diseño tiene como objetivo almacenar y consultar las preferencias musicales y cinematográficas de los usuarios. La estructura se basa en dos tablas principales: `UserTastes` y `tastesTags`.

## Tabla: UserTastes

### Atributos:

- **Partition Key (PK)**: Combinación de `Type`, `Rank`, y `Name`.
 - Ejemplo: `SONG#1#Radioactive` o `GENRE#1#Rock`.

- **Sort Key (SK)**: `UserID`
  - Descripción: Identificador único para cada usuario.

- **Type**:
  - Descripción: Tipo de preferencia (ej. musical, cinematográfica).
  - Ejemplo: 'ARTIST', 'GENRE', 'SONG', 'MOVIE'.

- **Rank**:
  - Descripción: Número que representa el rango de preferencia.
  - Ejemplo: 1, 2, 3...

- **Name**:
  - Descripción: El nombre del artista, género, canción o película.
  - Ejemplo: "Imagine Dragons", "Rock", "Radioactive", "Inception".

- **TagID**:
  - Descripción: Referencia al ID en `tastesTags` que indica cómo y de dónde se obtuvo esta preferencia.
  - Ejemplo: "1" (que podría corresponder a Spotify).

- **created_at**:
  - Descripción: Fecha de creación del registro.

- **updated_at**:
  - Descripción: Fecha de última actualización del registro.

### Nota sobre `last_accessed_at`:
No se recomienda agregar esta columna directamente en la tabla. Si se desea rastrear el último acceso de un usuario, debe ser gestionado a través de logs.

### Ejemplos de Entradas:

| PK                        | SK         | Type  | Rank | Name                 | TagID | created_at | updated_at |
|---------------------------|------------|-------|------|----------------------|-------|------------|------------|
| SONG#1#unforgiven        | Sebas      | SONG  | 1    | Radioactive          | 1     | 2023-01-01 | 2023-01-01 |
| GENRE#1#Rock              | Joselin    | GENRE | 1    | Rock                 | 4     | 2023-01-01 | 2023-01-01 |
| MOVIE#1#Inception         | Maria      | MOVIE | 1    | Inception            | 5     | 2023-01-01 | 2023-01-01 |
| SONG#1#Radioactive        | Junior     | SONG  | 1    | Radioactive          | 1     | 2023-01-01 | 2023-01-01 |
| ARTIST#1#ImagineDragons   | Sebas      | ARTIST| 1    | Imagine Dragons      | 1     | 2023-01-02 | 2023-01-02 |
| SONG#2#Demons             | Sebas      | SONG  | 2    | Demons               | 1     | 2023-01-03 | 2023-01-03 |
| MOVIE#2#Interstellar      | Maria      | MOVIE | 2    | Interstellar         | 5     | 2023-01-04 | 2023-01-04 |
| GENRE#2#Pop               | Joselin    | GENRE | 2    | Pop                  | 4     | 2023-01-05 | 2023-01-05 |
| SONG#3#Thunder            | Junior     | SONG  | 3    | Thunder              | 1     | 2023-01-06 | 2023-01-06 |
| MOVIE#3#DarkKnight        | Maria      | MOVIE | 3    | The Dark Knight      | 5     | 2023-01-07 | 2023-01-07 |
| ARTIST#2#Coldplay         | Junior     | ARTIST| 2    | Coldplay             | 2     | 2023-01-08 | 2023-01-08 |
| GENRE#3#Jazz              | Sebas      | GENRE | 3    | Jazz                 | 3     | 2023-01-09 | 2023-01-09 |


### Índice Secundario Global para UserTastes

#### Recomendación: 

Para facilitar consultas que buscan recuperar todos los gustos de un usuario específico, recomendamos agregar un Índice Secundario Global (GSI) en la tabla `UserTastes` utilizando `UserID` como la clave principal.

#### Especificaciones:

- **Nombre del GSI**: `UserIDIndex`
- **Clave Principal (PK) del GSI**: `UserID` (anteriormente SK en la tabla principal)
- **Clave de Ordenación (SK) del GSI**: `PK` de la tabla principal (combinación de `Type`, `Rank`, y `Name`).

#### Beneficios:

- Permitirá consultas eficientes para recuperar todos los gustos de un usuario específico sin tener que especificar el tipo o rango.
- Facilita el análisis de los gustos individuales de los usuarios.

#### Ejemplo de Uso:

Si deseas obtener todos los gustos del usuario "Sebas", simplemente realizarías una consulta en el `UserIDIndex` con `UserID = Sebas`.


## DynamoDB Streams para Historial de Cambios de Gustos

### Implementación:

1. **Habilitar DynamoDB Streams en la tabla `UserTastes`**:
    - Ve a la sección "DynamoDB Streams" en la configuración de la tabla.
    - Activa la opción y selecciona "Only capture modifications".

2. **Configurar AWS Lambda**:
    - Crea una nueva función Lambda.
    - Configura un disparador ("trigger") para esta función, seleccionando el stream de DynamoDB que acabas de crear.
    - Implementa el código de la función para procesar los cambios y almacenarlos en una tabla de backup o histórico.

3. **Tabla de Backup o Histórico `UserTastesHistory`**:

| PK                  | SK    | Type  | Rank | Name          | TagID | created_at | updated_at | ChangeDate          |
|---------------------|-------|-------|------|---------------|-------|------------|------------|---------------------|
| SONG#Radioactive    | Sebas | SONG  | 1    | Radioactive   | 1     | 2023-01-01 | 2023-05-01 | 2023-05-01 15:23:00 |
| GENRE#Rock          | Joselin | GENRE | 1    | Rock         | 4     | 2023-01-01 | 2023-06-10 | 2023-06-10 10:10:10 |

Nota: El campo `ChangeDate` representa la fecha y hora exacta del cambio registrado en la tabla principal.


---

## Tabla: tastesTags

### Atributos:

- **PK**: Identificador único para cada etiqueta.
- **Config**: Formato de la etiqueta. Ejemplo: `SONG#RANK#SONGID`.
- **DESCRIP**: Descripción breve de la etiqueta.
- **created_at**: Fecha de creación de la entrada.
- **version**: Versión de la etiqueta.
- **updated_at**: Fecha de última actualización.
- **status**: Estado de la etiqueta (activo, inactivo, obsoleto).
- **long_description**: Descripción detallada de la etiqueta.
- **IntegrationRef**: Identificador de la plataforma o integración de donde proviene la información (ej. Spotify, IMDB).

### Ejemplos de Entradas:

| PK | Config                | DESCRIP                                | created_at | version | updated_at | status | long_description                             | IntegrationRef | pricing  |
|----|-----------------------|----------------------------------------|------------|---------|------------|--------|----------------------------------------------|----------------|-----------|
| 1  | SONG#1#SONGID         | Favorite song rank 1                   | 2023-01-01 | 1.0     | 2023-01-01 | active | Rank of favorite song for each user.        | Spotify        | FREE      |
| 2  | SONG#2#SONGID         | Favorite song rank 2                   | 2023-01-01 | 1.0     | 2023-01-01 | active | Rank of 2 favorite songs for each user.     | Spotify        | PREMIUM   |
| 3  | SONG#3#SONGID         | Favorite song rank 3                   | 2023-01-01 | 1.0     | 2023-01-01 | active | Rank of 3 favorite songs for each user.     | Spotify        | PREMIUM   |
| 4  | GENRE#RANK#RANKID     | Genre                                  | 2023-01-01 | 1.0     | 2023-01-01 | active | Preferred genre rank for users.              | YouTubeMusic   | FREE      |
| 5  | MOVIE#RANK#MOVIEID    | Favorite movie rank                    | 2023-01-01 | 1.0     | 2023-01-01 | active | Rank of favorite movies for users.           | IMDB           | FREE      |
| 5  | FOLLOW#-#ArtistID    | Favorite movie rank                    | 2023-01-01 | 1.0     | 2023-01-01 | active | Rank of favorite movies for users.           | IMDB           | FREE      |
---

## consideraciones para la tabla tastesTags:


1. **Status**: Añade un atributo de `status` para gestionar la activación de etiquetas.
2. **Automatización**: Automatiza la adición de nuevas entradas si sigue un patrón.
3. **Validación**: Implementa validaciones al añadir o modificar entradas.

---

