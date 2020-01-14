<template>
<div class="container">
    <div class="row">
    <div class="col-sm-12">
        <table id="classifications_table" class="table table-hover table-sm">
        <thead>
            <tr>
                <th>Description</th>
                <th>Valid From</th>
                <th>Valid To</th>
                <th>Is Expense</th>
            </tr>
        </thead>
        Expenses
        <tr v-for="classification in expenses" v-bind:key="classification.id">
               <td scope="row"><input class="form-control" v-model="classification.description"></td>  
               <td scope="row"><input class="form-control" v-model="classification.from"></td>  
               <td scope="row"><input class="form-control" v-model="classification.to"></td>  
               <td scope="row"><input class="form-control" type="checkbox" v-model="classification.hidden"></td>  
               <td scope="row" v-on:click="saveClassification(classification)">save</td>  
        </tr>
        Non Expenses
        <tr v-for="classification in nonExpenses" v-bind:key="classification.id">
               <td scope="row"><input class="form-control" v-model="classification.description"></td>  
               <td scope="row"><input class="form-control" v-model="classification.from"></td>  
               <td scope="row"><input class="form-control" v-model="classification.to"></td>  
               <td scope="row"><input class="form-control" type="checkbox" v-model="classification.hidden"></td>  
               <td scope="row" v-on:click="saveClassification(classification)">save</td>  
        </tr>
        Old
        <tr v-for="classification in oldClassifications" v-bind:key="classification.id">
               <td scope="row"><input class="form-control" v-model="classification.description"></td>  
               <td scope="row"><input class="form-control" v-model="classification.from"></td>  
               <td scope="row"><input class="form-control" v-model="classification.to"></td>  
               <td scope="row"><input class="form-control" type="checkbox" v-model="classification.hidden"></td>  
               <td scope="row" v-on:click="saveClassification(classification)">save</td>  
        </tr>
        New
        <tr>
               <td scope="row"><input v-model="newClassification.description" class="form-control"></td>  
               <td scope="row"><input v-model="newClassification.from" class="form-control"></td>  
               <td scope="row"><input v-model="newClassification.to" class="form-control"></td>  
               <td scope="row"><input v-model="newClassification.hidden" class="form-control" type="checkbox"></td>  
               <td scope="row" v-on:click="addClassification()">new</td>  
        </tr>
        </table>
    </div>
    </div>
</div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'config',
  data: function() {
      return {
          rawClassifications: [],
          newClassification: {}
      }},
  components: {
  },
  methods: {
      loadClassifications: function() {
        axios.get("https://localhost:8000/expense_classifications")
          .then(response => {this.rawClassifications = response.data;


          })
      },
      saveClassification: function(classification) {
          axios.put("https://localhost:8000/expense_classifications/"+classification.id, classification)
      },
      addClassification: function() {
          axios.post("https://localhost:8000/expense_classifications/", this.newClassification)
          .then(this.loadClassifications, this.newClassification={})
      }
  },
  computed: {
      expenses: function() {
       var result = {}
       var today = new Date() 
          for (var classification, i = 0; (classification = this.rawClassifications[i++]);) {
              var d = new Date(classification.to)
              if (( classification.to.length === 0 || today <= d ) && classification.hidden) {
                  result[parseInt(classification.id)] = classification
              }
            }
            return result
      },
      nonExpenses: function() {
       var result = {}
       var today = new Date() 
          for (var classification, i = 0; (classification = this.rawClassifications[i++]);) {
              var d = new Date(classification.to)
              if (( classification.to.length === 0 || today <= d ) && !classification.hidden) {
                  result[parseInt(classification.id)] = classification
              }
            }
            return result
      },
      oldClassifications: function() {
       var result = {}
       var today = new Date() 
          for (var classification, i = 0; (classification = this.rawClassifications[i++]);) {
              var d = new Date(classification.to)
              if ( classification.to.length !== 0 && today > d ) {
                  result[parseInt(classification.id)] = classification
              }
            }
            return result
      }
  },
  mounted() {
      this.loadClassifications()
  }
}
</script>
<style>
</style>
