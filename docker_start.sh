docker-compose down || true
docker-compose pull
docker-compose up --remove-orphans --build