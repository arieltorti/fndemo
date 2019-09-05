.PHONY: start
start:
	docker swarm init || true
	docker stack deploy --compose-file=docker-compose.yml tls-demo

.PHONY: start-deploy
start-deploy:
	docker swarm init || true
	docker stack deploy --compose-file=docker-compose.full.yml tls-demo-full

.PHONY: stop
stop:
	docker stack rm tls-demo

.PHONY: stop-deploy
stop-deploy:
	docker stack rm tls-demo-full
