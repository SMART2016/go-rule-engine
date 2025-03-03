FROM postgres:17

RUN apt update && apt install -y postgresql-17-cron

CMD ["postgres", "-c", "shared_preload_libraries=pg_cron"]
