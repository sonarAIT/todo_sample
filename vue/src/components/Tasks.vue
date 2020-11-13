<template>
  <!-- タスク一覧テーブル -->
  <!-- html tableとかでググれば全部わかるので，説明は割愛．（というか，webプログラミングで習うでしょ．) -->
  <table>
    <thead>
      <tr>
        <th>ID</th>
        <th>タイトル</th>
        <th>コメント</th>
        <th>登録時間</th>
        <th>ラベル</th>
        <th>完了ボタン</th>
      </tr>
    </thead>
    <tbody>
      <!-- v-forを使って，タスクの数だけTaskを表示している． -->
      <!-- v-bind:taskでTaskのpropsのためにtaskitemを代入している． -->
      <!-- 詳しくはTask.vueを参照 -->
      <Task
        v-for="taskitem in tasks"
        v-bind:key="taskitem.id"
        v-bind:task="taskitem"
      />
    </tbody>
  </table>
</template>

<script>
// component Taskをimport
import Task from "./Task.vue";

export default {
  // このcomponentの名前はTasks
  name: "Tasks",

  // component Taskを定義．
  components: {
    Task,
  },

  // 親であるApp.vueからnowFilterを受け取る．
  // 即ち，現在のフィルターの設定．
  props: {
    nowFilter: Number,
  },

  // 加工したデータをお届けすることでお馴染みのcomputedゾーン
  computed: {
    // もしフィルタを使うのであれば，nowFilterの値に応じたフィルタリングを行う
    // もしフィルタを使わないのであれば（nowFilterが-1ならば）フィルタしていないそのままの結果をお届けする
    // filterメソッドについてはArray.prototype.filter()でググれ．
    tasks: function() {
      if (this.nowFilter != -1) {
        return this.$store.state.tasks.filter(
          (task) => task.label === this.nowFilter
        );
      }
      return this.$store.state.tasks;
    },
  },
};

/*
フィルタリングの処理の部分がわかりにくいかもしれないので，スペシャルな解説を行う．
まず，タスクにはそれぞれlabelという物が設定されている．labelは数字である．
これは，labelのテキスト自体ではなくidを登録しているためである．

次に，storeの中にはlabelsというラベル一覧の配列がある．
これは，id(数字)とlabelText(文字列)がセットになったラベルのオブジェクトの配列である．
idは，配列上で何番目になるかを表せる様に設定してある．

上記の事を踏まえると，labels[label]とすれば，タスクが設定しているラベルに応じたラベルのオブジェクトが取り出せることがわかる．
このことは，Task.vueで利用している様子を見ることができる．

さて，これをフィルタリングするには，絞りたいラベルのidを指定して，指定したidとtaskのラベルのidが等しい場合にtaskを取り出せばよい．
ラベル「緊急」のidは1である．「緊急」のタスクだけを取り出すには，
nowFilterを1にして，nowFilterとtaskのidが等しい場合だけタスクを取り出していく．
そうして，フィルタリングされた結果のtaskの配列が生まれる．

上のフィルタリングでやっていることは，そういうことである．
*/
</script>
