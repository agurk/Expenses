<template>
  <div class="container">
    <div class="row">
      <div class="col-sm-2"><h2>Document</h2></div>
      <div class="col-sm-10 h2">
        <div class="float-right">
          <button class="btn btn-secondary" v-on:click="reprocess()">Reprocess</button>
          &nbsp;
          <button class="btn btn-danger" v-on:click="del()">Delete</button>
        </div>
      </div>
    </div>

    <div class="row">
      <div class="col-sm-1">
        <b-form-checkbox v-model="document.starred" v-on:change="saveStarred()">Starred</b-form-checkbox>
      </div>
      <div class="col-sm-2">
        <b-form-checkbox v-model="document.archived" v-on:change="saveArchived()">Archived</b-form-checkbox>
      </div>
    </div>

    <div class="row">
      <div class="col-sm-6">
        <img class="img-fluid" alt="document image" :src="imageURL()">
      </div>
      <div class="col-sm-6">

        <div class="row">
          <div class="input-group">
            <span class="input-group-text">Expense ID</span>
            <input class="form-control" text="text" v-model="mergeId">
            <button class="btn btn-secondary" type="button" v-on:click="mergeExpense()">Attach</button>
          </div>
        </div>

        <div class="no-expenses" v-if="document.expenses.length === 0">No matched expenses</div>

        <expense-match v-for="expense in document.expenses"
          v-bind:key="expense.id"
          v-bind:id="expense.expenseId"
          v-bind:mapId="expense.id"
          v-on:del="loadDocument()"
          v-bind:confirmed="expense.confirmed"></expense-match>

        <div class="row"><code><pre>{{document.text}}</pre></code></div>
      </div>
    </div>

    <b-modal id="fail-modal" title="Error" ok-only>
      <p class="my-4">{{ failModalText }}</p>
    </b-modal>

  </div>

</template>

<script>
import axios from 'axios'
import ExpenseMatch from '@/components/ExpenseMatch.vue'

export default {
  name: 'expenses',
  props:  {
    id: { type: String}
  },
  data: function() { return {
    document: [],
    failModalText: "",
    mergeId: "" 
  }},
  components: {
    ExpenseMatch 
  },
  methods: {
    loadDocument: function() {
      axios.get(this.$backend + "/documents/"+this.id)
        .then(response => {this.document= response.data})
        .catch( error=> { this.requestFail(error) } )
    },
    imageURL: function() {
      return '/resources/documents/' + this.document.filename
    },
    del: function() {
      axios.delete(this.$backend + "/documents/"+this.id)
        .then(response => { if (response.status === 200) {
          this.document.deleted = true
        }})
        .catch( error=> { this.requestFail(error) } )
    },
    reprocess: function() {
      axios.post(this.$backend + "/processor", {"id":parseInt(this.id), "type":"document"})
    },
    mergeExpense: function() {
      axios.post(this.$backend + "/mappings/", {"expenseId":parseInt(this.mergeId), "documentId":parseInt(this.id), "confirmed": true})
        .then(response => { if (response.status === 200) {
          this.loadDocument()
        }})
        .catch( error=> { this.requestFail(error) } )
    },
    saveStarred: function() {
      axios.patch(this.$backend + "/documents/"+this.document.id, {starred: !this.document.starred})
    },
    saveArchived: function() {
      axios.patch(this.$backend + "/documents/"+this.document.id, {archived: !this.document.archived})
    },
    requestFail: function(error) {
      this.failModalText = error.response.data
      this.$root.$emit('bv::show::modal', "fail-modal")
    }
  },
  mounted() {
    this.loadDocument()
  }
}
</script>
<style>
.link:hover {color: #888; cursor: pointer}
.no-expenses {
    text-align: center;
    padding: 5px;
}
</style>
