## Task Management System


### Migrations
Goose is used [goose](https://github.com/pressly/goose) for migrations.
To run migrations


### Running the service
```
docker-compose up -d
```

## Setting up environment variables

```
source env
```

### Running migrations
Use the following command to run migrations
```
goose -dir=db/sql  postgres $DSN up 
```