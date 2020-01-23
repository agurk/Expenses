<template>
  <div class="exepense-section">
      <div class="row ex-header">
          <div class="col-sm-12 expense-section">{{ label }}</div>
      </div>
      <div v-if="expanded">
          <expense-item v-for="expense in expenses"
              v-bind:expense="expense"
              v-bind:groupedby="groupedby"
              v-bind:classifications="classifications"
              v-bind:groups="groups"
              v-bind:selectedId="selectedId"
              v-on:select="$emit('select', $event)"
              v-bind:key="expense.id"></expense-item>
      </div>
  </div>
</template>

<script>

import ExpenseItem from '@/components/ExpenseItem.vue'

export default {
  name: 'expense-section',
  data: function() { return{} },
  props: ['expenses', 'label', 'classifications', 'groupedby', 'groups', 'expanded', 'selectedId'],
  components: { ExpenseItem },
  methods: {}
}
</script>

<style>
.exepense-section {
    padding-bottom: 25px;
}
.ex-header {
    font-weight: bold;
    border-bottom: 2px solid #404040;
}
.unconfirmed {
    background: #E5EAF1;
}
.temporary {
    font-style: italic;
}
.link:hover {
    color: #888;
    cursor: pointer
}
</style>

