<template>
  <div class="row expense-item"
    v-bind:class="{ temporary: expense.metadata.temporary, unconfirmed: !expense.metadata.confirmed }" >
  <div v-if="expense.metadata.confirmed" class="col-sm-1"></div>
  <div v-else class="col-sm-1 link" v-on:click="confirmExpense(expense)">con</div>
  <div class="col-sm-4"> <router-link v-bind:to="linkURL()" >{{ expense.description }}</router-link></div>
  <div class="col-sm-2">{{ expense.amount}} ({{expense.currency}})</div>
  <div v-if="groupedby === groups.day || groupedby === groups.month || groupedby === groups.year" class="col-sm-2">{{ classifications[expense.metadata.classification].description}}</div>
  <div v-if="groupedby === groups.classification || groupedby === groups.month || groupedby === groups.year" class="col-sm-2">{{ expense.date}}</div>
  <div v-if="expense.documents" class="col-sm-1">
    <router-link v-for="doc in expense.documents" v-bind:key=doc.id v-bind:to="docURL(doc)">R </router-link>
  </div>
  </div>
</template>

<script>
import axios from 'axios'
export default {
  name: 'expense-item',
  props: ['expense', 'classifications', 'groupedby', 'groups'],
   methods: {
    confirmExpense: function(expense) {
      axios.patch("https://localhost:8000/expenses/"+expense.id, {"metadata":{"confirmed":true}})
        .then(function (response) { if (response.status === 200) {
          expense.metadata.confirmed = true
        }})
    },
    linkURL: function() {
        return '/expenses/' + this.expense.id
    },
    docURL: function(doc) {
        return '/documents/' + doc.documentId
    }
  }

}
</script>
<style>
.expense-item {
    border-bottom: 1px dashed #404040;
}

</style>
