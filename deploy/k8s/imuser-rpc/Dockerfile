FROM showurl/zerobase

WORKDIR /app

COPY ./bin /app/zeroservice

RUN chmod +x /app/zeroservice && mkdir /app/etc

CMD ["/app/zeroservice"]

