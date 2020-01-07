<template>
<div>
  <div class="row"
    v-bind:class="{ temporary: expense.metadata.temporary, unconfirmed: !expense.metadata.confirmed }" >
  <div class="col-sm-1">
    <div v-if="selectedId === expense.id || selectedId === ''">
        <div class="form-check">
              <input class="form-check-input" type="checkbox"
                    v-bind:id="expense.id"
                    v-on:click="$emit('select', expense.id)">
        </div>
      </div>
    </div>
  <div class="col-sm-4"> <router-link v-bind:to="linkURL()" >{{ expense.description }}</router-link></div>
  <div class="col-sm-2"><div style="float: right">{{ expense.amount | currency(expense.currency) }}</div></div>
  <div v-if="groupedby === groups.day || groupedby === groups.month || groupedby === groups.year" class="col-sm-2">{{ classifications[expense.metadata.classification].description}}</div>
  <div v-if="groupedby === groups.classification || groupedby === groups.month || groupedby === groups.year" class="col-sm-2">{{ expense.date}}</div>
  <div v-if="expense.documents" class="col-sm-1">
    <router-link v-for="doc in expense.documents" v-bind:key=doc.id v-bind:to="docURL(doc)">R </router-link>
  </div>
  </div>
  <div class="row expense-item"
    v-bind:class="{ temporary: expense.metadata.temporary, unconfirmed: !expense.metadata.confirmed }" >
    <div class="col-sm-1"></div>
    <div class="col-sm-1">
      <div v-if="!expense.metadata.confirmed" class="link" v-on:click="confirmExpense(expense)">con</div>
    </div>
    <div class="col-sm-1">
      <div v-if="selectedId !== '' && selectedId !== expense.id" class="link" v-on:click="merge()">merge</div>
    </div>
  </div>
  </div>
</template>

<script>
import axios from 'axios'
export default {
  name: 'expense-item',
  props: ['expense', 'classifications', 'groupedby', 'groups', 'selectedId'],
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
    },
    merge: function() {
        axios({ method: 'MERGE', url: "https://localhost:8000/expenses/"+this.selectedId, data: {"id":this.expense.id}})
        .then(function (response) { if (response.status === 200) {
            this.$emit('select', '')
        }})
    }
  }

}
</script>
<style>
.expense-item {
    border-bottom: 1px dashed #404040;
}

</style>
