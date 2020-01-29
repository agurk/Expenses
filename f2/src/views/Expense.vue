<template>
  <div class="container">
    <div class="row">
      <div class="col-sm-5"><h3>Expense <small>{{ id }}</small></h3></div>
      <div class="col-sm-7 h2">
        <div class="float-right">
          <button class="btn btn-danger btn-secondary" type="button"  v-on:click="deleteExpense()">Delete</button>
          &nbsp;
          <button class="btn btn-secondary" type="button" v-on:click="duplicateExpense()">Save as New</button>
          &nbsp;
          <button class="btn btn-secondary" type="button" v-on:click="saveExpense()" >Save</button>
        </div>
      </div>
    </div>

    <br>

    <div class="row">
      <div class="col-sm-8">
        <div class="input-group mb-3">
          <div class="input-group-prepend">
            <span class="input-group-text field-desc">Description</span>
          </div>
          <input class="form-control" id="exDesc" text="text" v-model="expense.description"> 
        </div>

        <div class="input-group mb-3">
          <div class="input-group-prepend">
            <span class="input-group-text field-desc">Details</span>
          </div>
          <textarea class="form-control" id="exDeetDesc" v-model="expense.detailedDescription" placeholder="none"></textarea>
        </div>
        <div class="input-group mb-3">
          <div class="input-group-prepend">
            <span class="input-group-text field-desc">Classification</span>
          </div>
          <select id="exClass" class="form-control" v-model="expense.metadata.classification">
            <option v-bind:key="key" v-bind:value="parseInt(key)" v-for="key in Object.keys(classifications)" >{{ classifications[key].description }}</option>
          </select>
        </div>
      </div>


      <div class="col-sm-4">
        <div class="row-sm-12">
          <div class="input-group">
            <span id="exCCY" class="input-group-text field-desc">{{ expense.currency }}</span>
            <input class="form-control" id="exAmount" text="number" v-model="expense.amount">
          </div>
        </div>
        <div class="row-sm-12">
          <div class="input-group">
            <span class="input-group-text field-desc">Date</span>
            <input class="form-control" id="exDate" text="text" v-model="expense.date" v-on:keydown="cursorDate(event)">
          </div>
        </div>
        <div class="row-sm-12">
          <div class="input-group">
            <span class="input-group-text field-desc">Process Date</span>
            <input class="form-control" id="procDate" text="text" v-model="expense.processDate">
          </div>
        </div>


        <br>
        <div class="row-sm-12">
          <div class="input-group">
            <span class="input-group-text field-desc">FX amount</span>
            <input class="form-control" id="exFXAmount" text="text" v-model="expense.fx.amount">
          </div>
        </div>
        <div class="row-sm-12">
          <div class="input-group">
            <span class="input-group-text field-desc">FX currency</span>
            <input class="form-control" id="exFXCCY" text="text" v-model="expense.fx.currency">
          </div>
        </div>
        <div class="row-sm-12">
          <div class="input-group">
            <span class="input-group-text field-desc">FX Rate</span>
            <input class="form-control" id="exFXRate" text="text" v-model="expense.fx.rate">
          </div>
        </div>
        <div class="row-sm-12">
          <div class="input-group">
            <span class="input-group-text field-desc">Commission</span>
            <input class="form-control" id="exCommission" text="text" v-model="expense.commission">
          </div>
        </div>
      </div>
    </div>
    <div class="row">
      <div class="col-sm-12">
        Splitwise Integration
        <external-record v-for="record in expense.externalRecords"
          v-bind:record="record"
          v-bind:key="record.reference">
        </external-record>
        <external-record></external-record>
      </div>
    </div>
  </div>

</template>

<script>
import axios from 'axios'
import ExternalRecord from '@/components/ExpenseExternalRecord.vue'

export default {
  name: 'expenses',
  props:  {
    id: { type: String }
  },
  components: { ExternalRecord },
  data: function() {return {
    expense: {metadata: {}, fx: {}},
    raw_classifications: []
  }},
  methods: {
    loadExpense: function() {
      axios.get(this.$backend + "/expenses/"+this.id)
        .then(response => {this.expense = response.data})
    },
    loadClassifications: function() {
      axios.get(this.$backend + "/expense_classifications?date="+this.expense.date)
        .then(response => {this.raw_classifications= response.data})
    },
    saveExpense: function() {
      axios.put(this.$backend + "/expenses/"+this.id, this.expense)
    },
    duplicateExpense: function() {
      axios.post(this.$backend + "/expenses/"+this.id, this.expense)
    },
    deleteExpense: function() {
      axios.delete(this.$backend + "/expenses/"+this.id)
    },
    cursorDate: function(e) {
      var d = new Date(this.expense.date);

      e = e || window.event;
      switch (e.keyCode) {
        case 38:
          d.setDate(d.getDate() + 1);
          break;
        case 40:
          d.setDate(d.getDate() -1);
          break;
      }
      this.expense.date  = (d.toISOString().slice(0,10));
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
.field-desc {min-width:125px}
</style>
