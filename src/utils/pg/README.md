# Postgres access API

To start a postgres server:

```
docker run -it --rm --name=pg --env=POSTGRES_PASSWORD=passwd \
  --volume=/var/run/postgresql:/var/run/postgresql postgres
```
