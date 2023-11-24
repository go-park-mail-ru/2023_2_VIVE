# Dockerfile for postgres with russian dicts for full text search 
FROM postgres:16.0-alpine

# Опционально: если нужно выполнить дополнительные команды для настройки словарей
# RUN wget https://bitbucket.org/Shaman_Alex/russian-dictionary-hunspell/downloads/ru_RU_UTF-8_20131101.zip
# RUN uzip ru_RU_UTF-8_20131101.zip
# RUN DICT=/usr/local/share/postgresql/tsearch_data
# RUN iconv -o $DICT/russian.affix ru_RU.aff
# RUN iconv -o $DICT/russian.dict ru_RU.dic

# Устанавливаем рабочую директорию в контейнере
WORKDIR /usr/local/share/postgresql/tsearch_data

# Копируем файлы словарей из локальной директории 'dicts' в рабочую директорию в контейнере
COPY dicts/ ./

# Настраиваем права доступа (если необходимо)
# RUN chmod -R 755 /usr/share/postgresql/tsearch_data

# Можно добавить скрипты для настройки PostgreSQL, чтобы использовать новые словари
# Например, можно добавить SQL-скрипты в директорию /docker-entrypoint-initdb.d/, которая автоматически выполняется при инициализации базы данных
# COPY init_scripts/ /docker-entrypoint-initdb.d/

