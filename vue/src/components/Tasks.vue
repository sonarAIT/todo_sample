<template>
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
      <Task
        v-for="taskitem in tasks"
        v-bind:key="taskitem.id"
        v-bind:task="taskitem"
      />
    </tbody>
  </table>
</template>

<script>
import Task from "./Task.vue";

export default {
  name: "Tasks",

  props: {
    nowFilter: Number,
  },

  components: {
    Task,
  },

  computed: {
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
</script>
