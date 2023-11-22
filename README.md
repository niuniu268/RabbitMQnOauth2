# Oauth2 N RabbitMQ

```docker run -d --hostname my-rabbit --name some-rabbit -p 15672:15672 -p 5672:5672 rabbitmq:3-management ```` 

Get token 

```http://localhost:9096/token?grant_type=client_credentials&client_id=000000&client_secret=999999&scope=read ```
Pay attention to the token check in the rabbitMQ server (line 70)
