<template>
  <div class="row expense-item">
  <div class="col-sm-1">
    <div v-if="confirmed" >C</div>
  </div>
  <div class="col-sm-1">
    <div >X</div>
  </div>
  <div class="col-sm-8"> <router-link v-bind:to="linkURL()" >{{ description }}</router-link></div>
  <div class="col-sm-2">{{ date }}</div>
  </div>
</template>

<script>
import axios from 'axios'
export default {
  name: 'expense-match',
  props: ['id', 'description', 'date', 'confirmed'],
   methods: {
    confirmExpense: function(expense) {
      axios.patch("https://localhost:8000/expenses/"+expense.ID, {"Metadata":{"Confirmed":true}})
        .then(function (response) { if (response.status === 200) {
          expense.Metadata.Confirmed = true
        }})
    },
    linkURL: function() {
        return '/expense/' + this.id
    },
  }

}
</script>
<style>
.expense-item {
    border-bottom: 1px dashed #404040;
}

</style>
