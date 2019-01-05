echo "Stopping...."
eval $(docker-machine env)
docker-compose down -v