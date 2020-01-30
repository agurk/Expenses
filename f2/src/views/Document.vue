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
        <b-form-checkbox v-model="document.starred" v-on:click="saveStarred()">Starred</b-form-checkbox>
      </div>
      <div class="col-sm-2">
        <b-form-checkbox v-model="document.archived" v-on:click="saveArchived()">Archived</b-form-checkbox>
      </div>
    </div>

    <div class="row">
      <div class="col-sm-6">
        <img class="img-fluid" alt="document image" :src="imageURL()">
      </div>
      <div class="col-sm-6">
        <expense-match v-for="expense in document.expenses"
          v-bind:key="expense.id"
          v-bind:id="expense.expenseId"
          v-bind:mapId="expense.id"
          v-on:del="loadDocument()"
          v-bind:confirmed="expense.confirmed"></expense-match>

        <div class="row">
          <div class="input-group">
            <span class="input-group-text">Expense ID</span>
            <input class="form-control" text="text" v-model="mergeId">
            <button class="btn btn-secondary" type="button" v-on:click="mergeExpense()">Attach</button>
          </div>
        </div>

        <textarea class="form-control" style="height: 100%" v-model="document.text"></textarea>
      </div>
    </div>
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
    mergeId: "" 
  }},
  components: {
    ExpenseMatch 
  },
  methods: {
    loadDocument: function() {
      axios.get(this.$backend + "/documents/"+this.id)
        .then(response => {this.document= response.data})
    },
    imageURL: function() {
      return '/resources/documents/' + this.document.filename
    },
    del: function() {
      axios.delete(this.$backend + "/documents/"+this.id)
        .then(response => { if (response.status === 200) {
          this.document.deleted = true
        }})
    },
    reprocess: function() {
      axios.post(this.$backend + "/processor", {"id":parseInt(this.id), "type":"document"})
    },
    mergeExpense: function() {
      axios.post(this.$backend + "/mappings/", {"expenseId":parseInt(this.mergeId), "documentId":parseInt(this.id), "confirmed": true})
        .then(response => { if (response.status === 200) {
          this.loadDocument()
        }})
    },
    saveStarred: function() {
      axios.patch(this.$backend + "/documents/"+this.document.id, {starred: !this.document.starred})
    },
    saveArchived: function() {
      axios.patch(this.$backend + "/documents/"+this.document.id, {archived: !this.document.archived})
    },
  },
  mounted() {
    this.loadDocument()
  }
}
</script>
<style>
.link:hover {color: #888; cursor: pointer}
</style>
