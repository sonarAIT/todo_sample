<template>
  <div>
    <h3>タスク一覧</h3>

    <div>
      <label>ラベルでフィルタ</label>
      <!-- v-forではいといいえの二種類のラジオボタンを表示 -->
      <!-- v-modelによってdata内のfilterFlagが変更される -->
      <!-- filterFlagに代入されるのは，v-bind:value="option.value"より，option.value -->
      <!-- ついでに，option.label（はい，いいえ）を表示 -->
      <label v-for="option in filterFlagOption" v-bind:key="option.label">
        <input
          type="radio"
          v-model="filterFlag"
          v-bind:value="option.value"
        />{{ option.label }}
      </label>
    </div>

    <!-- もし上のラジオボタンでいいえを選んでいたら，v-ifによって表示されない． -->
    <div v-if="filterFlag">
      <!-- v-modelによってnowFilterが変更される． -->
      <!-- @changeによって，値が変更されたらchangedNowFilterが実行される． -->
      <select v-model="nowFilter" @change="changedNowFilter">
        <!-- labelsの要素の数だけ選択肢を表示 -->
        <!-- v-bind:value="label.id"により，label.idをnowFilterに代入． -->
        <!-- v-bind:key="local.id"により，表示の順番を指定 -->
        <!-- v-bind:label="label.labelText"により，labelTextを表示している． -->
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
  // このコンポーネントの名前はTodoFilte．
  name: "TodoFilter",

  // 親，つまりApp.vueからpropFilterを受け取る．
  props: {
    propFilter: Number,
  },

  // 変数置き場．
  // filterFlagがtrueの時にfilterを使用
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

  // 関数置き場
  methods: {
    // 現在のフィルタが変更されたら，イベントを発行．
    // App.vueにて@changeに登録された関数が実行される．
    changedNowFilter: function() {
      this.$emit("change", this.nowFilter);
    },
  },

  // this.$store.state.labelsと書くのは長いので，（このファイル内では）labelsだけで呼べる様にする．
  computed: {
    labels: function() {
      return this.$store.state.labels;
    },
  },

  // 変数監視（変数が変更された時に実行される関数）
  watch: {
    // FilterFlagが変更されたら実行
    filterFlag: function() {
      // nowFilterが-1だと，Filterを使用しないという意味になる
      if (this.filterFlag) {
        this.nowFilter = 0;
      } else {
        this.nowFilter = -1;
      }
      // イベントを発行．
      // App.vueにて@changeに登録された関数が実行される．
      this.$emit("change", this.nowFilter);
    },
  },
};
</script>
