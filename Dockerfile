FROM golang

WORKDIR /app

COPY . /app

LABEL project="ASCII-ART-WEB" \
      authors="Aibek_kz, satybalding" \
      description="Forum" \
      link="https://git.01.alem.school/satybalding/forum"

RUN go build -o forum ./cmd/

EXPOSE 8080
CMD ["/app/forum"]
# echo "Building image"
# docker build -t ascii-art-web .
# echo "Running image"
# docker run -d -p 8080:8080 --name web ascii-art-web
# echo "Images list"
# docker images
# echo "Containers list"
# docker container ls