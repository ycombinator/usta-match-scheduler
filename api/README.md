## Running

1. Start API server.
```
go run main.go
```

API is served at http://localhost:8000/

## Routes

### Get USTA organization teams
Retrieves the list of teams for a specific USTA NorCal organization.

#### Endpoint
```
GET /usta/organization/{id}/teams
```

#### Path parameters
| Name | Type    | Description                 | Required? | Default value |
|------|---------|-----------------------------|-----------|---------------|
| `id` | Integer | USTA NorCal Organization ID | Yes       | -             |

### Query string parameters
| Name       | Type    | Description                                                                         | Required? | Default value |
|------------|---------|-------------------------------------------------------------------------------------|-----------|---------------|
| `upcoming` | Boolean | Filter results to include only teams whose season has not yet started               | No        | False         |


### Get USTA organization matches
Retrieves the list of matches for a specific USTA NorCal organization.

#### Endpoint
```
GET /usta/organization/{id}/matches
```

#### Path parameters
| Name | Type    | Description                 | Required? | Default value |
|------|---------|-----------------------------|-----------|---------------|
| `id` | Integer | USTA NorCal Organization ID | Yes       | -             |

### Query string parameters
| Name           | Type    | Description                                                                         | Required? | Default value |
|----------------|---------|-------------------------------------------------------------------------------------|-----------|---------------|
| `is_scheduled` | Boolean | Filter results to include only scheduled (`true`) or unscheduled (`false`) matches  | No        | -             |
| `location`     | String  | Filter results to include only `home` or `away` matches                             | No        | -             |
| `after`        | String  | Filter results to include only matches after this date (inclusive; ISO8601 format)  | No        | -             |
| `before`       | String  | Filter results to include only matches before this date (inclusive; ISO8601 format) | No        | -             |
