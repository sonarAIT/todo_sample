# 30 分でわかれ Go と SQL

**30 分でやるの普通に無理だと思う**．爆速でやって 3 日ぐらいかな…

## はじめに

このドキュメントではバックエンドの初歩的な知識をなるべく迅速に身に付けてもらい，Web アプリ開発におけるバックエンドの必要最低限な能力を身に付けてもらうことを期待しています．頑張ってください．

## バックエンドエンジニアが何をするか

**ぶっちゃけバックエンドの仕事なんざ無限種類あるのだが**，本ドキュメントにおける仕事は，Web アプリのサーバーサイド側のプログラムを書くことである．（データベースの設計もあるが，今回はこれは保留する．）

サーバーサイドのプログラムが何をするかというと，

```
・クライアントからのリクエストを受け取る．
・リクエストの種別(GET, POST, PUT, DELETE)を識別する．
・リクエストの内容(body, query)を解読する．
・それに合わせたデータベース等の処理を行う
・処理結果を送信できるように変換する．
・クライアントに送信する．
```

という一連の流れを行う．**が，データベース等の処理以外はほとんどコピペで済む．ただデータベース等の処理を書くのにどえらい時間を食う**．そういう仕事を本ドキュメントでは扱う．（外の世界に出たら実際にサーバを管理したり他にも仕事がある．でも，**とりあえずハッカソンとかやるのにその知識はいらない．**）

## 知っておいて欲しい概念

### Go

今回サーバーサイドのプログラムを書くのに使う言語．えらくシンプル．Go じゃないとできないということはない．（他にも選択肢がある．）

### リレーショナルデータベース

そういうデータベースの種類がある．めっちゃよく使われる．（大体データベースと言ったらこれ．）

### SQL

リレーショナルデータベースを操作するのに使用する言語．このドキュメントを勉強目的で読む人は大体 C 言語とかのプログラミング言語しかやったことがないと思うけど，**SQL はそもそもプログラミング言語ですらない**．なので初学者からしてみれば新概念を大量に詰め込まれることになる．頑張れ．

### http

Vue.js 編である程度書いたのでコピペしながら書く．サーバとクライアントの間で通信するためのプロトコル（通信規程）．

クライアントから受信する http のリクエストには 4 つの種類がある．

```
POST: データの送信
GET: データの受信
PUT: データの更新（で使われることが多い）
DELETE: データの削除
```

この 4 つを使い分けることによって，同じ URL に対して 4 種類の操作を実現することが可能である．

なので，例えば，http://hoge/tasks という 1 つの URL には

```
・タスクを追加する
・全タスクを取得する
・タスクを更新する（タスク名の変更など）
・タスクを削除する
```

という 4 種類の処理を割り当てることができる．

### body

http のリクエストにはデータを添付することができる．それが body．例えば，タスクをサーバに登録するために，body にタスクの情報を入れて送信したりする．**GET だと送れないので注意**（厳密には curl を使えば送れる．が，axios などでは基本的には送れないので，多分非推奨なんだろう．多分．）

余談だが，当サンプルでは新しく登録する一つのタスクの情報ではなく，存在する全てのタスクをいちいちサーバに送信している．（なぜそんな無駄な仕様になっているかというと，それを直す問題を出題しているから．last-homework.vue を参照．）

### クエリパラメータ

GET でもデータを送りたいことがある．その時はクエリパラメータを使用する．

```
http://hoge/tasks?user=sonarAIT&hoge=fuga
```

という感じに書けば，user は sonarAIT で hoge は fuga というデータが送信できる．

ただ，**外から丸見えなので注意．**

### curl

URL と POST, GET, PUT, DELETE と送る body を指定して，実際にリクエストをターミナル上から送信できるコマンド．使い方は後ほど．

## Go の書き方

**なか゛いたひ゛か゛ はし゛まる…**

### 文法

**ググれば出る**．配列とか宣言の仕方を忘れたらその都度「Golang 配列」とかでググればいいだけである．というかサンプルコードを眺めていればなんとなくわかってもらえると思う．

本ドキュメントでは，ググっても出ない単純な Web アプリケーションの実装方法を説明していく．

### クライアントからのリクエストを受け取る

まず main.go に注目して欲しい．import はどうでもいい．main 関数が記述されている．main 関数は Go 言語で一番はじめに実行される関数である．

```
func main() {
	http.HandleFunc("/tasks", apifuncs.TasksFunc)
	http.HandleFunc("/labels", apifuncs.LabelsFunc)

	http.ListenAndServe(":80", nil)
}
```

これは，要するに，

「クライアントに http://(サーバの IP アドレス)/tasks にアクセスされたら apifuncs の TasksFunc を実行します．」

「クライアントに http://(サーバの IP アドレス)/labels にアクセスされたら apifuncs の LabelsFunc を実行します．」

「80 番ポートで受け取ります」

という意味である．

80 番ポートに関してはどうでも良くて，ここで覚えて欲しいのは，**クライアントがリクエストできる URL の種類を増やしたかったら http.HandleFunc と対応する関数を増やせばいい**ということである．

例えば，このアプリケーションにログイン機能を実装して，ログインの http リクエストを受け取ることになった場合，

```
http.HandleFunc("/login", apifuncs.LoginFunc)
```

のように書けばよい．

次項からは，その実行される関数の書き方について記述していく．

### リクエストの種別(GET, POST, PUT, DELETE)を識別する．

apiFuncs のファイル（どっちでもいい）の関数を見てみよう．

```
w.Header().Set("Access-Control-Allow-Origin", "*")                       // Allow any access.
w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE") // Allowed methods.
w.Header().Set("Access-Control-Allow-Headers", "*")

r.Header.Set("Content-Type", "application/json")
```

**これは　もう　覚えろ．**（ちょっとだけ説明すると，上 3 つは受け取るリクエストのセキュリティ的な設定である．これを見ると**全種類のリクエストを受け取る**というめちゃくちゃガバガバなセキュリティになっていることがわかる．しかし，実際にデプロイしないハッカソンとかでは**どうでもいい**）

（`r.Header.Set("Content-Type", "application/json")`は，返信を JSON 形式で行いますという意味である．この Web アプリケーションでは，絶対に JSON で返信するので，1 番上にこのことを記述している．場合によってはこれを返信する直前に書く．）

問題は次だ．

```

if r.Method == http.MethodGet {
    //省略
} else if r.Method == http.MethodPost {
    //省略
}

```

**なんとわかりやすいのだろう**．上記で散々述べた GET,POST,PUT,DELETE はここで分岐させて処理していることがわかる．

これによって，同じ URL でも GET,POST,PUT,DELETE で全く異なった処理を実現することが可能になる．

識別できたので，リクエストの中身を読み取る．

### body を読み取る．

クライアントがサーバにリクエストを送る際， body とクエリパラメータの二種類の送り方があることは先ほど述べた通りである．

とりあえず，body の受け取り方をやる．apifuncs/tasks.go の POST を参照する．

/tasks の POST では「クライアントからタスク一覧を受け取り，受け取ったタスク一覧に合わせてデータベースを更新する」という処理を行っている．つまり，以下の処理は「クライアントから受け取ったタスク一覧を解読する」という処理だ．

```
jsonBytes, err := ioutil.ReadAll(r.Body)
if err != nil {
	w.WriteHeader(http.StatusServiceUnavailable)
	fmt.Fprintln(w, `{"status":"Unavailable"}`)
	fmt.Println("Can't catch Tasks(io error)", err)
	return
}

var recTasks []dbctl.Task
if err := json.Unmarshal(jsonBytes, &recTasks); err != nil {
	w.WriteHeader(http.StatusServiceUnavailable)
	fmt.Fprintln(w, `{"status":"Unavailable"}`)
	fmt.Println("Can't catch Tasks(JSON Unmarshal error)", err)
	return
}
```

長くて申し訳ないのだが，dbctl/tasks.go のこの部分も参照する．

```
// Task はタスク1つ分の構造体
type Task struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SubmitTime  string `json:"submitTime"`
	Label       int    `json:"label"`
}
```

とりあえず，この処理について説明していく．

何はともあれ，リクエストを人が読める形にしないといけない．その第一段階がこれである．

```
jsonBytes, err := ioutil.ReadAll(r.Body)
```

送られてきたリクエストを Byte の列にする．（まだ読めない．）

ところで，この処理は失敗することがある．失敗した場合，エラー処理を行わないといけない．なので，

```
if err != nil {
	w.WriteHeader(http.StatusServiceUnavailable)
	fmt.Fprintln(w, `{"status":"Unavailable"}`)
	fmt.Println("Can't catch Tasks(io error)", err)
	return
}
```

もしエラーが null でなければ，失敗したことをクライアントに告げ，ログにも失敗したことを記述し，処理を終了する．（ちなみに nil は null という意味である．なんで nil なんだろう．）

`if err != nil`と書いてあるところは全てエラー処理である．

次に，読み取ったリクエストを受け取るための変数を宣言する．/tasks の POST には，タスクの配列（タスク一覧）が送られてくる．だから，

```
var recTasks []dbctl.Task
```

タスクの配列を受け取るために，dbctl に定義されている Task の配列を宣言している．

次に，recTasks にこれを代入する．

```
if err := json.Unmarshal(jsonBytes, &recTasks); err != nil {
```

ややこしい書き方だが，要するに，jsonBytes を読める形にしてから recTasks に代入，もしエラーが起きたらエラー処理をするという処理が書かれている．

実際，if 文の中身を見てみるとエラー処理が行われていることがわかる．（上記とほぼ同じである．）

ところで，関数が error を返すか（つまり，err := しないといけないか）はその関数のドキュメントを見ればわかるが，大体 VSCode の go 拡張機能 が教えてくれるのであまり気にする必要はない．

まとめると，

・`jsonBytes, err := ioutil.ReadAll(r.Body)`で body をバイト列に変換．

・`var recTasks []dbctl.Task`で body を受け取るための配列変数を宣言．

・`if err := json.Unmarshal(jsonBytes, &recTasks); err != nil {`でバイト列を変換し変数に代入．

・それぞれ適所でエラー処理を行う

という一連の処理が行われている．

### クエリパラメータを受け取る

本サンプルにはクエリパラメータを受け取って行われる処理はない．

だが，受け取り方は body と比べて簡単である．

```
query := r.URL.Query()
log.Print(query.Get("user"))
```

先ほど，`http://hoge/tasks?user=sonarAIT&hoge=fuga`という URL の例を記載した．この処理が実行されれば，sonarAIT という値が表示されるはずだ．

### リクエストに合わせたデータベース等の処理を行う

/tasks の POST の処理はややこしいので，/tasks の GET を見てみる．(apifuncs/tasks.go)

/tasks の GET では，「クライアントにデータベースのタスク一覧を送信する」という処理を行っている．（必要がないので，body やクエリパラメータは受け取っていない．）

今回行うデータベースの処理は，「データベースからタスク一覧を取得する」である．

```
Tasks, err := dbctl.GetTasks()
if err != nil {
	w.WriteHeader(http.StatusServiceUnavailable)
	fmt.Fprintln(w, `{"status":"Unavailable"}`)
	fmt.Println("database error(GetTasks)", err)
	return
}
```

1 番上の行以外は先ほどやったエラー処理である．（全く同じ．）

:=とは，宣言と代入を同時に行うという意味である．dbctl.GetTasks()の戻り値の型は Task の配列と error であることが dbctl に宣言されているので，型をこちらでわざわざ書かなくても自動的に型推論が行われる．

```
var Tasks []dbctl.task
var err error
Tasks, err = dbctl.GetTasks()
```

`Tasks, err := dbctl.GetTasks()`は，これと全く同じ意味である．

さて，上記によれば dbctl の GetTasks()の結果を Tasks に代入しているらしい．なので，ここから飛んで，dbctl/tasks.go の GetTasks を参照する．

```
rows, err := db.Query("select id, name, description, submitTime, label from tasks")
if err != nil {
	pc, file, line, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	log.Printf(errFormat, err, f.Name(), file, line)
	return nil, err
}
defer rows.Close()
```

真ん中の`if err != nil`の中身がいつもとちょっと違うが，結局はエラーをログに書いて，エラーを呼び出し元（つまり apifuncs.TaskFunc）に返すという処理である．（クライアントに通達するのは呼び出し元の仕事なので，ここでは行わない．）

つまり，本質的には，

```
rows, err := db.Query("select id, name, description, submitTime, label from tasks")
defer rows.Close()
```

こう．

db.Query はデータベースからデータを取得する際に書く処理である．引数にはみんな大好き SQL が記述されている．（後述．**こんな隙間で済むほど単純ではない．**）

ちなみに，データベースを更新する際には，db.Exec()という関数を実行する．（同ファイルの InsertTasks や DeleteTasks などを参照されたい．）

rows に問い合わせの結果を，err にエラーを代入する．

`defer rows.Close()`は，後で(つまり，dbctl.GetTasks が終わる際に)rows.Close()を行うという意味である．**defer rows.Close は忘れてはいけない**．忘れるとメモリがリークする．

問い合わせが終わったので，それを変数に取り出していく．

```
var Tasks []Task

for rows.Next() {
	var id int
	var name string
	var description string
	var submitTime string
	var label int
	rows.Scan(&id, &name, &description, &submitTime, &label)
	Tasks = append(Tasks, Task{ID: id, Name: name, Description: description, SubmitTime: submitTime, Label: label})
}

return Tasks, nil
```

まず，Task の配列を宣言する．（ここは dbctl なので，先頭に dbctl.とつける必要はない．）

先ほどの db.Query では全タスクの一覧を取得した．なので，`for rows.Next() {`で取得したタスクの数だけ「取り出す」という処理を繰り返す．

Task には ID と名前と説明と登録時間とラベルが登録されているので，受け取るための変数を宣言する．

```
var id int
var name string
var description string
var submitTime string
var label int
```

そして，受け取る．

```
rows.Scan(&id, &name, &description, &submitTime, &label)
```

C 言語の scanf を思い出してくれた人もいるのではないだろうか．ここでも&が必要である．

そして，受け取ったタスクを Tasks 配列に追加する．

```
Tasks = append(Tasks, Task{ID: id, Name: name, Description: description, SubmitTime: submitTime, Label: label})
```

要するに，Tasks に新しい Task オブジェクトを加えた新しい Tasks を生成し，それを 古い Tasks に代入するということをしている．

動作原理は複雑だが，結局やっているのは「Tasks 配列に新しいタスクを追加する」という行為である．

あとは，返す．

```
return Tasks, nil
```

nil はエラーの nil である．（つまり，エラーがなかったということを返している．）

まとめると，

・データベースに問い合わせる

・その結果を変数に代入する

・結果を apifunc 側で受け取る

という流れで処理を行っている．

データベースの更新とかだと，エラー以外にデータベースから受け取るものがないので，2 番目が消えたり 3 番目がエラーだけを返したりといろいろ変わってくる（そこは　ケースバイケース）

### 処理結果を送信できるように変換する．

**あとちょっとで Go 編終わるぞ！** データベースを操作したのであとは結果をクライアントに返信するだけである．

/tasks の GET 処理の続きを見ていく．(apifuncs/tasks.go TasksFunc)

```
jsonBytes, err := json.Marshal(Tasks)
if err != nil {
	w.WriteHeader(http.StatusServiceUnavailable)
	fmt.Fprintln(w, `{"status":"Unavailable"}`)
	fmt.Println("JSON Marshal error(Tasks)", err)
	return
}

jsonString := string(jsonBytes)

w.WriteHeader(http.StatusOK)

if Tasks == nil {
	jsonString = "[]"
}

fmt.Fprintln(w, jsonString)
```

エラー処理は説明したので削る．

```
jsonBytes, err := json.Marshal(Tasks)

jsonString := string(jsonBytes)

w.WriteHeader(http.StatusOK)

if Tasks == nil {
	jsonString = "[]"
}

fmt.Fprintln(w, jsonString)
```

これだけ．

```
jsonBytes, err := json.Marshal(Tasks)
```

先ほどと真逆で，データを送信するために Byte 列に Tasks を変換する．

```
jsonString := string(jsonBytes)
```

それを文字列に変換する．（ちなみに，jsonString は普通に log.Print などを用いて読むことができる．なので，この一連の操作を文字列操作による実装で行うことも可能だが，大抵の場合は json.Marshal してから string に変換した方が断トツではやい．）

```
w.WriteHeader(http.StatusOK)
```

処理が成功したことをヘッダーに記述する．

```
if Tasks == nil {
	jsonString = "[]"
}
```

**これいらんかも（は？）**

```
fmt.Fprintln(w, jsonString)
```

最後の最後に処理の結果の集大成である jsonString を返信に記述して終わり．**ふう．**

まとめると，

・json.Marshal して Byte 列に変換してから， string(jsonBytes) で読める JSON に変換する．

・w.WriteHeader(http.StatusOK)しておく

・fmt.Fprintln(w, jsonString)でクライアントに送信するために結果を書き込む．

おしまい．

### Go の動作確認をする

テストする手法には様々な方法があるが，ここでは curl を使用したものを紹介する．

curl とは，さっき紹介した通り，サーバに http リクエストを送信することができるコマンドである．ターミナル上から入力する．

書き方は，こう．

```
curl -H "content-type: application/json" -X (POST,GET,PUT,DELETE) -d '(JSON)' (URL)
```

わかりにくいので，実際にこのサンプルで使用できるコマンドを紹介する．

```
curl -H "content-type: application/json" -X POST -d '[{"id":1, "name":"ダンス", "description":"さやかちゃんとダンスする", "submitTime":"2050-12-25 21:00:00", "label": 1},{"id":2, "name":"ダンス2", "description":"かけるくんとダンスする", "submitTime":"2050-12-25 21:00:03", "label": 1}]' http://localhost:8081/tasks
```

/tasks の POST はタスク一覧の配列を受け付けている．タスク一覧を送りつけると，データベースがそれに上書きされる．だから，タスクの配列をこちらから送信する．

読みにくいが，この JSON は要するに，ダンスとダンス 2 というタスクが入った配列を意味している．

docker container が立ち上がってる状態で上記のコマンドを実行してもらうと．確かに Tasks が上書きされて，ダンスとダンス 2 だけになることがわかる．（ただし，id=1,2 のタスクが存在している場合は上手くいかない．**いわゆる仕様．**)

```
curl -H "content-type: application/json" -X GET http://localhost:8081/tasks
```

こっちの方が簡単だ．これはデータベースからタスク一覧を取得してくる curl だ．

実行してみると，コンソール上にタスクの一覧が表示される．

### Go のおまけ

ところで，dbctl や apifuncs の関数は大体大文字から始まっていることに気がついただろうか．これは，Go の仕様で，他のパッケージから呼び出すには大文字で始まる名前にしなければならないというルールがあるからである．

例えば，apifuncs の TasksFunc から dbctl の GetTasks を呼び出す場面などがそれに該当する．

また，大文字で始まるものの前にはコメントで説明をつけなければならない．

```

// GetTasks タスク一覧を取得する関数
func GetTasks() ([]Task, error) {

```

`// 関数名 説明`という感じでコメントを書く．**書かないと lint がうるさい．**（うっとおしい）

---

話は変わる．

```

type Task struct {
    ID int `json:"id"`
    Name string `json:"name"`
    Description string `json:"description"`
    SubmitTime string `json:"submitTime"`
    Label int `json:"label"`
}

```

タスクの構造体である．これも ID だの Name だのの変数名は先頭を大文字にしないと怒られる．

ちなみに，フロントエンドで受信したデータは，`json:"hoge"`の hoge の部分がそのまま変数名になる．

```

console.log(res.data.id)
console.log(res.data.name)

```

こんな感じで，フロントエンドで使用可能．

## SQL の書き方

**正直，SQL というのは結構な時間をかけて勉強するものである．それを数分で詰めるというのはだいぶ無謀だと思うが，まあやるだけやってみよう．**

### データベースのかたち

本サンプルには，tasks と labels という二つのテーブルがある．それを実際に見てもらった方がわかりやすい．

tasks

![sql-1](images/sql-1.png)

labels

![sql-2](images/sql-2.png)

データベースとは言ってしまえば**ただの表**であることがお分かりいただけただろうか．

一つのタスクには id, name, description, submitTime, label の 5 つの属性がある．タスクが集まったものが tasks だ．

一つのラベルには id, label_text の 2 つの属性がある．ラベルが集まったものが labels だ．

データベースからデータを取得するというのは，要するに，これらのテーブルから全ての行を取り出したり，条件を指定して一部の行だけを取り出したりすることである．

データベースを更新するというのは，このテーブルに新しい行を挿入したり，行を削除したり，行の内容を変えたりすることである．

これからやるのはそういうことだ．

### データベースに失礼するゾ〜^

基本的に，Go からデータベースに SQL によるアクセスを行うことになる．が，いちいち Go に SQL を書いて docker-compose up してサーバーが立ち上がるのを待って SQL の結果を見て修正して docker-compose up して…と繰り返すのは**クッッッッッッッッッッッッッソ面倒である．やってられない．**

というわけで，データベースに直にアクセスして，バンバン SQL を試す方法を紹介する．

以下のコマンドを全て実行していく:

```

docker exec -it todo_sample_mysql bash
mysql -u gopher -p
(パスワードを入力．本サンプルのパスワードは'password'である．)
use prod_db

```

というコマンドを全て間違いなく入力してようやく SQL の試し撃ちができるようになる．

試しに`select * from tasks;`と入力すれば，現在のタスク一覧を確認できる．

しかし，docker-compose up で立ち上げたコンテナを終了させる度に強制ログアウトさせられるので，**毎回上 4 つを入力しないといけない．こっちはこっちでめんどいぞ！？！？？！？！**

慣れるとこれらのコマンドだけめちゃくちゃ早く入力できるようになる．暗記した頃には君も立派なバックエンドエンジニアになっているだろう（知らんけど．）

でも，コンテナが立ち上がっている間はサーバの中に入って SQL を好きなだけ撃てるので，Go に書いていちいちサーバを立ち上げる方法よりはずっとマシである．

### select

ここからは具体的に SQL について学んでいく．

`select id, name from tasks;`

こうすると

![sql-3](images/sql-3.png)

こうなる．

`select * from tasks;`

こうすれば

![sql-4](images/sql-4.png)

こうなる．

select．そのまんまである．

### update

`update tasks set description='UHOUHO';`

こうすると

![sql-5](images/sql-5.png)

task の全ての説明欄が**UHOUHO**になってしまった．

これは困るので(?)，id が 1 のタスクを「yokuneru」にしてみる．

`update tasks set description = 'yokuneru' where id = 1;`

![sql-6](images/sql-6.png)

このように，where をつけると更新する行を絞ることができる．
ちなみに，select でも where は使用できる．

`select * from tasks where id = 2;`

![sql-7](images/sql-7.png)

条件式には>=だの!=だのよく使う記号も使用できる．

### insert

`insert into tasks values(3, 'UHOUHO', 'uhouhouhouho.', '3000-01-01 00:00:00', 1)`

![sql-8](images/sql-8.png)

insert で行を増やすことができる．

### delete

`delete from tasks where id = 3;`

![sql-6](images/sql-6.png)

消せる．

### 結合

ところで，tasks の label は，ラベルの文字列そのものではなく labels テーブル の id を表している．

labels

![sql-2](images/sql-2.png)

というわけで，この二つの表は結合することが可能である．**やってみよう．**

`select * from tasks, labels where tasks.label = labels.id;`

![sql-9](images/sql-9.png)

なるほど．では，label と id が重複していて邪魔なので，select から除外してみる．

`select tasks.id, name, description, submitTime, label_text from tasks, labels where tasks.label = labels.id;`

![sql-10](images/sql-10.png)

このように，リレーショナルデータベースは結合演算を行うことが可能である．

なお，SQL を見てもらうと，`tasks.id`という様に，`テーブル.属性名`という記述の仕方がなされている．これは，tasks と labels の 2 つのテーブルを使用して演算しているため，tasks の id か labels の id かどちらの id かわからなくなるのを防ぐためにこういう記述をしている．

### group by

テーブルをグループ化させる．グループ化は SQL の重要な概念のひとつだが，理解するのが難しいので，順を追って説明していく．

とりあえず，上で散々遊び倒した tasks テーブルをリセットしつつ，適当なデータを増やす．

![sql-11](images/sql-11.png)

ちなみに，上のようにするには，以下を順に行う

```
・delete from tasks;をmysql上で実行
・ターミナル上で以下のコマンドを実行

curl -H "content-type: application/json" -X POST -d '[{"id":1, "name":"睡眠", "description":"よく寝る", "submitTime":"2020-08-25 21:00:00", "label": 0},{"id":2, "name":"ダンス", "description":"かけるくんとダンスする", "submitTime":"2050-12-25 21:00:03", "label": 1},{"id":3, "name":"電車ごっこ","description":"死ぬまでやめない", "submitTime":"2060-01-01 11:00:00", "label": 2},{"id":4, "name":"バーベキュー", "description":"焼くのはかけるくん", "submitTime":"2070-10-31 00:45:00", "label": 0},{"id":5, "name":"写生大会", "description":"余生を謳歌", "submitTime":"2080-04-01 00:00:14", "label": 2}]' http://localhost:8081/tasks
```

次に，ラベル で tasks を group 化する．こう書く:
`select label from tasks group by label;`

![sql-12](images/sql-12.png)

ラベル での tasks のグループ化に成功した．label の数字が並んでいるだけに見えるが，0 の行には id=1,4 の，1 の行には id=2 の，2 の行には id=3,5 のタスクの情報が含まれている．

なので，以下の SQL は失敗する．

`select label, id from tasks group by label;`

```

ERROR 1055 (42000): Expression #2 of SELECT list is not in GROUP BY clause and contains nonaggregated column 'prod_db.tasks.id' which is not functionally dependent on columns in GROUP BY clause; this is incompatible with sql_mode=only_full_group_by`

```

なぜなら，1 つの行に複数の行の情報が含まれているので，「どの id を表示すればいいんだ？」と MySQL が困惑してしまうからである．

つまり，グループ化すると，グループ化に使用した属性（今回の場合は label）しか表示できなくなる．**んな不便なことあるか？**

もちろん，グループ化すると嬉しいことがある．以下の SQL を実行してみて欲しい．

`select label, max(submitTime) id from tasks group by label;`

![sql-13](images/sql-13.png)

max は文字通り最大の値を求める関数である．だから，max(submitTime)で，グループの中で最大の submitTime（つまり，最も新しく登録されたタスクの登録日時）を取り出すことができる．

min もやってみよう

![sql-14](images/sql-14.png)

グループの中で最小の submitTime（つまり，最も昔に登録されたタスクの登録日時）を取り出すことができる．

関数は他にも sum だの ave だの count だの色々ある．試してみるといい．

### order by

並び替える．

`select * from tasks order by submitTime desc;`

![sql-15](images/sql-15.png)

desc は降順（大きいもの順）という意味である．だから，この SQL では，submitTime の大きいものから Task を取り出したということになる．

asc は小さいもの順だ．やってみよう．

`select * from tasks order by label asc;`

![sql-16](images/sql-16.png)

ラベルの数字が小さいものから取り出した．

ちなみに，asc は省略可能．

`select * from tasks order by label;`

あと，複数の条件を指定することもできる．

`select * from tasks order by label asc, submitTime desc`

### exists

存在するものを取り出す．

唐突だが，id=2 のダンスタスクには消えてもらう．

`delete from tasks where id = 2;`

![sql-17](images/sql-17.png)

label の列に注目してほしい．label が 1 のタスクは一つも存在しなくなった．

さて，labels から，tasks に存在しているラベルだけを取り出したいと思ったとする．（つまり，tasks には「なし」と「期限なし」のラベルがついたタスクだけが存在しているので，この 2 つのラベルを取り出す．）それを実現するには，こうする:

`select * from labels where exists (select * from tasks where tasks.label = labels.id);`

![sql-18](images/sql-18.png)

これが exists である．tasks.label の中に labels.id が存在した場合，そのラベルを取り出すという処理を行っている．

labels.id = 0 は tasks.label の中に存在している（睡眠タスクとバーベキュータスク）

labels.id = 1 は tasks.label の中に存在していない

labels.id = 2 は tasks.label の中に存在している（電車ごっこと写生大会）

よって，labels.id = 0,2 の 2 つのラベルだけが取り出される．

一方，not exists とすると，tasks.label に存在しない ラベル を取り出そうとする．

`select * from labels where not exists (select * from tasks where tasks.label = labels.id);`

![sql-19](images/sql-19.png)

### in

**書き方が違うだけで exists とほぼ一緒なので雑に紹介する．**

`select * from labels where labels.id in (select tasks.label from tasks);`

![sql-18](images/sql-18.png)

`select * from labels where labels.id not in (select tasks.label from tasks);`

![sql-19](images/sql-19.png)

こっちの方がわかりやすいという人もいるかもしれない．好きな方を選ぶといい．

### left outer join

突然だが，こんなタスクを挿入する．

![sql-20](images/sql-20.png)

**どういう Todo リストの使い方してるんだ**というツッコミはさておき，注目して欲しいのは label だ．

**id=3 のラベルは存在しない．**

![sql-2](images/sql-2.png)

では，このタスクが追加された tasks と labels を結合してみると，どうなる？

`select tasks.id, name, description, submitTime, label_text from tasks, labels where tasks.label = labels.id;`

![sql-21](images/sql-21.png)

新しく挿入した死のタスクがなかったことになってしまった．

ところが，なかったことにされては困る場合もある．

というわけで，この問題に対処するため，left outer join(左外部結合)というものを使ってみる．

`select tasks.id, name, description, submitTime, label_text from tasks left outer join labels on tasks.label = labels.id;`

![sql-22](images/sql-22.png)

存在しない label_text が NULL として現れた．

left outer join とは，左側（今回は tasks）を軸として，右側(今回は labels)を結合し，右側に存在しないものは NULL として表す演算である．

文法がえらくややこしく見えるが，覚えてしまえば書く難易度は低い．

`… from 左側テーブル名 left outer join 右側テーブル名 on 左側テーブル名.属性 = 右側テーブル名.属性`

または

`… from 左側テーブル名 left outer join 右側テーブル名 on 右側テーブル名.属性 = 左側テーブル名.属性`

である．

### between

**ぶっちゃけ便利そうな割にあんまり使わないのでざっくりと説明する．**

見た方がわかりやすい．

`select * from tasks where submitTime between '2060-01-01 00:00:00' and '2070-12-31 23:59:59';`

![sql-23](images/sql-23.png)

[属性]が〇〇から〇〇の間である行を取り出すという意味．今回は，[登録時間]が '2060-01-01 00:00:00' から '2070-12-31 23:59:59' の間である Task を取り出した．

### like

**これも使わん まれに使うことはあるかもしれない**

`select * from tasks where submitTime like '%2020%';`

![sql-24](images/sql-24.png)

submitTime に 2020 という文字列を含む行だけを取り出した．

ちなみに，Go から SQL を使用する際は日本語も使用できる．が，コンソール上では日本語を受け付けてくれないことが多い．

`select * from tasks where description like '%死%';`

この SQL を入力すれば，説明に**死**という文字列を含んだタスクだけを取り出すはずだ．**テストしてないけどどうせ動く．**

### case

先ほど，こんな SQL を紹介した．

`select tasks.id, name, description, submitTime, label_text from tasks left outer join labels on tasks.label = labels.id;`

![sql-22](images/sql-22.png)

この id 6 の label_text は NULL である．が，「LabelNashi」と表示したかったとする．

そのためには，次の SQL を入力する．（上の SQL からの変更点は，submitTime, ~ from tasks の~の部分である．）

```

select tasks.id, name, description, submitTime,
case when label_text is null then "LabelNashi" else label_text end
from tasks left outer join labels on tasks.label = labels.id;

```

![sql-25](images/sql-25.png)

このように，select 文や update 文の中で if 文のようなものを書くことができる．

`case when [条件] then [条件が真だった時] else [条件が偽だった時] end`

という文法である．

## as

表や属性に名前をつけることができる

`select tasks.id as 'UHOUHOID', uho.label_text as 'UHOUHOLABEL' from tasks, (select * from labels) as uho where tasks.label = uho.id;`

![sql-26](images/sql-26.png)

これは，tasks.id に `UHOUHOID` という名前を， label_text に`UHOUHOLABEL`という名前をつけ，

更に`select * from labels`の結果に`uho`という名前をつけた例である．

from に select を持ってくるのは，まあまあ使う．

ちなみに，as 自体は省略できる．

`select tasks.id 'UHOUHO', uho.label_text 'UHOUHO' from tasks, (select * from labels) uho where tasks.label = uho.id;`

### SQL のおまけ（と言いながら，割と重要な話．）

Go で SQL を使用する時の注意点を 2 つ説明する．

まず，Go で SQL を使用する際に，`db.Query`や`db.Exec`を使用することは既に紹介した．

SQL では複雑な条件を立てることができる．だから，Go の変数を SQL 内で使用したいという時もあるはずだ．

そういう時は，こう記述する．

```

hoge := 1
rows, err := db.Query("select * from tasks where id = ?", hoge)

```

?の中に変数の値が入る．この SQL では id が 1 のタスクを検索する．

たくさん変数を入れたい時は，どんどん?と変数を増やしていけばよい．

```

hoge := 6
fuga := "死"
piyo := "素晴らしい人生でした"
rows, err := db.Query("select * from tasks where id = ? and name = ? and description = ?", hoge, fuga, piyo)

```

---

話は変わる．

```
select tasks.id, label_text from tasks left outer join labels on tasks.label = labels.id;
```

![sql-27](images/sql-27.png)

これを rows.Scan で受け取るとする．

先に申し上げると，この rows.Scan は**失敗する．**原因は，id が 6 の label_text が**NULL**であることである．

SQL の結果に NULL が含まれる可能性がある場合，受け取る側で工夫をしなければならない．このようにする:

```

var id int
var labelText sql.NullString = sql.NullString{}
rows.Scan(&id, &label)
HogeArray = append(HogeArray, Hoge{ID: id, LabelText: labelText.String})

```

sql.NullString という型で変数を宣言し，それを rows.Scan で使用する．変数の値を使用する時は`変数名.String`と記述する．

名前の通り，String 専用である．数字を扱いたい時は`sql.NullInt64 = sql.NullInt64{}`を使用すればよい．

---

話は変わる

init.sql というファイルが mysql/db の中にあるのを見たことがあるかもしれない．これは，データベースの初期化に使用される一連の SQL だ．

```

create table labels(
id int primary key not null,
label_text varchar(64) not null
);

create table tasks(
id int primary key not null,
name varchar(128) not null,
description varchar(256) not null,
submitTime datetime not null,
label int not null
-- foreign key (label) references labels (id) on delete cascade
);

```

文字通りテーブルを作っている．primary key は主キーという意味で，そのテーブルのなかで唯一で無ければならない．

つまり，既に tasks に存在している ID のタスクを insert しようとすると，エラーとして処理される．

not null は　まあ　どんなに英語ができない人でもなんとなくわかるだろう．

最後のコメントアウトされている foreign key は，tasks.label は labels の id ですよ〜という意味だと思っておけばいい．

これによって，本当であれば，labels に存在しないラベルのタスクは tasks に挿入することができなくなる．

本当は書いておかなければならないが，学習の都合上（主に left outer join のせい）でコメントアウトしている．

```

insert into
labels(id,label_text)
values
(0,"なし");

...

```

あとは初期値を insert しているだけ．

## おめでとう！これで君も SQL マスターだ！

**当然嘘である**．だが座学で教えられるのはこれが限界だ．

実際に書いて，慣れていって欲しい．残念ながら SQL だけは慣れてもらうしかない．

## さいごに

これを読み終わったら app ディレクトリの中の main.go, apifuncs/tasks.go, apifuncs/labels.go, dbctl/tasks.go, dbctl/labels.go を読んでみて欲しい．多分大体読めると思う．

あと，backend-homework.md という形で宿題も用意しておいた．やると多分力がつく．多分．
