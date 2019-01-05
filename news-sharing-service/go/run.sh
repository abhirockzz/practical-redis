echo "Starting...."
eval $(docker-machine env)
docker-compose down -v
docker-compose up --build