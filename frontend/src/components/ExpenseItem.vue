<template>
  <div>
    <div class="row"
      v-bind:class="{ temporary: expense.metadata.temporary, unconfirmed: !expense.metadata.confirmed }" >
      <b-col cols="1">
        <div v-if="selectedId === expense.id || selectedId === ''">
          <div class="form-check">
            <input class="form-check-input" type="checkbox"
              v-bind:id="expense.id"
              v-on:click="$emit('select', expense.id)"/>
          </div>
        </div>
        <div v-else style="padding-bottom: 0px; padding-top: 0px">
          <b-dropdown
            split
            split-variant="outline-secondary"
            variant="secondary"
            text="merge"
            style="font-size: 0.5em; padding-top:0px; padding-bottom:0px"
            size="sm"
            v-on:click="merge()"
          >
            <b-dropdown-item size="sm" style="font-size: 0.8em;" href="#" v-on:click="mergeCommission()">as commission</b-dropdown-item>
          </b-dropdown>
        </div>
      </b-col>
      <b-col cols="5" md="4"><router-link v-bind:to="linkURL()" >{{ expense.description }}</router-link></b-col>
      <b-col cols="3" md="4"><div style="float: right">{{ expense.amount | currency(expense.currency) }}</div></b-col>
      <b-col cols="auto" md="2"  v-if="groupedby === groups.day || groupedby === groups.month || groupedby === groups.year">{{ classifications[expense.metadata.classification].description }}</b-col>
      <b-col cols="auto" v-if="groupedby === groups.classification || groupedby === groups.month || groupedby === groups.year">{{ expense.date}}</b-col>
      <b-col cols="auto" md="1" v-if="expense.documents">
        <b-button v-for="doc in expense.documents"
          size="sm"
          style="font-size: 0.7em; padding-top:0px; padding-bottom:0px; margin-right: 5px"
          v-bind:key=doc.id
          v-on:click="$emit('showdocument', doc.documentId)"
          v-b-modal.document>Receipt</b-button>
      </b-col>
    </div>
    <div class="row expense-item"
      v-bind:class="{ temporary: expense.metadata.temporary, unconfirmed: !expense.metadata.confirmed }" >
      <b-col cols="1"></b-col>
      <b-col cols="2" lg="1">
        <button v-if="!expense.metadata.confirmed" class="btn btn-outline-secondary btn-sm suggestion-btn"  v-on:click="confirmExpense(expense)">confirm</button>
      </b-col>
      <b-col cols="9" v-if="!expense.metadata.confirmed">
        &nbsp;&nbsp;
        <button v-for="suggestion in suggestions" v-bind:key="suggestion.value" class="btn btn-outline-secondary btn-sm suggestion-btn" v-on:click="useSuggestion(suggestion)">{{ getSuggestionDescription(suggestion) }}</button>  
      </b-col>
    </div>
  </div>



</template>

<script>
import axios from 'axios'
export default {
  name: 'expense-item',
  data: function() {
    return {
      suggestions: [],
    }
  },
  props: ['expense', 'groupedby', 'groups', 'selectedId', 'classifications'],
  methods: {
    confirmExpense: function(expense) {
      axios.patch(this.$backend + "/expenses/"+expense.id, {"metadata":{"confirmed":true}})
        .then(function (response) { if (response.status === 200) {
          expense.metadata.confirmed = true
        }})
    },
    linkURL: function() {
      return '/expenses/' + this.expense.id
    },
    merge: function() {
      axios({ method: 'MERGE', url: this.$backend + "/expenses/"+this.expense.id, data: {"id":this.selectedId}})
        .then(response => { if (response.status === 200) {
          this.$emit('select', 'MERGED')
        }})
    },
    mergeCommission: function() {
      axios({ method: 'MERGE', url: this.$backend + "/expenses/"+this.expense.id, data: {"id":this.selectedId, "parameters":"commission"}})
        .then(response => { if (response.status === 200) {
          this.$emit('select', 'MERGED')
        }})
    },
    getSuggestions: function() {
      axios.get(this.$backend + "/expenses/suggestions/" + this.expense.id)
        .then(response => {this.suggestions = response.data})
    },
    getSuggestionDescription: function(s) {
      return this.classifications[s.value].description
    },
    useSuggestion: function(s) {
      axios.patch(this.$backend + "/expenses/"+this.expense.id, {"metadata":{"classification":parseInt(s.value)}})
        .then(response => { if (response.status === 200) {
          this.expense.metadata.classification = parseInt(s.value)
        }})
    }
  },
  mounted() {
    if (!this.expense.metadata.confirmed) {
      this.getSuggestions()
    }
  }

}
</script>
<style>
.expense-item {
  border-bottom: 1px dashed #404040;
}

.btn-sm {
  padding-top: 0px;
  padding-bottom: 0px;
  font-size: 1.6em;
}

.suggestion-btn {
  padding-top: 0px;
  margin-top: 2px;
  margin-bottom: 5px;
  font-size: 1em
}

</style>
