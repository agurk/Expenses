<template>
<div class="container">
    <div class="row">
        <div class="col-sm-5"><h3>Expense <small>{{ id }}</small></h3></div>
        <div class="col-sm-7 h2">
            <button class="btn btn-danger btn-secondary" type="button" >Delete</button>
            <div class="float-right">
                <button class="btn btn-secondary" type="button" v-on:click="duplicateExpense()">Save as New</button>
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
              <input class="form-control" id="exDesc" text="text" v-model="expense.description"> 
          </div>

          <div class="input-group mb-3">
              <div class="input-group-prepend">
                  <span class="input-group-text expense-addon">Details</span>
              </div>
              <textarea class="form-control" id="exDeetDesc" v-model="expense.detailedDescription" placeholder="none"></textarea>
          </div>
          <div class="input-group mb-3">
              <div class="input-group-prepend">
                  <span class="input-group-text expense-addon">Classification</span>
              </div>
              <select id="exClass" class="form-control" v-model="expense.metadata.classification">
                  <option v-bind:key="key" v-bind:value="parseInt(key)" v-for="key in Object.keys(classifications)" >{{ classifications[key].description }}</option>
              </select>
          </div>
      </div>


      <div class="col-sm-4">
          <div class="row-sm-12">
              <div class="input-group">
                  <span id="exCCY" class="input-group-text expense-addon">{{ expense.currency }}</span>
                  <input class="form-control" id="exAmount" text="number" v-model="expense.amount">
              </div>
          </div>
          <div class="row-sm-12">
              <div class="input-group">
                  <span class="input-group-text expense-addon">Date</span>
                  <input class="form-control" id="exDate" text="text" v-model="expense.date" v-on:keydown="cursorDate(event, 'exDate')">
              </div>
          </div>
          <div class="row-sm-12">
              <div class="input-group">
                  <span class="input-group-text expense-addon">Process Date</span>
                  <input class="form-control" id="procDate" text="text" v-model="expense.processDate" v-on:keydown="cursorDate(event, 'procDate')">
              </div>
          </div>


          <br>
          <div class="row-sm-12">
              <div class="input-group">
                  <span class="input-group-text expense-addon">FX amount</span>
                  <input class="form-control" id="exFXAmount" text="text" v-model="expense.fx.amount">
              </div>
          </div>
          <div class="row-sm-12">
              <div class="input-group">
                  <span class="input-group-text expense-addon">FX currency</span>
                  <input class="form-control" id="exFXCCY" text="text" v-model="expense.fx.currency">
              </div>
          </div>
          <div class="row-sm-12">
              <div class="input-group">
                  <span class="input-group-text expense-addon">FX Rate</span>
                  <input class="form-control" id="exFXRate" text="text" v-model="expense.fx.rate">
              </div>
          </div>
          <div class="row-sm-12">
              <div class="input-group">
                  <span class="input-group-text expense-addon">Commission</span>
                  <input class="form-control" id="exCommission" text="text" v-model="expense.commission">
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
                id: { type: String }
            },
            data: function() {return {
                expense: {metadata: {}, fx: {}},
                raw_classifications: []
                }},
        methods: {
                    loadExpense: function() {
                        axios.get("https://localhost:8000/expenses/"+this.id)
                            .then(response => {this.expense = response.data})
                    },
                    loadClassifications: function() {
                        axios.get("https://localhost:8000/expense_classifications?date="+this.expense.date)
                            .then(response => {this.raw_classifications= response.data})
                    },
                    saveExpense: function() {
                        axios.put("https://localhost:8000/expenses/"+this.id, this.expense)
                    },
                    duplicateExpense: function() {
                        axios.post("https://localhost:8000/expenses/"+this.id, this.expense)
                    },
                    cursorDate: function(e, id) {
                        var d = new Date(document.getElementById(id).value);

                        e = e || window.event;
                        switch (e.keyCode) {
                        case 38:
                            d.setDate(d.getDate() + 1);
                        break;
                        case 40:
                            d.setDate(d.getDate() -1);
                        break;
                        }
                        document.getElementById(id).value = (d.toISOString().slice(0,10));
                    }
        },
        computed: {
                    classifications: function() {
                        var result = {}
                        for (var classification, i = 0; (classification = this.raw_classifications[i++]);) {
                            result[parseInt(classification.id)] = classification
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
</style>
