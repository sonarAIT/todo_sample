<template>
  <div>
    <h3>タスク一覧</h3>

    <div>
      <label for="filter">ラベルでフィルタ</label>
      <div id="filter">
        <input
          type="radio"
          name="filter"
          value="OnFilter"
          v-model="FilterFlagString"
        />する
        <input
          type="radio"
          name="filter"
          value="NoFilter"
          v-model="FilterFlagString"
          checked
        />しない
      </div>
    </div>

    <div v-if="FilterFlagString == 'OnFilter'">
      <select v-model="filterIndex">
        <option
          v-for="label in labels"
          v-bind:key="label.id"
          v-bind:value="label.id"
          >{{ label.labelText }}</option
        >
      </select>
    </div>

    <br />

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
  </div>
</template>

<script>
import Task from "./Task.vue";

export default {
  name: "Tasks",

  data() {
    return {
      FilterFlagString: "NoFilter",
      filterIndex: 0,
    };
  },

  components: {
    Task,
  },

  computed: {
    tasks: function() {
      if (this.FilterFlagString == "OnFilter") {
        return this.$store.state.tasks.filter(
          (task) => task.label === this.filterIndex
        );
      }
      return this.$store.state.tasks;
    },

    labels: function() {
      return this.$store.state.labels;
    },
  },
};
</script>
