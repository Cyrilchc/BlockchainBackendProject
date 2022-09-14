FROM postgres:14.5-alpine

ENV POSTGRES_PASSWORD=cyril
ENV POSTGRES_USER=postgres
ENV POSTGRES_DB=blockchainbackend

WORKDIR /app

# Éxécute le script au lancement
COPY script.sql /docker-entrypoint-initdb.d/

EXPOSE 5432 5432