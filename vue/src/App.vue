<template>
  <div id="app">
    <!-- component TaskFormを設置 -->
    <TaskForm />
    <!-- brは空白の意 -->
    <br />
    <!-- component TodoFilterを設置
    v-bindでpropFilterに変数nowFitlerを代入．
    changeイベントを取得したら，nowFilterChangedを実行する．
    -->
    <TodoFilter v-bind:propFilter="nowFilter" @change="nowFilterChanged" />
    <br />
    <!-- 上二つ言ったらこれ説明しなくてもいいでしょ -->
    <Tasks v-bind:nowFilter="nowFilter" />
  </div>
</template>

<script>
// 3つのcomponentを使用するためにimportする
import TaskForm from "./components/TaskForm.vue";
import TodoFilter from "./components/TodoFilter.vue";
import Tasks from "./components/Tasks.vue";

export default {
  name: "App",

  // component3つを定義
  components: {
    TaskForm,
    TodoFilter,
    Tasks,
  },

  // 変数置き場．
  data() {
    return {
      nowFilter: -1,
      // 現在使用しているフィルタを指定するnowFilter
      // -1はフィルタを使用しないという意味．0以降は，数字に対応したlabelでフィルタリングする．
    };
  },

  // 関数置き場．
  methods: {
    // TodoFilterのnowFilterが変更されたら，こっちのnowFilterも変更する．
    nowFilterChanged: function(nowFilter) {
      this.nowFilter = nowFilter;
    },
  },

  // ページを開いた時に実行される．
  created() {
    this.$store.dispatch("getTasks"); // サーバからタスク一覧を取得
    this.$store.dispatch("getLabels"); // サーバからラベル一覧を取得
    // 詳しい処理はstore/index.jsのactionにて．
  },
};
</script>
