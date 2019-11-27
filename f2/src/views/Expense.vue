<template>
<div class="container">
    <div class="row">
        <div class="col-sm-5"><h3>Expense <small>{{ id }}</small></h3></div>
        <div class="col-sm-7 h2">
            <button class="btn btn-danger btn-secondary" type="button" >Delete</button>
            <div class="float-right">
                <button class="btn btn-secondary" type="button" v-on:click="duplicateExpense()">Save as New</button>
                <button class="btn btn-secondary" type="button" >Reprocess</button>
                <div class="btn-group">
                    <button class="btn btn-secondary" type="button" >Merge</button>
                    <button type="button" class="btn btn-secondary dropdown-toggle dropdown-toggle-split" data-toggle="dropdown">
                        <span class="sr-only">Toggle Dropdown</span>
                        </button>
                    <div class="dropdown-menu">
                        <a href="#" class="btn btn-default" type="button" >Merge as Commission</a>
                    </div>
                </div>
                <button class="btn btn-secondary" type="button" v-on:click="saveExpense()" >Save</button>
            </div>
        </div>
    </div>

  <br>

    <div class="row">
      <div class="col-sm-8">
          <div class="input-group mb-3">
              <div class="input-group-prepend">
                  <span class="input-group-text expense-addon">Description</span>
              </div>
              <input class="form-control" id="exDesc" text="text" v-model="expense.Description"> 
          </div>

          <div class="input-group mb-3">
              <div class="input-group-prepend">
                  <span class="input-group-text expense-addon">Details</span>
              </div>
              <textarea class="form-control" id="exDeetDesc" v-model="expense.DetailedDescription" placeholder="none"></textarea>
          </div>
          <div class="input-group mb-3">
              <div class="input-group-prepend">
                  <span class="input-group-text expense-addon">Classification</span>
              </div>
              <select id="exClass" class="form-control" v-model="expense.Metadata.Classification">
                  <option v-bind:key="key" v-bind:value="key" v-for="key in Object.keys(classifications)" >{{ classifications[key].Description }}</option>
              </select>
          </div>
      </div>


      <div class="col-sm-4">
          <div class="row-sm-12">
              <div class="input-group">
                  <span id="exCCY" class="input-group-text expense-addon">{{ expense.Currency }}</span>
                  <input class="form-control" id="exAmount" text="number" v-model="expense.Amount">
              </div>
          </div>
          <div class="row-sm-12">
              <div class="input-group">
                  <span class="input-group-text expense-addon">Date</span>
                  <input class="form-control" id="exDate" text="text" v-model="expense.Date" onkeydown="cursor_date(event, 'exDate')">
              </div>
          </div>


          <br>
          <div class="row-sm-12">
              <div class="input-group">
                  <span class="input-group-text expense-addon">FX amount</span>
                  <input class="form-control" id="exFXAmount" text="text" v-model="expense.FX.Amount">
              </div>
          </div>
          <div class="row-sm-12">
              <div class="input-group">
                  <span class="input-group-text expense-addon">FX currency</span>
                  <input class="form-control" id="exFXCCY" text="text" v-model="expense.FX.Currency">
              </div>
          </div>
          <div class="row-sm-12">
              <div class="input-group">
                  <span class="input-group-text expense-addon">FX Rate</span>
                  <input class="form-control" id="exFXRate" text="text" v-model="expense.FX.Rate">
              </div>
          </div>
          <div class="row-sm-12">
              <div class="input-group">
                  <span class="input-group-text expense-addon">Commission</span>
                  <input class="form-control" id="exCommission" text="text" v-model="expense.Commission">
              </div>
          </div>
      </div>
</div>
</div>

</template>

<script>
import axios from 'axios'

export default {
  name: 'expenses',
          props:  {
                id: { type: String, default: '10875'}
            },
            data: function() {return {
                expense: [],
                raw_classifications: []
                }},
        methods: {
                    loadExpense: function() {
                        axios.get("https://localhost:8000/expenses/"+this.id)
                            .then(response => {this.expense = response.data})
                    },
                    loadClassifications: function() {
                        axios.get("https://localhost:8000/expense_classifications")
                            .then(response => {this.raw_classifications= response.data})
                    },
                    saveExpense: function() {
                        axios.put("https://localhost:8000/expenses/"+this.id, this.expense)
                    },
                    duplicateExpense: function() {
                        axios.post("https://localhost:8000/expenses/"+this.id, this.expense)
                    }
        },
        computed: {
                    classifications: function() {
                        var result = {}
                        for (var classification, i = 0; (classification = this.raw_classifications[i++]);) {
                            result[classification.ID] = classification
                        }
                        return result
                    }
        },
        mounted() {
                    this.loadExpense()
                    this.loadClassifications()
        }
}
</script>
<style>
.expense-item{}
.unconfirmed{background: #E5EAF1;}
.temporary{font-style: italic;}
.expense-section{font-weight: bold; border-bottom: 2px solid #404040;}
.link:hover {color: #888; cursor: pointer}
</style>
