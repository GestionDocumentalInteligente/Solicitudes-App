# Crear primer imagen

aws ecr get-login-password --region us-east-2 | docker login --username AWS --password-stdin 396608773260.dkr.ecr.us-east-2.amazonaws.com

docker build -t sg-ui:0.0.1 .

docker tag sg-ui:0.0.1 396608773260.dkr.ecr.us-east-2.amazonaws.com/sg-backend/sg-ui:0.0.1

docker push 396608773260.dkr.ecr.us-east-2.amazonaws.com/sg-backend/sg-ui:0.0.1
