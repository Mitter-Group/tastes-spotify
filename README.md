# Spotify


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


---

#### 1. Descripción 

Se conectará con Spotify para obtener datos musicales de los usuarios, como sus canciones recientes, artistas top, géneros. Con esta info analizaremos los gustos musicales del usuario a lo largo del tiempo.

Usaremos el API de Spotify:
https://developer.spotify.com/documentation/web-api

---
#### 2. Endpoints

##### Autenticación
- **/api/spotify/auth/start (GET)**
  - **Descripción**: Inicia el proceso de autenticación con Spotify, redirigiendo al usuario a la página de autenticación.
  - **Recibe**: N/A
  - **Devuelve**: Redirección a la página de autenticación de Spotify.
  - **Ejemplo**: `GET /api/spotify/auth/start`

- **/api/spotify/auth/callback (GET)**
  - **Descripción**: Endpoint al que Spotify redirige tras la autenticación. Aquí se procesan y almacenan los tokens.
  - **Recibe**: `code` (Código de autorización de Spotify).
  - **Devuelve**: Confirmación del estado de autenticación y posible token de acceso.
  - **Ejemplo**: `GET /api/spotify/auth/callback?code=ABC123`

- **/api/spotify/auth/refresh (POST)**
  - **Descripción**: Refresca el token de acceso utilizando el token de actualización.
  - **Recibe**: `userID`
  - **Devuelve**: Nuevo `accessToken`.
  - **Ejemplo**: `POST /api/spotify/auth/refresh { "userID": "user123" }`

##### Datos
- **/api/spotify/data/tracks (GET)**
  - **Descripción**: Obtiene las pistas más recientes o favoritas del usuario.
  - **Recibe**: `userID`
  - **Devuelve**: Lista de pistas.
  - **Ejemplo**: `GET /api/spotify/data/tracks?userID=user123`
  
- **/api/spotify/data/artists (GET)**
  - **Descripción**: Obtiene los artistas top del usuario.
  - **Recibe**: `userID`
  - **Devuelve**: Lista de artistas.
  - **Ejemplo**: `GET /api/spotify/data/artists?userID=user123`

- **/api/spotify/data/genres (GET)**
  - **Descripción**: Obtiene los géneros favoritos o más escuchados por el usuario.
  - **Recibe**: `userID`
  - **Devuelve**: Lista de géneros.
  - **Ejemplo**: `GET /api/spotify/data/genres?userID=user123`

- **/api/spotify/data/playlists (GET)**
  - **Descripción**: Obtiene las listas de reproducción del usuario.
  - **Recibe**: `userID`
  - **Devuelve**: Lista de listas de reproducción.
  - **Ejemplo**: `GET /api/spotify/data/playlists?userID=user123`

Spotify solo permite consulta TOP tracks y artist, debemos calcular de nuestro lado TOP tracks.

_Para analizar: solo podemos consultar las últimas 50 canciones, si queremos un historial más específico debemos guardar los datos en nuestra base de datos y refrescarlos cada cierto tiempo. Tomar en cuenta la quota de Spotify._


---

#### 3. Diseño de las Tablas


- **Tabla de Autenticación (Tokens)**
  
  
  **NOTA**: Esta tabla no la neceistamos en este servicio, debemos pasarla al servicio de usuarios. Investigar cognito.

  - **Descripción**: Almacena la información de autenticación de cada usuario.
  - **Datos Ejemplo**:

    | userID | accessToken | refreshToken | expirationTime | 
    |--------|-------------|--------------|----------------|
    | user1  | tok123      | ref123       | 1677812359     |
    | user2  | tok456      | ref456       | 1677898759     |



- **Tabla UserSpotifyData**
  - **Descripción**: Almacena la información de Spotify relacionada con el usuario.
  - **Datos Ejemplo**:

    | userID | dataType | data                                         | updatedAt              | createdAt              |
    |--------|----------|----------------------------------------------|------------------------|------------------------|
    | user1  | tracks   | ["song1", "song2", "song3"]                  | 2023-08-27T15:30:00Z   | 2023-08-26T14:20:00Z   |
    | user1  | artists  | ["artist1", "artist2"]                       | 2023-08-27T15:45:00Z   | 2023-08-26T14:25:00Z   |
    | user2  | genres   | ["rock", "pop"]                              | 2023-08-28T16:32:00Z   | 2023-08-27T15:22:00Z   |



---

#### 5. Nota

Este diseño es básicamente para empezar con un MVP. En el futuro, queremos agregar una tabla para guardar cómo cambian los gustos musicales del usuario. Para esto, usaremos streams de DynamoDB que tomaran esos cambios en tiempo real en la tabla de música.
