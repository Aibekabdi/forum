FROM golang

WORKDIR /app

COPY . /app

LABEL project="Forum" \
      authors="Aibek_kz, satybalding" \
      description="Forum" \
      link="https://git.01.alem.school/satybalding/forum"

RUN go build -o forum ./cmd/

EXPOSE 8080
CMD ["/app/forum"]