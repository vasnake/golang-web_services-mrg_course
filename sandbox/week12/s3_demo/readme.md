https://docs.min.io/docs/minio-docker-quickstart-guide.html
https://docs.min.io/docs/deploy-minio-on-docker-compose.html

docker run -p 9000:9000 \
  -e "MINIO_ACCESS_KEY=access_123" \
  -e "MINIO_SECRET_KEY=secret_123" \
  -v "$(pwd)"/data:/data \
  minio/minio server /data


https://raw.githubusercontent.com/minio/minio/master/docs/orchestration/docker-compose/docker-compose.yaml

https://docs.min.io/docs/distributed-minio-quickstart-guide
https://docs.min.io/docs/minio-erasure-code-quickstart-guide