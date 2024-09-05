docker compose up --build

docker exec -it network_show-mynode-1 sh

python echo-client.py mynode 3000 1
