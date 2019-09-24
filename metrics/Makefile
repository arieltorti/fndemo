NAME ?= fn-metrics

.PHONY: start
start:
	docker swarm init || true
	docker stack deploy --compose-file=docker-compose.yml ${NAME}

.PHONY: stop
stop:
	docker stack rm ${NAME}
