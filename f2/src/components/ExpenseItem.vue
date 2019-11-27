<template>
  <div class="row expense-item"
    v-bind:class="{ temporary: expense.Metadata.Temporary, unconfirmed: !expense.Metadata.Confirmed }" >
  <div v-if="expense.Metadata.Confirmed" class="col-sm-1"></div>
  <div v-else class="col-sm-1 link" v-on:click="confirmExpense(expense)">con</div>
  <div class="col-sm-4"> <router-link v-bind:to="{ path: linkURL() }" >{{ expense.Description }}</router-link></div>
  <div class="col-sm-2">{{ expense.Amount}} ({{expense.Currency}})</div>
  <div v-if="groupedby === groups.day || groupedby === groups.month || groupedby === groups.year" class="col-sm-2">{{ classifications[expense.Metadata.Classification].Description}}</div>
  <div v-if="groupedby === groups.classification || groupedby === groups.month || groupedby === groups.year" class="col-sm-2">{{ expense.Date}}</div>
  </div>
</template>

<script>
import axios from 'axios'
export default {
  name: 'expense-item',
  props: ['expense', 'classifications', 'groupedby', 'groups'],
   methods: {
    confirmExpense: function(expense) {
      axios.patch("https://localhost:8000/expenses/"+expense.ID, {"Metadata":{"Confirmed":true}})
        .then(function (response) { if (response.status === 200) {
          expense.Metadata.Confirmed = true
        }})
    },
    linkURL: function() {
        return '/expense/' + this.expense.ID
    }
  }

}
</script>
<style>
.expense-item {
    border-bottom: 1px dashed #404040;
}

</style>
