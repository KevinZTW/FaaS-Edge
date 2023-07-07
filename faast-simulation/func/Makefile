.PONY: redeploy
redeploy:
	docker-compose build
	docker-compose push
	kubectl rollout restart deploy func