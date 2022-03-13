echo "Building image"
docker build -t forum .
echo "Running image"
docker run -d -p 8080:8080 --name forum forum
echo "Images list"
docker images
echo "Containers list"
docker container ls