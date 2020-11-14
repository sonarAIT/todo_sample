# 30 分でわかれ Go と SQL

**30 分でやるの普通に無理だと思う**

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

http のリクエストにはデータを添付することができる．それが body．**GET だと送れないので注意**（厳密には curl を使えば送れる．が，axios などでは基本的には送れないので，多分非推奨なんだろう．多分．）

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

80 番ポートに関してはどうでも良くて，ここで覚えて欲しいのは，**受け取る URL の数を増やしたかったら http.HandleFunc と対応する関数を増やせばいい**ということである．

例えば，このアプリケーションにログイン機能を実装して，ログインの http リクエストを受け取ることになった場合，

```
http.HandleFunc("/login", apifuncs.LoginFunc)
```

のように書けばよい．

次項からは，その実行される関数の書き方について記述していく．

### リクエストの種別(GET, POST, PUT, DELETE)を識別する．

apiFuncs のファイル（どっちでもいい）を見てみよう．

```
w.Header().Set("Access-Control-Allow-Origin", "*")                       // Allow any access.
w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE") // Allowed methods.
w.Header().Set("Access-Control-Allow-Headers", "*")

r.Header.Set("Content-Type", "application/json")
```

**これは　もう　書け．覚えろ．**（ちょっとだけ説明すると，上 3 つは受け取るリクエストのセキュリティ的な設定である．これを見ると**全種類のリクエストを受け取る**というめちゃくちゃガバガバなセキュリティになっていることがわかる．しかし，実際にデプロイしないハッカソンとかでは**どうでもいい**）

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

/tasks の POST では「クライアントからタスク一覧を受け取り，受け取ったタスク一覧に合わせてデータベースを更新する」という処理を行っている．

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

次に，リクエストを受け取るための変数を宣言する．/tasks の POST には，タスクの配列（タスク一覧）が送られてくる．だから，

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

・`var recTasks []dbctl.Task`で body を受け取るための変数を宣言．

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

さて，上記によれば dbctl の GetTasks()の結果を Tasks に代入しているらしい．ここから飛んで，dbctl/tasks.go の GetTasks を参照する．

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

真ん中の`if err != nil`の中身がいつもとちょっと違うが，結局はエラーをログに書いてエラーを呼び出し元（つまり apifuncs.TaskFunc）に返すという処理である．（クライアントに通達するのは呼び出し元の仕事なので，ここでは行わない．）

つまり，本質的には，

```
rows, err := db.Query("select id, name, description, submitTime, label from tasks")
defer rows.Close()
```

こう．

db.Query はデータベースからデータを取得する際に書く処理である．引数にはみんな大好き SQL が記述されている．（後述．**こんな隙間で済むほど単純ではない．**）

rows に問い合わせの結果を，err にエラーを代入する．

`defer rows.Close()`は，後で(dbctl.GetTasks が終わる際に)rows.Close()を行うという意味である．**defer rows.Close は忘れてはいけない**．忘れるとメモリがリークする．

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

先ほどの db.Query では全タスクの一覧を取得した．なので，`for rows.Next() {`で取得したタスクの数だけ処理を繰り返す．

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

要するに，Tasks に新しい Task オブジェクトを加えた新しい Tasks を生成し，それを 古い Tasks に代入するということをしている．（動作原理は複雑だが，結局やっているのは「Tasks に新しいタスクを追加する」という行為である．）

あとは，返す．

```
return Tasks, nil
```

nil はエラーの nil である．（つまり，エラーがなかったということを返している．）

まとめると，

・データベースに問い合わせる

・その結果を変数に代入する

・結果を apifunc 側で受け取る

という流れで処理を行っている．

データベースの更新とかだと，エラー以外に受け取るものがないので，2 番目が消えたり 3 番目がエラーだけを返したりといろいろ変わってくる（そこは　ケースバイケース）

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

`jsonBytes, err := json.Marshal(Tasks)`は先ほどと真逆で，データを送信するために Byte 列に Tasks を変換する．

`jsonString := string(jsonBytes)`でそれを文字列に変換する．（ちなみに，jsonString は普通に読むことができる．なので，この一連の操作を文字列操作による実装で行うことも可能だが，大抵の場合は json.Marshal してから string に変換した方が断トツではやい．）

`w.WriteHeader(http.StatusOK)`で処理が成功したことをヘッダーに記述する．

```
if Tasks == nil {
	jsonString = "[]"
}
```

**これいらんかも（は？）**

`fmt.Fprintln(w, jsonString)`

最後に返信に処理の結果の集大成である jsonString を記述して終わり．**ふう．**

まとめると，

・json.Marshal して Byte 列に変換してから string で 読める JSON に変換する．

・w.WriteHeader(http.StatusOK)しておく

・fmt.Fprintln(w, jsonString)でクライアントに送信するために結果を書き込む．

おしまい．

### Go のテストをする

テストする手法には様々な方法があるが，ここでは curl を使用したものを紹介する．

curl とは，さっき紹介した通り，サーバに http リクエストを送信することができるコマンドである．ターミナル上から入力する．

書き方は，こう．

```
curl -H "content-type: application/json" -X (POST,GET,PUT,DELETE) -d '(JSON)' (URL)
```

わかりにくいので，実際にこのサンプルで使用できるコマンドを紹介する．

```
curl -H "content-type: application/json" -X POST -d '[{"id":1, "name":"ダンス", "description":"さやかちゃんとダンスする", "submitTime":"2050/12/25 21:00:00", "label": 1},{"id":2, "name":"ダンス2", "description":"かけるくんとダンスする", "submitTime":"2050/12/25 21:00:03", "label": 1}]' http://localhost:8081/tasks
```

/tasks の POST はタスク一覧の配列を受け付けている．タスク一覧を送りつけると，データベースがそれに上書きされる．だから，タスクの配列をこちらから送信する．

読みにくいが，この JSON は要するに，ダンスとダンス 2 というタスクが入った配列を意味している．

docker container が立ち上がってる状態で上記のコマンドを実行してもらうと．確かに Tasks が上書きされて，ダンスとダンス 2 だけになることがわかる．

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

```

type Task struct {
    ID int `json:"id"`
    Name string `json:"name"`
    Description string `json:"description"`
    SubmitTime string `json:"submitTime"`
    Label int `json:"label"`
}

```

タスクの構造体である．これも変数名は先頭を大文字にしないと怒られる．

ちなみに，フロントエンドで受信したデータは，`json:"hoge"`の hoge の部分がそのまま変数名になる．

```

console.log(res.data.id)

```

こんな感じ．

## SQL の書き方

あとで書く．

## さいごに

これを読み終わったら app ディレクトリの中の main.go, apifuncs/tasks.go, apifuncs/labels.go, dbctl/tasks.go, dbctl/labels.go を読んでみて欲しい．多分大体読めると思う．

あと，backend-homework.md という形で宿題も用意しておいた．やると多分力がつく．多分．
