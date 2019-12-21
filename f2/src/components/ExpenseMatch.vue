<template>
  <div class="row expense-item">
  <div class="col-sm-1">
 <div v-if="confirmed" >C</div>
  </div>
  <div class="col-sm-1">
    <div >X</div>
  </div>
  <div class="col-sm-8"> <router-link v-bind:to="linkURL()" >{{ expense.description }}</router-link></div>
  <div class="col-sm-2">{{ expense.date }}</div>
  </div>
</template>

<script>
import axios from 'axios'
export default {
  name: 'expense-match',
  props: ['id', 'confirmed'],
  data: function() { return {
      expense: []
  }},
   methods: {
    confirmExpense: function(expense) {
      axios.patch("https://localhost:8000/expenses/"+expense.id, {"metadata":{"confirmed":true}})
        .then(function (response) { if (response.status === 200) {
          expense.metadata.confirmed = true
        }})
    },
    linkURL: function() {
        return '/expenses/' + this.id
    },
    getExpense: function() {
        axios.get("https://localhost:8000/expenses/"+this.id)
            .then(response => {this.expense= response.data})
    }
  },
  mounted() {
      this.getExpense()
  }

}
</script>
<style>
.expense-item {
    border-bottom: 1px dashed #404040;
}

</style>
