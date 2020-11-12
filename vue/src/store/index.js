import Vue from "vue";
import Vuex from "vuex";
import axios from "axios";

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    tasks: [
      {
        id: 1,
        name: "aaa",
        description: "iii",
        submitTime: "uuu",
        label: 0,
      },
      {
        id: 2,
        name: "iii",
        description: "uuu",
        submitTime: "eee",
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

  mutations: {
    setTasks(state, tasks) {
      state.tasks = tasks;
    },
    setLabels(state, labelTexts) {
      state.labels = labelTexts;
    },
    addTask(state, task) {
      state.tasks.push(task);
    },
    removeTask(state, index) {
      state.tasks.splice(index, 1);
    },
  },

  actions: {
    async getTasks(context) {
      await axios
        .get(process.env.VUE_APP_API_TASKS)
        .then((res) => {
          context.commit("setTasks", res.data);
          console.log(context.state.tasks);
        })
        .catch(() => {
          console.log("Tasksのgetに失敗しました.");
        });
    },

    async postTasks(context) {
      await axios
        .post(
          process.env.VUE_APP_API_TASKS,
          JSON.stringify(context.state.tasks)
        )
        .then(() => {})
        .catch(() => {
          console.log("Tasksのpostに失敗しました.");
        });
    },

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

    addTask(context, task) {
      context.commit("addTask", task);
      context.dispatch("postTasks");
    },

    completeTask(context, id) {
      let index = context.state.tasks.findIndex((element) => element.id === id);
      context.commit("removeTask", index);
      context.dispatch("postTasks");
    },
  },
});
