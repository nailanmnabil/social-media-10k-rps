# marketplace-10k-rps

Make new migration script

```migrate create -ext sql -dir migration -seq init```

Migrate database

up local : ```../migrate -database postgres://postgres:root@localhost:5432/satu?sslmode=disable -path ./db/migrations up```

up : ```../migrate -database postgres://postgres:ohN6Nei0ugiRena5@project-sprint-postgres.cavsdeuj9ixh.ap-southeast-1.rds.amazonaws.com:5432/postgres?sslmode=verify-full?sslrootcert=ap-southeast-1-bundle.pem  -path ./db/migrations up```

DB_NAME=postgres
DB_PORT=5432
DB_HOST=project-sprint-postgres.cavsdeuj9ixh.ap-southeast-1.rds.amazonaws.com
DB_USERNAME=postgres
DB_PASSWORD=ohN6Nei0ugiRena5

down : ```migrate -database postgres://postgres:root@localhost:5432/marketplace?sslmode=disable -path ./db/migrations down```