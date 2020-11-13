<template>
  <!-- html tableでググれ -->
  <tr>
    <th>{{ task.id }}</th>
    <td>{{ task.name }}</td>
    <td>{{ task.description }}</td>
    <td>{{ task.submitTime }}</td>
    <td>{{ labels[task.label].labelText }}</td>
    <!-- ボタンを押すとcompleteを実行 -->
    <td><button @click="complete">完了</button></td>
  </tr>
</template>

<script>
export default {
  // componentの名前はTask
  name: "Task",

  // 親，Tasks.vueからtask一つ分のデータを受け取る．
  props: {
    task: {
      id: Number,
      name: String,
      description: String,
      submitTime: String,
      label: Number,
    },
  },

  //関数置き場
  methods: {
    complete: function() {
      // store内のタスクが完了した際のactionを実行している．
      // index.jsのaction内のcompleteTaskを参照．
      this.$store.dispatch("completeTask", this.task.id);
    },
  },

  // this.$store.state.labelsと書くのは長いので，（このファイル内では）labelsだけで呼べる様にする．
  computed: {
    labels: function() {
      return this.$store.state.labels;
    },
  },
};
</script>
