start_or_run () {
    docker inspect crawltrip_nats > /dev/null 2>&1

    if [ $? -eq 0 ]; then
        echo "Starting Crawltrip NATS container..."
        docker start crawltrip_nats
    else
        echo "Crawltrip NATS container not found, creating a new one..."
        docker run -d --name crawltrip_nats -p 4222:4222 -p 8222:8222 -v crawltrip_nats_vol:/data nats:2.10-alpine --jetstream --store_dir /data -m 8222
    fi
}

case "$1" in
    start)
        start_or_run
        ;;
    stop)
        echo "Stopping Crawltrip NATS container..."
        docker stop crawltrip_nats
        ;;
    logs)
        echo "Fetching logs for Crawltrip NATS container..."
        docker logs -f crawltrip_nats
        ;;
    status)
        docker ps -a --filter "name=$CONTAINER_NAME"
        ;;
    *)
        echo "Usage: $0 {start|stop|logs|status}"
        exit 1
esac