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
                <button class="btn btn-secondary" onclick="">Reprocess</button>
                <button class="btn btn-secondary" onclick="">Reclassify</button>
                <button class="btn btn-secondary" href="" >New Expense</button>
                <button class="btn btn-danger" onclick="">Delete</button>
            </div>
        </div>
    </div>

    <div class="row">
        <div class="col-sm-6">
        <img class="img-fluid" alt="image2" src="../assets/documents/IMG_1368.jpg">
        </div>
        <div class="col-sm-6">
        <expense-match v-for="expense in document.expenses" v-bind:key="expense.ID" v-bind:id="expense.expenseId" v-bind:confirmed="expense.confirmed"></expense-match>

            <textarea class="form-control" style="height: 100%" v-model="document.Text"></textarea>
        </div>
    </div>

    <div class="row">
    <p>{{ imageURL() }}</p>
        <img alt="image" :src="imageURL()">
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
                        return '../assets/documents/' + this.document.Filename
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
