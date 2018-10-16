up: ## Docker stack deploy
	docker pull muninn/restdemo && docker stack deploy -c docker-compose.yml restdemo

down: ## Docker stack rm
	docker stack rm restdemo
