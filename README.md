# todo_sample

## 起動方法

`git clone git@github.com:sonarAIT/todo_sample.git`で clone する

`cd todo_sample`で移動．

`docker network create todo_network`を実行．

`cd vue && npm install`を実行．

2 つターミナルを立ち上げ，`todo_sample`内で`docker-compose up`を

`vue`内で`npm run serve`をそれぞれ実行することで起動できる．

（以後，`docker-compose up`と`npm run serve`だけでよい．）

データベースをリセットしたくなった時は，docker のサーバを stop した上で，

```
docker container rm `docker container ls -aq`
```

を実行する．

document 内の study.md が教科書に該当する．homework.md は問題集だ．

ソースコードにはコメントも書いてある．色んなことを学んで欲しい．
