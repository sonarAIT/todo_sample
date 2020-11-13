<template>
  <div>
    <!-- data内の変数pageTitleの値を表示 -->
    <h3>{{ pageTitle }}</h3>

    <!-- タスクの名前設定 -->
    <div>
      <!-- id nameIDのためにラベルを設定 -->
      <label for="nameID">名前</label>
      <!-- v-modelによってdata内の変数nameと同期 -->
      <input type="text" id="nameID" v-model="name" />
    </div>

    <!-- タスクの説明設定．本質的には名前設定と完全に一緒 -->
    <label for="descriptionID">説明</label>
    <div>
      <textarea
        rows="4"
        cols="40"
        id="descriptionID"
        v-model="description"
      ></textarea>
    </div>

    <!-- ラベルの設定．本質的には名前設定と一緒 -->
    <div>
      <label for="labelID">ラベル </label>
      <!-- 選択した値によってdata内のlabelを変更する． -->
      <!-- v-bind:valueの値がdata内のlabelに代入される． -->
      <!-- this.$store.state.labelsはstore内にあるラベル一覧である -->
      <!-- つまり，選んだラベルのIDがdata内のlabelに代入される -->
      <select id="labelID" v-model="label">
        <option
          v-for="labelitems in this.$store.state.labels"
          v-bind:key="labelitems.id"
          v-bind:value="labelitems.id"
          >{{ labelitems.labelText }}</option
        >
      </select>
    </div>

    <!-- 登録ボタンを押すと，methodsのsubmitが実行される． -->
    <button @click="submit">登録</button>
  </div>
</template>

<script>
export default {
  // コンポーネントの名前を設定．
  name: "TaskForm",

  // 変数置き場．
  // pageTitle以外は，v-modelによってユーザが入力した値を受け取る変数である．
  data() {
    return {
      pageTitle: "タスク登録フォーム",
      name: "",
      description: "",
      label: 0,
    };
  },

  // 関数置き場
  // タスク登録のための関数がひとつだけある．
  methods: {
    submit: function() {
      // 現在日時を取得
      let date = new Date();
      // 登録するタスクのためのIDを決める
      // もしタスクがひとつもなかったら，idは1．
      // そうでなければ，タスク一覧の最も後ろのタスクのIDに1を足した物を次のIDとする．
      let nextID;
      if (this.$store.state.tasks.length == 0) {
        nextID = 1;
      } else {
        nextID = this.$store.state.tasks.slice(-1)[0].id + 1;
      }
      // 新しいタスクのタスクのオブジェクトを作る．
      // submitTimeは登録した日付と時間．
      // console.log(date.toISOString().substr(0, 10) + " " + date.getHours() + ":" + date.getMinutes() + ":" + date.getSeconds())
      // を実行してもらえれば，何をしたいのかわかるはず．
      let task = {
        id: nextID,
        name: this.name,
        description: this.description,
        label: this.label,
        submitTime:
          date.toISOString().substr(0, 10) +
          " " +
          date.getHours() +
          ":" +
          date.getMinutes() +
          ":" +
          date.getSeconds(),
      };
      // storeにて，タスクを追加する際の処理を実行．
      // index.jsのactions内のaddTaskを参照．
      this.$store.dispatch("addTask", task);
    },
  },
};
</script>
