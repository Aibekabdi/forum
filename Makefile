build:
	docker build -t forum .
run:
	docker run -p 8080:8080 forum
delete:
	docker image prune forum
	docker container prune forum