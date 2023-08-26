# Geo


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

## Endpoints:

### 1. Guardar Ubicación de Usuario:
- **URL**: `/location`
- **Method**: `POST`
- **Descripción**:  
Este endpoint permite guardar la ubicación actual del usuario, pero en lugar de guardar la latitud y longitud exactas, mapea esas coordenadas a una zona específica, como un barrio, ciudad o punto de interés.
  
- **Parámetros**:  
  - **userID**: Identificador del usuario.
  - **lat**: Latitud actual del usuario.
  - **long**: Longitud actual del usuario.

- **Respuesta**:  
Un mensaje de confirmación si la ubicación fue guardada con éxito, o un mensaje de error en caso contrario.

### 2. Buscar Usuarios Cercanos con Gustos Similares:
- **URL**: `/location/nearbies`
- **Method**: `GET`
- **Descripción**:  
Este endpoint permite encontrar usuarios que se encuentren en la misma zona y que tengan gustos similares basados en un filtro proporcionado. Primero, filtra a los usuarios que están en la misma zona, y luego consulta a otro servicio (Tastes) para verificar si tienen gustos similares según el filtro proporcionado.

- **Parámetros**:  
  - **userID**: Identificador del usuario.
  - **lat**: Latitud actual del usuario.
  - **long**: Longitud actual del usuario.
  - **filter**: Tipo de filtro a aplicar (por ejemplo, canción más escuchada, género musical, etc.).

- **Respuesta**:  
Una lista de usuarios cercanos que cumplen con el criterio de búsqueda.

## Diseño de Tablas PostgreSQL con Amazon RDS:

### Tabla: ZonesConfig
- **Atributos**:
  - **ZoneID** (PK): Identificador único para cada zona.
  - **ZoneName**: Nombre de la zona o punto de interés.
  - **ZoneType**: Tipo de la zona (ej., country, POI, neighborhood, city).
  - **created_at**: Fecha de creación del registro.
  - **updated_at**: Fecha de la última actualización del registro.

**Ejemplo**:  
| ZoneID | ZoneName            | ZoneType    |
|--------|---------------------|-------------|
| Z001   | New York City       | city        |
| Z002   | Central Park        | POI         |
| Z003   | Brooklyn            | neighborhood|
| Z004   | Los Angeles         | city        |
| Z005   | Hollywood           | POI         |
| Z006   | Santa Monica        | neighborhood|
| Z007   | Buenos Aires        | city        |
| Z008   | River Plate Stadium | POI         |
| Z009   | Villa Crespo        | neighborhood|
| Z010   | Tokyo               | city        |

### Tabla: UserZones
- **Atributos**:
  - **ZoneID** (PK): Identificador único para cada zona.
  - **UserID** (SK): Identificador único para cada usuario.
  - **Timestamp**: Fecha y hora exactas en que se guardó la zona para el usuario.
  - **created_at**: Fecha de creación del registro.
  - **updated_at**: Fecha de la última actualización del registro.

**Ejemplo**:  
| ZoneID | UserID  | Timestamp               | createdAt              | updatedAt              |
|--------|---------|-------------------------|------------------------|------------------------|
| Z001   | U1001   | 2023-08-25 08:00:00     | 2023-08-25 08:00:00    | 2023-08-25 08:00:00    |
| Z002   | U1002   | 2023-08-25 09:15:00     | 2023-08-25 09:15:00    | 2023-08-25 09:15:00    |
| Z003   | U1003   | 2023-08-25 10:30:00     | 2023-08-25 10:30:00    | 2023-08-25 10:30:00    |
| Z004   | U1004   | 2023-08-25 11:45:00     | 2023-08-25 11:45:00    | 2023-08-25 11:45:00    |
| Z005   | U1005   | 2023-08-25 12:00:00     | 2023-08-25 12:00:00    | 2023-08-25 12:00:00    |
| Z006   | U1006   | 2023-08-25 13:15:00     | 2023-08-25 13:15:00    | 2023-08-25 13:15:00    |
| Z007   | U1007   | 2023-08-25 14:30:00     | 2023-08-25 14:30:00    | 2023-08-25 14:30:00    |
| Z008   | U1008   | 2023-08-25 15:45:00     | 2023-08-25 15:45:00    | 2023-08-25 15:45:00    |
| Z009   | U1009   | 2023-08-25 16:00:00     | 2023-08-25 16:00:00    | 2023-08-25 16:00:00    |
| Z010   | U1010   | 2023-08-25 17:15:00     | 2023-08-25 17:15:00    | 2023-08-25 17:15:00    |

**Nota:** se tiene que considerar cual es el mejor metodo para almacenar zonas y realizar las busquedas mas eficientes, algunos metodos que podmos considerar son:
- Utilizar PostGIS con PostgreSQL
- Utilizar un servicio geoespacial externo
- Método de "ray casting" o "point-in-polygon"
- Método de cuadros delimitadores (bounding boxes)

Segun el metodo elegido se cambiara el diseño de la tabla.

### Tabla: UserLocationHistory
Guardara el historial de ubicaciones de los usuarios. Solo cuando cambie de ZoneID se guardara un nuevo registro. (trigger)
- **Atributos**:
  - **UserID** (PK): Identificador único para cada usuario.
  - **Timestamp** (SK): Fecha y hora exactas en que se guardó la zona para el usuario.
  - **ZoneID**: Identificador único para cada zona tomado de la tabla `ZonesConfig`.
  - **created_at**: Fecha de creación del registro.
  - **updated_at**: Fecha de la última actualización del registro.

**Ejemplo**:  
| UserID | Timestamp               | ZoneID |
|--------|-------------------------|--------|
| U1001  | 2023-08-25 08:00:00     | Z001   |
| U1002  | 2023-08-25 09:15:00     | Z002   |
| U1003  | 2023-08-25 10:30:00     | Z003   |
| U1004  | 2023-08-25 11:45:00     | Z004   |
| U1005  | 2023-08-25 12:00:00     | Z005   |
| U1006  | 2023-08-25 13:15:00     | Z006   |
| U1007  | 2023-08-25 14:30:00     | Z007   |
| U1008  | 2023-08-25 15:45:00     | Z008   |
| U1009  | 2023-08-25 16:00:00     | Z009   |
| U1010  | 2023-08-25 17:15:00     | Z010   |
## CREATE TABLE (pseudo-code):

```sql
CREATE TABLE ZonesConfig (
    ZoneID UUID PRIMARY KEY,
    ZoneName VARCHAR(255) NOT NULL,
    ZoneType VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE UserZones (
    ZoneID UUID REFERENCES ZonesConfig(ZoneID),
    UserID UUID NOT NULL,
    Timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (ZoneID, UserID)
);

CREATE TABLE UserLocationHistory (
    UserID UUID NOT NULL,
    Timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    ZoneID UUID REFERENCES ZonesConfig(ZoneID),
    PRIMARY KEY (UserID, Timestamp)
);

CREATE OR REPLACE FUNCTION insert_into_userlocationhistory_on_zoneid_change()
RETURNS TRIGGER AS $$
BEGIN
    -- Inserta un nuevo registro en UserLocationHistory solo si el ZoneID ha cambiado
    IF NEW.ZoneID <> OLD.ZoneID THEN
        INSERT INTO UserLocationHistory (UserID, Timestamp, ZoneID)
        VALUES (NEW.UserID, CURRENT_TIMESTAMP, OLD.ZoneID);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_insert_into_userlocationhistory
AFTER UPDATE OF ZoneID ON UserZones
FOR EACH ROW
WHEN (OLD.ZoneID IS DISTINCT FROM NEW.ZoneID)
EXECUTE FUNCTION insert_into_userlocationhistory_on_zoneid_change();

```

**Nota**: Hay que tener en cuenta que el uso intensivo de triggers puede afectar el rendimiento. Deberemos monitorear el rendimiento y buscar una solucion mejor si es necesario.

### Pendientes:

**Seguridad**: Asegurar que la información de ubicación esté segura, y considera cualquier regulación o ley de privacidad aplicable. 

**REDIS**: Analizar si es necesario utilizar REDIS para guardar la información de ubicación de los usuarios.

**Encriptar**??: datos sensibles como las coordenadas de ubicación, estén encriptados???