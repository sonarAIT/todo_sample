import Vue from "vue";
import Vuex from "vuex";
import axios from "axios";

Vue.use(Vuex);

// たのしいVuex Storeへようこそ
export default new Vuex.Store({
  // 変数置き場．
  // この変数へは，this.$store.state.変数名でどこからでもアクセスできる．
  state: {
    tasks: [
      {
        id: 1,
        name: "タスク読み込みエラー",
        description: "タスクを読み込めなかった時に出るタスクです",
        submitTime: "2020/01/01 00:00:00",
        label: 2,
      },
    ],
    labels: [
      {
        id: 0,
        labelText: "none",
      },
      {
        id: 1,
        labelText: "none",
      },
      {
        id: 2,
        labelText: "none",
      },
    ],
  },

  // mutations
  // stateの変数に代入するときは，このメソッドを使って代入しないといけない．
  mutations: {
    // stateのtasksに引数のtasksを代入
    setTasks(state, tasks) {
      state.tasks = tasks;
    },
    setLabels(state, labelTexts) {
      state.labels = labelTexts;
    },
    // タスクを加える．配列に加えるときはpushを使う．
    addTask(state, task) {
      state.tasks.push(task);
    },
    // タスクを取り除く．index番目から1つ取り除くという意味．
    removeTask(state, index) {
      state.tasks.splice(index, 1);
    },
  },

  // actions
  // stateの変数が主役の複雑な処理を実行する．
  actions: {
    // Tasksをサーバから取得する．
    // process.env.VUE_APP_API_TASKSはhttp://localhost:8081/task である．(.env.developmentにそれが書いてある．)
    async getTasks(context) {
      await axios
        .get(process.env.VUE_APP_API_TASKS)
        .then((res) => {
          // 通信が成功したら，受け取ったデータでsetTask(mutationsにある)を実行．
          // つまり，サーバから受け取った物をstateのTasksに代入する．
          context.commit("setTasks", res.data);
        })
        .catch(() => {
          // 失敗したらこれを表示．
          console.log("Tasksのgetに失敗しました.");
        });
    },

    async postTasks(context) {
      // Tasksをサーバに送信する．
      // URLが同じだが，上と違ってこちらはPOST．
      // stateのtasksをjsonに変換してから送る．
      await axios
        .post(
          process.env.VUE_APP_API_TASKS,
          JSON.stringify(context.state.tasks)
        )
        .then(() => {})
        .catch(() => {
          // 失敗したらこれを表示．
          console.log("Tasksのpostに失敗しました.");
        });
    },

    // Labelsをサーバから取得．
    // まあ，上のことがわかってるならわかるでしょ．
    async getLabels(context) {
      await axios
        .get(process.env.VUE_APP_API_LABELS)
        .then((res) => {
          context.commit("setLabels", res.data);
        })
        .catch(() => {
          console.log("Labelsのgetに失敗しました.");
        });
    },

    // taskをtasksに新しく追加する．（この関数では通信しないので，asyncはいらない）
    addTask(context, task) {
      // addTask(mutationsにある)を実行．
      context.commit("addTask", task);
      // postTasksを実行し，サーバにTasksを送信する．
      context.dispatch("postTasks");
    },

    // taskの完了ボタンを押した時の関数．
    completeTask(context, id) {
      // 指定されたidのtaskがtasks配列の何番目かを取得する
      let index = context.state.tasks.findIndex((element) => element.id === id);
      // 取得したら，そのtaskをtasksから削除
      context.commit("removeTask", index);
      // postTaskを実行し，サーバにTasksを送信する．
      context.dispatch("postTasks");
    },
  },
});
