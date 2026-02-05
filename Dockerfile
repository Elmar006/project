FROM ubuntu:latest

WORKDIR /app

COPY scheduler .
COPY web ./web/

EXPOSE 7540

CMD ["./scheduler"]