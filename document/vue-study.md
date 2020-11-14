# 30 分でわかれ Vue.js

## はじめに

このドキュメントはなるべく早く Vue.js を書けるようになってもらうことを目標としています．頑張ってください．

なお，厳密な理解を求めている人やより深く学びたい人は，Vue.js のエロ本を買うなりしてください．

## JavaScript とは？

Web サイトで使われがちなスクリプト言語．C が書ければボチボチ書ける．
C と違って困惑しそうな点は，配列が完全に別物な点とオブジェクトという新概念がある点である．その辺は各自で調べておいて欲しい．（これは，Vue.js のドキュメントなので…）

大抵の場合，JavaScript の部分と html の部分にわけて記述する．

## Vue.js とは？

すげえ JavaScript．

## Vue.js JavaScript の知っておいて欲しい概念

### console.log

JavaScript の printf．

### デベロッパー・ツール

大体のブラウザに備わっている機能．Google Chrome では F12 を押すことで見えるようになる．（Macbook では Fn キーと F12 の同時押し）
これの Console というタブで上記の Console.log の中身を見ることができる．

### App.vue

ページの親玉．ここにパーツを置いていってページを構成する感じになる．

### components

component とは，ページに置くパーツである．これを App.vue に置いたり component の中に component を置いたりしながらページを作りあげていく．

もちろん，component を一切使わずに，App.vue に全てを記述することも可能なのだが，component を使うメリットもある．ページの要素をパーツに分けることによって，作業を分担しやすくなったり，読みやすくなったりする．

本サンプルには，TaskForm.vue，Filter.vue，Tasks.vue，Task.vue の 4 つの component がある．

### Vuex Store

C 言語で言うところのグローバル変数．全てのページや component から参照できる．欠点として，どこから書き換えてるか store 側からはわからないので，使うとバグの詮索がしにくくなっていく．しかし，適切に使うと，コードが短くなる．

store フォルダの中の index.js がこれにあたる．

### axios

Vue.js がサーバやその他諸々と通信するために使用する．URL を指定し，データを添付することでサーバにデータを送ったりデータを受け取ったりできる．

ちなみに，http の知識になるが，送信する http のリクエストには 4 つの種類がある．

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

という 4 種類の処理を割り当てることができる．受け取った4 種類のリクエストを割り振るのは，バックエンドの仕事になる．（仕事と言っても，そう難しくは無いが…）

## JavaScript 部分の書き方

App.vue，TaskForm.vue，Filter.vue，Tasks.vue，Task.vue の下半分部分について書いていく．

### ざっくりと見てみる

上記 5 つのどれでもいいから見てもらうと，

```
hoge: {
}
```

というものがずらっと並んでいる．data だの methods だの色々種類があるので，それぞれなんなのか大事な物から説明していく．

### data

変数置き場．

```
data() {
    return {
        hoge: 0,
        fuga: "Hello"
    }
}
```

という感じで書いていくと，hoge や fuga がそのファイル内で使用できる．

```
console.log(this.hoge)
```

**気をつけて欲しいのだが，Script 部分では `this.hoge`という感じで this を付けないといけない．**（付けないと変数が定義されていないと言われて怒られる．）

### props

親から変数を受け取る機能．本サンプルでは TodoFilter.vue などで使用されている．
先ほど，component はパーツと述べた．なので，そのパーツを置いている場所（親．つまり，TodoFilter.vue から見ると App.vue）から変数が受け取れたら便利な場合がある．そこで，

```
props: {
    propHoge: Number,
    propFuga: String,
},
```

と書くと， prop の名前と型を書けば親から変数を受け取ることができる．
親では，

```
<ComponentHoge v-bind:propHoge="hoge" v-bind:propFuga="fuga">
```

と記述する．（hoge と fuga は変数．）
なお，v-bind の特性によって，親の変数の値が変更されると子の値も変更される．親の hoge の値が変わったら，子の propHoge の値も変わるということだ．（v-bind に関する詳しい話は後ほど行う．）

### methods

関数置き場．method ではなく method**s** である．（これを間違えるとエラーが全く出ないのにうまく動作しない地獄に陥るので注意して欲しい．）

```
methods: {
    hogeFunc: function() {
        console.log("ほげほげ〜")
    },
    fugaFunc: function(fuga) {
        console.log(fuga);
    }
}
```

という感じで記述する．この関数は同じファイル内のあらゆる所から呼び出せる．例えば

```
<button @click="hogeFunc()">ほげる</button>
```

と書けば，押せば hogeFunc が実行されるボタンが出来上がる．(他にも色々呼び出し方がある．後述する．)
例として，TaskForm.vue では，タスク登録ボタンを押した時に実行される関数が記述されている．

### computed

computed とは，つまる所，データを加工したものに名前をつけた物である．Tasks.vue に，tasks をフィルタリングした結果を返す computed の例がある．

```
computed: {
    numInc: function() {
        return this.num+1
    },
},
```

このように書く．numInc は，data（もしくは props）の num に 1 を足した結果を返す．即ち，num というデータを加工した結果を返していると考えられる．
computed に登録したものは，変数のようにして扱うことが可能である．例えば，

```
method: {
    catNumInc: function() {
        console.log(this.numInc)
    }
}
```

という様に書く．

computed を使用するメリットは，結果がキャッシュされることである．即ち，ファイル内でどれだけ numInc を呼び出しても，実際に関数が実行された回数（num に 1 を足すという演算を行った回数）は 1 回だけである．num 自体が変更されるまでは二度と実行されることはない．これは，即ち，処理の軽量化に貢献する．

### watch

変数の値が変更された時に呼び出される関数である．TodoFilter.vue に例がある．

```
watch: {
    hoge: function() {
        console.log("hogeの値が変更されました．")
    }
}
```

という様に記述する．
便利なのだが，正直，挙動不審である．（変数の値が変更されたのに実行されないということがたまによくある．）
実際，仕様通りに動いているだけなのだが，配列の値が変わった時に反応しなかったり，オブジェクトの中の値が変更された時に反応しなかったり，他にも色んな時に反応しない．もちろん対処法もあるのだが，その対処法を調べるのに結構な時間を取られることがあるので，**ぶっちゃけ使わない実装をするのが一番早い．**

### name

コンポーネントの名前を定義する．これを書かないと認識されない．忘れて時間を取られがちなので注意．

### components

使用する component を定義する．App.vue が一番わかりやすい．
component を使用する親では，import して，component に定義するという一連の流れを行わないといけない．

```
import HogeComponent from "./components/HogeComponent.vue";

export default {
  components: {
    HogeComponent,
  },
};
```

と書いて，ようやく html で HogeComponent を使用することができる．

### created

ページ及び component が表示された際（起動時）に実行される関数．（厳密に言えば表示よりもっと早いタイミングだったはずだが，どうでもいい．）App.vue に例がある．
これだけ定義の仕方が他と明らかに違う．

```
created() {
    console.log("表示しました")
},
```

という感じで，関数を直に書く．

似たようなものに，`mounted` があるが，大体同じようなものである．（違いが気になった人は調べて欲しい．）
また，逆に，ページを閉じる際に実行される `destoryed` がある．

## HTML 部分の書き方

下半分を説明したので，今度は上半分を説明する．

### {{hoge}}

このように書くことで，data や computed 内の変数の値をページ上に表示することができる．
TaskForm.vue にて~~無理やりねじ込んだ~~使用例を見ることができる．

### タグに付与できるおもしれー属性

タグとは，h1 だの p だの a だの input だの img だの label だのとにかく html の授業で死ぬほどやったであろう<>で包んでるアレである．これらには属性を付与することができる．（例えば，a の href 属性はリンクで飛ぶ先を指定する属性である．）
ところで，Vue.js を説明する文句の一つに「html の中に JavaScript を書く感じ」という説明文句がある．（ちなみに，React は真逆で，JavaScript の中に html を書く感じ．）

その html の中に JavaScript を書いている感じを**ゴリゴリ**に体感させてくれるのが，Vue.js 特有のタグに付与できる属性である．html を*ぐにゃぐにゃ*に変形させる様々なタグを見ていこう．

### v-bind

**一番大事**．これを使えなかったら Vue.js をやったことがあると名乗ってはいけない．

v-bind とは，つまるところ，タグの属性の設定に変数を使用することができる機能である．

例えば，data に`URLString: "http://google.com"`と定義されていたとする．これを a で使うには，

```
<a v-bind:href="URLSring">googleへ</a>
```

と書けば良い．
また，JavaScript 部分によって URLString が変更されると，a によって飛ぶリンク先も変化する．html を動的に変化させることができるのである．**ビバ，Web プログラミング！**

先ほどの例で述べたように一番重要なので，もう一つ例を紹介する．<input>タグにおいて

```
<input type="text" value="HOGE">
```

と書けば，最初から HOGE と入力されている入力欄が出来上がる．これを v-bind を使ってバインディングさせると

```
<input type="text" v-bind:value="hoge">
```

data 内の変数 hoge の値が入力欄に表示されることになる．
もちろん，hoge の値が変われば input 内の表示も変わる．

ただし，欠点がある．input は文字通り入力欄で，type="text"ということは 文字を入力するための入力欄ということになる．つまり，入力欄を書き換えたら当然 data の hoge も書き換えて欲しいところなのだが，それは叶わない願いである．なぜなら，v-bind はあくまで data などの変数を html 上のタグの属性の値として使用できるだけであって，html 上で値が変わろうと JavaScript の方を書き換えることはできないからである．JavaScript から html に対して一方通行なのだ．

次項でこれに対する対処法を説明する．

### v-model

上記の問題に対処する物だと思えばいい．v-bind では JavaScript 側から HTML 側に一方通行であるが，v-model ではそうでない．

```
<input type="text" v-model="hoge">
```

と書けば，input の入力欄は変数 hoge の値になるのはもちろんのこと，入力欄をユーザ側で書き換えると，なんと，JavaScript の hoge の値も書き変わる．

例を見た方が早い．TaskForm.vue にて，名前入力欄と説明入力欄とラベル選択メニューで使用されている．

### v-if

条件式によってそのタグが表示されるかどうかを決める．

data 内に`hogeBool: false`のように記述されていたとする．

```
<p v-if="hogeBool">ウオオオオ</p>
```

このタグは表示されない．なぜなら，hogeBool は false だからである．
JavaScript 側で hogeBool が true になった瞬間に，初めてこの**ウオオオオ**は表示される．

ちなみに，v-if の中はあくまで条件式なので，

```
<p v-if="hogeNum >= 0">ウオオオオ</p>
```

のような書き方も可能である．
一番わかりやすい例は，TodoFilter.vue のフィルタ選択部分である．フィルタを使用しない場合，フィルタ選択のセレクタは表示すらされない．

### v-for

配列の数だけ要素を繰り返す．

data 内に`hogeArray: ['aaa', 'iii', 'uuu']`のように記述されていたとする．

```
<p v-for="hogeString in hogeArray" v-bind:key="hogeString">{{hogeString}}</p>
```

と書けば，きっちり 3 つの p が表示される．

```
aaa

iii

uuu
```

ここでの hogeString は，配列のうちの 1 つの要素を意味している．**(data 内の変数では無いことに注意)**

ちなみに気をつけなければいけないのが，v-bind:key である．keyを設定することにより，vueが配列の要素を識別することが可能になり，動作が安定する．今回のkeyは文字列そのものである．v-bind:key は設定しなくとも動くが，基本的に設定しなければならない．

```
<p v-for="(hogeString, index)  in hogeArray" v-bind:key="index">{{hogeString}}</p>
```

という書き方をすると，index（つまり，配列で何番目か）を key にするが，**この書き方をしてはいけない**．配列の要素を削除すると，削除した要素の後の要素の index が全て変わるため，即ち key が変わることになる．これによって，様々なバグが発生する．

Tasks.vue ではタスクの数だけタスクを表示する v-for の記述がある．ここでは，Task オブジェクトの id をキーとしている．

### @click

クリックされた際に実行される関数を指定する．

```
<button @click="hogeFunc()"></butotn>
```

とすれば，button を押した際に hogeFunc が実行される．

### @change

値が変更された際に実行される関数を指定する．

```
<input type="text" @change="hogeFunc()" />
```

とすれば，入力欄を書き換えた際に hogeFunc が実行される．

## store の使い方

先に述べたように store とはグローバル変数置き場のような物である．store ディレクトリ内の index.js に設定を記述していく．その使い方を説明する．

### store から値を呼び出す方法

`this.$store.state.(変数名)`と記述する．長い．（長いので，サンプルの一部の component の computed では， this.\$store.state.labels が labels だけで呼べるように定義されている．）

### state

変数置き場．

```
state: {
    hoge: "hoge"
},
```

data と似たようなものだ．見ればわかる．

### mutations

今まで一度も述べたことが無いが，

```
this.$store.state.hoge = hoge
```

のように，直接 state の変数に代入するのは**違法行為**である．
なので，mutations という関数を介して代入する必要がある．

```
mutations: {
    setHoge(state, hoge) {
        state.hoge = hoge;
    }
}
```

こう定義する．代入したいときは

```
var hoge = 1
this.$store.commit('setHoge', hoge)
```

こう．

## acitons

state 内のデータをサーバに送信したり，とにかく state の変数でややこしいことをする時に使用する．

定義は

```
action: {
    hogeAction(context, hoge){

    },
}
```

こうで，呼び出すときは

```

this.$store.dispatch('hogeAction', hoge)

```

こう．

見ての通り，問答無用で第一引数に context をとる．（取らないこともある）context からは，

```

context.state.hoge              //state 内の hoge を呼び出す
context.commit('setHoge', hoge) //mutations 内の関数を呼び出す
context.dispatch('fugaAction')  //他の action を呼び出す

```

こういったことができる．

ちなみに，action は context ともう一つしか引数を受け取ることができない．なので，2 つ以上 action に引数を送りたければ

```

this.$store.dispatch('hogeAction', {hoge1: 0, hoge2: "hogehoge"})

```

という様に，オブジェクトで送る必要がある．

なお，通信を行う関数は，先頭に async を付けなければならない．**そういうもんだ．**（理由はググれ．）

### axios

厳密には store の機能では無いが，このサンプルでは store の中でしか使用されていないので，ここに記述する．

かなり前の方で述べたとおり，サーバなどとの通信に使用する．GET POST PUT DELETE の 4 種類のリクエストから選んで通信を行う．（どんなリクエストをどういう形式で送るか，というのは，普通設計段階で決める．）

記述の方法は，

```
await axios.get(URL)
    .then((サーバから帰ってきたデータ) => {
        //サーバから帰ってきた後の処理
    })
    .catch((サーバから帰ってきたデータ) => {
        //エラーが発生した時の処理
    })

await axios.post(URL, 送信するデータ)
    .then((サーバから帰ってきたデータ) => {
        //サーバから帰ってきた後の処理
    })
    .catch((サーバから帰ってきたデータ) => {
        //エラーが発生した時の処理
    })


```

という感じ．delete も put も似た様な物である．
ちなみに，**get ではデータを送ることができないので注意．**（クエリパラメータといって URL 上にパラメータを書く必要がある．）

## おまけ

### .env.developmentとは

開発用 URL 置き場．本当は開発用の URL だの本番用の URL だの様々なケースに応じた URL を定義することができる．ここに URL を書いておけば，本番に切り替える際に store 等のソースコードを書き換える必要がなくなる．

でもサンプルでは開発用の URL しか定義していない．どうせデプロイしないし．

## さいごに

これを読み終わったら App.vue と components と store の中身を読んでみて欲しい．多分大体読めると思う．

あと，vue-homework.md という形で宿題も用意しておいた．やると多分力がつく．多分．
