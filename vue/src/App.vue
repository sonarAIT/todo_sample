<template>
  <div id="app">
    <TaskForm />
    <br />
    <TodoFilter v-bind:propFilter="nowFilter" @change="nowFilterChanged" />
    <br />
    <Tasks v-bind:nowFilter="nowFilter" />
  </div>
</template>

<script>
import TaskForm from "./components/TaskForm.vue";
import TodoFilter from "./components/Filter.vue";
import Tasks from "./components/Tasks.vue";

export default {
  name: "App",

  components: {
    TaskForm,
    TodoFilter,
    Tasks,
  },

  data() {
    return {
      nowFilter: -1,
    };
  },

  methods: {
    nowFilterChanged: function(nowFilter) {
      this.nowFilter = nowFilter;
    },
  },

  created() {
    this.$store.dispatch("getTasks");
    this.$store.dispatch("getLabels");
  },
};
</script>
