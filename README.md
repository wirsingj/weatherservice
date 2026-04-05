# Weather Service

Go HTTP service that fetches today's forecast from the National Weather Service API.

Project is broken into 3 parts:
Main: Handles initialization of weatherClient and HTTP server,  and handles graceful shutdown.
WeatherClient: Handles fetching weather data from the NWS API.
API: Handles incoming HTTP requests, validates query parameters, and returns weather data.

A swagger document is also included to provide the API specification.

## Run

```bash
go run .
```

Server starts on `:8080`.

## Run (Precompiled Binary)

Prebuilt binaries are in `build/`:
- `build/weatherservice-windows-amd64.exe`
- `build/weatherservice-darwin-amd64`
- `build/weatherservice-darwin-arm64`

Windows (PowerShell):
```powershell
.\build\weatherservice-windows-amd64.exe
```

macOS (Intel):
```bash
chmod +x ./build/weatherservice-darwin-amd64
./build/weatherservice-darwin-amd64
```

macOS (Apple Silicon):
```bash
chmod +x ./build/weatherservice-darwin-arm64
./build/weatherservice-darwin-arm64
```

## Endpoint

`GET /weather?lat={latitude}&lon={longitude}`

Example:

```bash
curl "http://localhost:8080/weather?lat=39.7392&lon=-104.9903"
```

Response example:

```json
{
  "short_forecast": "Partly Sunny",
  "temperature": 72,
  "temperature_unit": "F",
  "temperature_characterization": "moderate"
}
```
