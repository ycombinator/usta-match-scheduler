## Running

1. Start API server.
```
go run main.go
```

API is served at http://localhost:8000/

## Routes

### Get organization matches

#### Endpoint
```
GET /organization/{id}/matches
```

#### Path parameters
| Name | Type    | Description                 | Required? | Default value |
|------|---------|-----------------------------|-----------|---------------|
| `id` | Integer | USTA NorCal Organization ID | Yes       | -             |

### Query string parameters
| Name           | Type    | Description                                                                        | Required? | Default value |
|----------------|---------|------------------------------------------------------------------------------------|-----------|---------------|
| `is_scheduled` | Boolean | Filter results to include only scheduled (`true`) or unscheduled (`false`) matches | No        | -             |
| `location`     | String  | Filter results to include only `home` or `away` matches                            | No        | -             |
