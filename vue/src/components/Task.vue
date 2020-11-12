<template>
  <tr>
    <th>{{ task.id }}</th>
    <td>{{ task.name }}</td>
    <td>{{ task.description }}</td>
    <td>{{ task.submitTime }}</td>
    <td>{{ labels[task.label].labelText }}</td>
    <td><button @click="complete">完了</button></td>
  </tr>
</template>

<script>
export default {
  name: "Task",

  props: ["task"],

  methods: {
    complete: function() {
      this.$store.dispatch("completeTask", this.task.id);
    },
  },

  computed: {
    labels: function() {
      return this.$store.state.labels;
    },
  },
};
</script>
