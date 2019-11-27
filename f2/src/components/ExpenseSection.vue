<template>
  <div class="exepense-section">
      <div class="row ex-header">
          <div class="col-sm-10 expense-section">{{ label }}</div>
          <div class="col-sm-2 expense-section">{{ total }}</div>
      </div>
      <div v-if="expanded">
          <expense-item v-for="expense in expenses"
              v-bind:expense="expense"
              v-bind:groupedby="groupedby"
              v-bind:groups="groups"
              v-bind:classifications="classifications"
              v-bind:key="expense.ID"></expense-item>
      </div>
  </div>
</template>

<script>

import ExpenseItem from '@/components/ExpenseItem.vue'

export default {
  name: 'expense-section',
  props: ['expenses', 'label', 'total', 'classifications', 'groupedby', 'groups', 'expanded'],
  components: { ExpenseItem }
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

