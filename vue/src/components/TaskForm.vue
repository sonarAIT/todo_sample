<template>
  <div>
    <h3>タスク登録フォーム</h3>

    <div>
      <label for="name">名前</label>
      <input type="text" id="name" v-model="name" />
    </div>

    <label for="description">説明</label>
    <div>
      <textarea
        rows="4"
        cols="40"
        id="description"
        v-model="description"
      ></textarea>
    </div>

    <div>
      <label for="attribute">ラベル </label>
      <select id="attribute" v-model="label">
        <option
          v-for="label in labels"
          v-bind:key="label.id"
          v-bind:value="label.id"
          >{{ label.labelText }}</option
        >
      </select>
    </div>

    <button @click="submit">登録</button>
  </div>
</template>

<script>
export default {
  name: "TaskForm",

  data() {
    return {
      name: "",
      description: "",
      label: 0,
    };
  },

  methods: {
    submit: function() {
      let date = new Date();
      let nextID;
      if (this.$store.state.tasks.length == 0) {
        nextID = 1;
      } else {
        nextID = this.$store.state.tasks.slice(-1)[0].id + 1;
      }
      let task = {
        id: nextID,
        name: this.name,
        description: this.description,
        label: this.label,
        submitTime:
          date.toISOString().substr(0, 10) +
          " " +
          date.getHours() +
          ":" +
          date.getMinutes() +
          ":" +
          date.getSeconds(),
      };
      this.$store.dispatch("addTask", task);
    },
  },

  computed: {
    labels: function() {
      return this.$store.state.labels;
    },
  },
};
</script>
