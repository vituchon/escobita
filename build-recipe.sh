# local build
go build .
docker build -f Dockerfile -t escobita-image .
docker run --name escobita-container -d -p 9090:9090 --security-opt seccomp=unconfined escobita-image
// si quiero explorar docker exec -i -t escobita-container bash

# deploy to docker hub

docker login --username vitus
docker tag escobita-image vitus/escobita-image:1.0.0
docker push vitus/escobita-image:1.0.0
