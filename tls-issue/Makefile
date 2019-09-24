NAME ?= tls-demo
NAME_FULL ?= tls-demo-full

.PHONY: start
start:
	docker swarm init || true
	docker stack deploy --compose-file=docker-compose.yml ${NAME}

.PHONY: start-deploy
start-deploy:
	docker swarm init || true
	docker stack deploy --compose-file=docker-compose.full.yml ${NAME_FULL}

.PHONY: stop
stop:
	docker stack rm ${NAME}

.PHONY: stop-deploy
stop-deploy:
	docker stack rm ${NAME_FULL}
