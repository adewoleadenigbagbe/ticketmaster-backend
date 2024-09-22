FROM gomicro/goose

WORKDIR /migrations
COPY ./shared/database/migrations .

#RUN chmod +x /migrations/entrypoint.sh

ENTRYPOINT ["/migrations/entrypoint.sh"]