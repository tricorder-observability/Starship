# DAO

[DAO/Data Access Object](https://en.wikipedia.org/wiki/Data_access_object)
are ORM types for easier management of database data table.

Include types integrating with GORM and corresponding APIs for serializing and
accessing sqlite DB data tables.

**NOTE**: There is no DAO layer for PostgreSQL username & password, as the postgreSQL
is accessed through username & password directly, no need to produce new API key
like Grafana.

Grafana API keys seems are stored in its own PVC.

Grafana API keys can only be created once, later creation request with the same username
can only have 1 key live inside Grafana. Later requests to obtain API key will return
already exists, but wont' return the key.

TODO:
We can instead always fetch new API Key from Grafana with the username+password.
And remove the old key if it already exists.
