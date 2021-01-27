FROM postgres

ENV POSTGRES_PASSWORD=password
COPY db_init.sql /docker-entrypoint-initdb.d/db_init.sql

EXPOSE 5432
