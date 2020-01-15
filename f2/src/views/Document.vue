<template>
<div class="container">
    <div class="row">
        <div class="col-sm-2"><h2>Document</h2></div>
        <div class="col-sm-10 h2">
            <div class="float-right">
                <div class="btn-group">
                    <a class="btn btn-secondary" href="">
                        Previous
                    </a>
                    <a class="btn btn-secondary" href="">
                        Next
                    </a>
                </div>
                &nbsp;
                <button class="btn btn-secondary" v-on:click="reprocess()">Reprocess</button>
                &nbsp;
                <button class="btn btn-danger" v-on:click="del()">Delete</button>
            </div>
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
                mergeId: 0
                }},
            components: {
                            ExpenseMatch 
            },
        methods: {
                    loadDocument: function() {
                        axios.get("https://localhost:8000/documents/"+this.id)
                            .then(response => {this.document= response.data})
                    },
                    imageURL: function() {
                        return '/resources/documents/' + this.document.filename
                    },
                    del: function() {
                        axios.delete("https://localhost:8000/documents/"+this.id)
                        .then(response => { if (response.status === 200) {
                            this.document.deleted = true
                        }})
                    },
                    reprocess: function() {
                        axios.post("https://localhost:8000/processor", {"id":parseInt(this.id), "type":"document"})
                    },
                    mergeExpense: function() {
                        axios.post("https://localhost:8000/mappings/", {"expenseId":parseInt(this.mergeId), "documentId":parseInt(this.id), "confirmed": true})
                        .then(response => { if (response.status === 200) {
                            this.loadDocument()
                        }})
                    }
        },
        mounted() {
                    this.loadDocument()
        }
}
</script>
<style>
.link:hover {color: #888; cursor: pointer}
</style>
