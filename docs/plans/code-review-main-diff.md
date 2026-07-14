# Code Review: diff `main..HEAD` en better-auth-sdk-go

**Fecha:** 2026-07-14
**Rama revisada:** `add-plugins` (HEAD)
**Base:** `main`
**Autor del review:** Kimi Code CLI

---

## Resumen ejecutivo

El PR refactoriza el SDK desde un cliente monolítico hacia una arquitectura basada en plugins, añade documentación para agentes (`AGENTS.md`) y para desarrolladores (`DEVELOPMENT.md`), elimina el directorio `claude/` (que parecía un artefacto de una sesión anterior) e introduce tres plugins iniciales: `admin`, `session` y `tenancy`.

La arquitectura es limpia, la API es coherente y el código compila y pasa los tests. Sin embargo, hay varios problemas que deberían resolverse antes de mergear: documentación que no compila, un posible `panic` en runtime, una data race y una API de errores tipificados que no se usa completamente.

---

## Checks ejecutados

| Comando | Resultado |
|---|---|
| `go build ./...` | ✅ |
| `go vet ./...` | ✅ |
| `go test ./...` | ✅ |
| `go test -race -cover ./...` | ✅ (solo `tenancy` tiene tests; cobertura global baja) |

---

## Hallazgos por severidad

### 🔴 Críticos / bloqueantes para merge

#### 1. README menciona plugins que no existen

**Archivo:** `README.md`

La sección *Plugins* importa y usa un plugin `phonenumber`:

```go
phone := phonenumber.New(client)
phone.SendOTP(ctx, "+34600000000")
```

La tabla de plugins también incluye `qrauth`, `expopasskey`, `checkphone` y `googlemapsproxy`.

En `plugins/` solo existen realmente:

- `admin/`
- `session/`
- `tenancy/`

**Impacto:** el código de ejemplo no compila y la documentación es engañosa.

#### 2. Ejemplo de middleware no compila

**Archivo:** `README.md`

```go
func authMiddleware(client *betterauth.Client) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // ...
            sessionData, err := sess.Get(r.Context()) // ❌ 'sess' no está declarado
        })
    }
}
```

**Impacto:** el ejemplo copiable está roto.

#### 3. `Do` hace panic si `SessionToken` es `nil`

**Archivo:** `client.go:73`

```go
cookie := c.SessionToken.Cookie
```

`NewClient` acepta `sessionToken *SessionToken`, por lo que `nil` es un uso válido (por ejemplo, en tests o cuando solo se usa bearer token). No se comprueba `c.SessionToken != nil` antes de acceder a `.Cookie`.

**Impacto:** panic en runtime.

#### 4. Data race en `SetBearerToken` + `Do`

**Archivo:** `client.go`

`SetBearerToken` escribe `c.bearerToken` mientras que `Do` lo lee, todo sin sincronización. El SDK recomienda reutilizar el cliente, por lo que el acceso concurrente es esperable.

**Impacto:** data race. Se puede reproducir fácilmente con un test que llame a `SetBearerToken` y a `Do` concurrentemente bajo `go test -race`.

---

### 🟡 Medios

#### 5. `Do` no mapea errores de red/timeout a los tipos propios

**Archivo:** `client.go:86`

```go
return fmt.Errorf("failed to perform request: %w", err)
```

Los errores de red o timeout se envuelven en un error genérico de Go, no en `*betterauth.Error` con tipo `network` o `timeout`. Como consecuencia, los helpers `IsNetworkError` / `IsTimeoutError` nunca devolverán `true`.

**Impacto:** la API de errores tipificados pierde utilidad.

#### 6. Construcción de URL frágil

**Archivo:** `client.go:62`

```go
url := c.config.BaseURL + c.config.BasePath + path
```

Si `BaseURL` termina en `/` o si `BasePath`/`path` se solapan, pueden generarse barras dobles (`//`). Se recomienda usar `url.JoinPath` o normalizar con `strings.TrimSuffix` / `strings.TrimPrefix`.

#### 7. Inconsistencia en validaciones de inputs

`AGENTS.md` y `DEVELOPMENT.md` piden validar inputs en el cliente antes de enviar la petición:

> Validate required inputs before the request and return `betterauth.NewError(betterauth.ErrorTypeValidation, ...)`.

`session.Verify` cumple con esto, pero `admin.CreateUser`, `tenancy.Organization.Create`, `tenancy.Team.Create` y otros no validan campos requeridos.

#### 8. Typo en godoc exportado

**Archivo:** `client.go:49`

```go
// Do perform an HTTP request with proper headers and error handling.
```

Debería ser `// Do performs ...`.

#### 9. `DEVELOPMENT.md` referencia archivos inexistentes

El documento menciona:

> See `plugins/phonenumber/phonenumber_test.go` for the pattern.

Ese archivo no existe.

También se mencionan comentarios `// ponytail:` para marcar rutas no verificadas, pero no hay ninguno en el código actual, a pesar de que el README dice que algunos endpoints son inferidos.

#### 10. Plugins `admin` y `session` carecen de tests

Solo `tenancy` tiene tests (`tenancy_test.go`). La cobertura global es prácticamente nula fuera de `tenancy`.

---

### 🟢 Menores / estilo

#### 11. `Config.Debug` no está cableado

`AGENTS.md` lo admite como *future concern*, pero exponer un flag público sin efecto puede confundir a los usuarios.

#### 12. `SetTimeout` muta el `http.Client` compartido

Si el usuario proporciona su propio `HTTPClient`, `SetTimeout` modifica el timeout de ese cliente externo. Esto puede ser sorprendente; debería documentarse o crear una copia interna.

#### 13. `parseErrorResponse` usa el cuerpo crudo si el JSON falla

**Archivo:** `errors.go:161`

```go
Message: string(body),
```

Si el servidor devuelve HTML (por ejemplo, un error 502 de un proxy), el mensaje de error puede ser muy grande. Sería más seguro truncarlo o usar un mensaje genérico.

#### 14. Uso de `map[string]interface{}` en `admin.ListUsers`

Contradice la guía de estilo de preferir structs tipados. Aunque en este caso se puede justificar por la flexibilidad del query, convendría definir un struct `ListUsersQuery`.

#### 15. `Permission.Check` recibe `DenyInput`

Funciona, pero el nombre del tipo es confuso para una operación de comprobación de permisos. Considerar un tipo neutral como `CheckInput`.

---

## Recomendaciones

### Mínimo viable antes de mergear

1. Corregir `README.md`: eliminar referencias a plugins inexistentes y arreglar el ejemplo de middleware.
2. Corregir `DEVELOPMENT.md`: quitar la referencia a `plugins/phonenumber/phonenumber_test.go` y aclarar si los comentarios `// ponytail:` aplican o no.
3. Proteger `Do` contra `SessionToken == nil`.
4. Añadir sincronización (`sync.RWMutex`) en `SetBearerToken` y en la lectura de `bearerToken` en `Do`.
5. Corregir el typo `Do perform` → `Do performs`.

### Mejoras recomendadas adicionales

- Añadir tests mínimos para `session` y `admin` siguiendo el patrón de `tenancy_test.go`.
- Añadir validaciones de campos requeridos en `admin` y `tenancy`.
- Normalizar la construcción de URL con `url.JoinPath`.
- Mapear errores de red/timeout a `ErrorTypeNetwork` / `ErrorTypeTimeout`.
- Limitar o sanitizar el mensaje de error cuando el cuerpo no es JSON.

---

## Notas sobre el diff

- La eliminación del directorio `claude/` parece correcta: en `main` todo el código relevante ya estaba en la raíz y `claude/` contenía duplicados y artefactos.
- El refactor a plugins es consistente con la filosofía de Better Auth y mejora la testabilidad.
- La arquitectura `Requester` + `Do` es adecuada y permite plugins de terceros sin modificar el core.
