#!/bin/bash

# Configuración de variables
BUCKET_NAME="test"
AWS_ACCESS_KEY_ID="test"
AWS_SECRET_ACCESS_KEY="test"
AWS_REGION="us-east-1"
LOCALSTACK_ENDPOINT="http://localhost:4566"
CONTAINER_NAME="localstack_main"

function start() {
  echo "Iniciando LocalStack con Docker..."
  docker run -d \
    --name $CONTAINER_NAME \
    -p 4566:4566 \
    -e SERVICES=s3 \
    -e DEFAULT_REGION=$AWS_REGION \
    -e DEBUG=1 \
    localstack/localstack

  echo "Esperando a que LocalStack esté listo..."
  until curl -s $LOCALSTACK_ENDPOINT/_localstack/health | grep -q '"s3": "available"'; do
    sleep 2
  done

  echo "Configurando credenciales para AWS CLI (perfil localstack)..."
  aws configure set aws_access_key_id $AWS_ACCESS_KEY_ID --profile localstack
  aws configure set aws_secret_access_key $AWS_SECRET_ACCESS_KEY --profile localstack
  aws configure set region $AWS_REGION --profile localstack

  echo "Creando el bucket: $BUCKET_NAME"
  aws --endpoint-url=$LOCALSTACK_ENDPOINT s3 mb s3://$BUCKET_NAME --profile localstack

  echo "Haciendo el bucket público (permitiendo lectura pública)..."
  aws --endpoint-url=$LOCALSTACK_ENDPOINT s3api put-bucket-policy --bucket $BUCKET_NAME --policy '{
      "Version": "2012-10-17",
      "Statement": [
          {
              "Sid": "PublicRead",
              "Effect": "Allow",
              "Principal": "*",
              "Action": ["s3:GetObject"],
              "Resource": ["arn:aws:s3:::'$BUCKET_NAME'/*"]
          }
      ]
  }' --profile localstack

  aws --endpoint-url=$LOCALSTACK_ENDPOINT s3api put-public-access-block \
      --bucket $BUCKET_NAME \
      --public-access-block-configuration "BlockPublicAcls=false,IgnorePublicAcls=false,BlockPublicPolicy=false,RestrictPublicBuckets=false" \
      --profile localstack

  echo "¡LocalStack configurado!"
  echo "Bucket: $BUCKET_NAME"
  echo "Endpoint: $LOCALSTACK_ENDPOINT"
}

function stop() {
  echo "Deteniendo y eliminando el contenedor $CONTAINER_NAME y sus volúmenes..."
  docker stop $CONTAINER_NAME
  docker rm -v $CONTAINER_NAME
  echo "¡LocalStack detenido y eliminado!"
}

case "$1" in
  start)
    start
    ;;
  stop)
    stop
    ;;
  *)
    echo "Uso: $0 {start|stop}"
    exit 1
    ;;
esac
