<template>
  <div id="app">
    <TaskForm />
    <br />
    <Tasks />
  </div>
</template>

<script>
import TaskForm from "./components/TaskForm.vue";
import Tasks from "./components/Tasks.vue";

export default {
  name: "App",

  components: {
    TaskForm,
    Tasks,
  },

  created() {
    this.$store.dispatch("getTasks");
    this.$store.dispatch("getLabels");
  },
};
</script>
