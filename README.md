# go-stress-test

1. Na pasta raiz do projeto

2. Compilar imagem docker
```shell
docker build -t stresstest .
```

3. Executar a aplicação
```shell
docker run stresstest --url=http://google.com --requests=1000 --concurrency=10
```
