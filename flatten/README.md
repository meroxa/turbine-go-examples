## Flatten

This example Turbine app leverages the `transforms` package to flatten a nested JSON record. This can be useful when
going from a resource that supports nested objects (e.g. a Document Store such as MongoDB) to a relational database
(such as Postgres or Redshift).

### Input Record
```json
{
  "id": 1,
  "user": {
    "id": 100,
    "name": "alice",
    "email": "alice@example.com"
  },
  "actions": ["register", "purchase"]
}
```

### Output Record
```json
{
    "actions.0": "register",
    "actions.1": "purchase",
    "id": 1,
    "user.email": "alice@example.com",
    "user.id": 100,
    "user.name": "alice"
}
```