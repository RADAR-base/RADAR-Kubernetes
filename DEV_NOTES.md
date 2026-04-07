# Nessie

## Create database

```sql
create database nessie;
create user nessie with password '';
GRANT ALL ON SCHEMA public TO nessie;
ALTER SCHEMA public OWNER TO nessie;
```

