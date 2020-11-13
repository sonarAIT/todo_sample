<template>
  <div>
    <h3>タスク一覧</h3>

    <div>
      <label>ラベルでフィルタ</label>
      <label v-for="option in filterFlagOption" v-bind:key="option.label">
        <input
          type="radio"
          v-model="filterFlag"
          v-bind:value="option.value"
          @change="changedFilterFlag"
        />{{ option.label }}
      </label>
    </div>

    <div v-if="filterFlag">
      <select v-model="nowFilter" @change="changedNowFilter">
        <option
          v-for="label in labels"
          v-bind:key="label.id"
          v-bind:value="label.id"
          v-bind:label="label.labelText"
        ></option>
      </select>
    </div>
  </div>
</template>

<script>
export default {
  name: "TodoFilter",

  props: {
    propFilter: Number,
  },

  data() {
    return {
      nowFilter: this.propFilter,
      filterFlag: false,
      filterFlagOption: [
        { value: true, label: "はい" },
        { value: false, label: "いいえ" },
      ],
    };
  },

  methods: {
    changedFilterFlag: function() {
      if (this.filterFlag) {
        this.nowFilter = 0;
      } else {
        this.nowFilter = -1;
      }
      this.$emit("change", this.nowFilter);
    },

    changedNowFilter: function() {
      this.$emit("change", this.nowFilter);
    },
  },

  computed: {
    labels: function() {
      return this.$store.state.labels;
    },
  },
};
</script>
