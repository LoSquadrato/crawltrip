
start_or_run () {
    docker inspect crawltrip_mongodb > /dev/null 2>&1

    if [ $? -eq 0 ]; then
        echo "Starting Crawltrip MongoDB container..."
        docker start crawltrip_mongodb
    else
        echo "Crawltrip MongoDB container not found, creating a new one..."
        docker run -d --name crawltrip_mongodb -p 27017:27017 -d mongodb/mongodb-community-server:latest
        
    fi
}

case "$1" in
    start)
        start_or_run
        ;;
    stop)
        echo "Stopping Crawltrip MongoDB container..."
        docker stop crawltrip_mongodb
        ;;
    logs)
        echo "Fetching logs for Crawltrip MongoDB container..."
        docker logs -f crawltrip_mongodb
        ;;
    status)
        docker ps -a --filter "name=$CONTAINER_NAME"
        ;;
    *)
        echo "Usage: $0 {start|stop|logs|status}"
        exit 1
esac