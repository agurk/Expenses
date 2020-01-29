<template>
  <div class="container">
    <div class="row">
      <div class="col-sm-12 topbar">
        <div class="input-group">
          <input type="text" class="form-control date-box" id="query" v-model="query" v-on:change="loadExpenses()">
          <div class="input-group-prepend">
            <button type="button" class="btn btn-outline-secondary"> Search </button>
          </div>
        </div>
      </div>
    </div>
    <div class="row details-header">
      <div class="col-sm-12">
        <input id="ccy" style="width: 80px" v-model="displayCCY" v-on:change="loadExpenses()">
        <div style="float: right">
          <div class="btn-group">
            <button type="button" class="btn btn-secondary" v-bind:class="{ active : reverseOrder }"
              @click="reverseOrder = true" data-toggle="button">
              Newest
            </button>
            <button type="button" class="btn btn-secondary" v-bind:class="{ active : !reverseOrder }"
              @click="reverseOrder = false" data-toggle="button">
              Oldest 
            </button>
          </div>
          &nbsp;
          <button type="button" class="btn btn-secondary" v-bind:class="{ active : expanded }"
            aria-pressed="false" @click="expanded= !expanded" data-toggle="button">
            Details
          </button>
          &nbsp;
          <div class="btn-group btn-group-toggle" data-toggle="buttons">
            <label class="btn btn-secondary" v-bind:class="{ active : groupedBy === groups.day }">
              <input type="radio" name="options" autocomplete="off" value="0" v-model="groupedBy">Day
            </label>
            <label class="btn btn-secondary" v-bind:class="{ active : groupedBy === groups.month}">
              <input type="radio" name="options" autocomplete="off" value="1" v-model="groupedBy">Month
            </label>
            <label class="btn btn-secondary" v-bind:class="{ active : groupedBy === groups.year}">
              <input type="radio" name="options" autocomplete="off" value="2" v-model="groupedBy">Year
            </label>
            <label class="btn btn-secondary" v-bind:class="{ active : groupedBy === groups.classification}">
              <input type="radio" name="options" autocomplete="off" value="3" v-model="groupedBy">Classification
            </label>
          </div>
          &nbsp;
          <button type="button" class="btn btn-secondary" v-bind:class="{ active : showHidden }"
            aria-pressed="false" @click="showHidden = !showHidden" data-toggle="button">
            All
          </button>
        </div>
      </div>
    </div>

    <expense-section v-for="key in sectionKeys()"
      v-bind:expenses="groupedExpenses[key]"
      v-bind:label="key"
      v-bind:groupedby="groupedBy"
      v-bind:groups="groups"
      v-bind:expanded="expanded"
      v-bind:classifications="classifications"
      v-bind:selectedId="selectedId"
      v-on:select="select"
      v-bind:key="key"></expense-section>

  </div>
</template>

<script>
import ExpenseSection from '@/components/ExpenseSection.vue'
import axios from 'axios'

export default {
  name: 'expenses',
  data: function() {
    return {
      expenses: [],
      raw_classifications: [],
      query: "",
      groups: {day: "0", month: "1", year: "2", classification: "3"},
      groupedBy: "1",
      showHidden: true,
      expanded: true,
      displayCCY: "GBP",
      reverseOrder: true,
      selectedId: "",
    }},
  components: {
    ExpenseSection
  },
  methods: {
    loadExpenses: function() {
      if (this.query === "" ) {
        this.expenses = []
      } else {
        axios.get("https://localhost:8000/expense_classifications")
          .then(response => {this.raw_classifications = response.data; 
            axios.get("https://localhost:8000/expenses?search=" + this.query)
              .then(response => {this.expenses = response.data})
          })}
    },
    select: function(id) {
      if (this.selectedId === "" ) {
        this.selectedId = id
      } else if (this.selectedId === id) {
        this.selectedId = ""
      } else {
        this.selectedId = id
      }
    },
    sectionKeys: function() {
      if (this.reverseOrder === true ) {
        return Object.keys(this.groupedExpenses).sort().reverse()
      } else {
        return Object.keys(this.groupedExpenses).sort()
      }
    }
  },
  computed: {
    groupedExpenses: function() {
      var lookup = {};
      var key
      var expense
      var i = 0

      for (;(expense = this.expenses[i++]);) {
        if ( !this.showHidden && !this.classifications[expense.metadata.classification].hidden ) {
          continue 
        }
        if ( this.groupedBy === this.groups.classification ) {
          key = this.classifications[expense.metadata.classification].description;
        } else if (this.groupedBy === this.groups.day ) {
          key = expense.date;
        } else if (this.groupedBy === this.groups.month ) {
          key = expense.date.substr(0, 7);
        } else if (this.groupedBy === this.groups.year) {
          key = expense.date.substr(0, 4);
        } else {
          key = expense.date;
        }

        if (!(key in lookup)) {
          lookup[key]= []
        }
        lookup[key].push(expense)
      }
      return lookup
    },
    classifications: function() {
      var result = {}
      for (var classification, i = 0; (classification = this.raw_classifications[i++]);) {
        result[classification.id] = classification
      }
      return result
    }
  },
  mounted() {
    this.loadExpenses()
  }
}
</script>
<style>
.topbar {
  font-weight: bold;
  border-top: 2px solid #404040;
  padding-top: 5px;
  margin-top: 5px;
  padding-bottom: 5px;
  margin-bottom: 5px;
}

.details-header{
  border-bottom: 1px solid #404040;
  padding-top: 10px;
  margin-top: 10px;
  padding-bottom: 10px;
  margin-bottom: 10px;
}

.date-box {
  background-color: rgb(90, 98, 104);
  border-color: rgb(90, 98, 104);
  color: #FFFFFF;
  font-weight: bold;
  text-align: center;
}
</style>
