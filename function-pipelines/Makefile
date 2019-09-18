.PHONY: start
start:
	docker swarm init || true
	docker stack deploy --compose-file=docker-compose.yml function-pipelines

.PHONY: stop
stop:
	docker stack rm function-pipelines
