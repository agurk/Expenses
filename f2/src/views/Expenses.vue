<template>
<div class="container">
<div style="font-weight: bold">
    <button type="button" class="btn btn-secondary" v-on:click="change_date('monthBack')" > A </button>
    <input id="dateFrom" style="width: 100px" v-model="from" v-on:change="loadExpenses()">
    —
    <input id="dateTo" style="width: 100px" v-model="to" v-on:change="loadExpenses()">
    <button type="button" class="btn btn-secondary"  v-on:click="change_date('monthForward')"> > </button>

    <div style="float: right">
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
    <expense-section v-for="key in Object.keys(groupedExpenses).sort()"
                     v-bind:expenses="groupedExpenses[key]"
                     v-bind:total="groupTotal(groupedExpenses[key])"
                     v-bind:label="key"
                     v-bind:groupedby="groupedBy"
                     v-bind:groups="groups"
                     v-bind:expanded="expanded"
                     v-bind:classifications="classifications"
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
          raw_fx_rates: [],
          from: "",
          to: "",
          groups: {day: "0", month: "1", year: "2", classification: "3"},
          groupedBy: "0",
          showHidden: false,
          expanded: true,
          displayCCY: "GBP"
            }},
  components: {
    ExpenseSection
  },
        methods: {
          loadExpenses: function() {
            axios.get("https://localhost:8000/expenses?from=" + this.from + "&to=" + this.to)
              .then(response => {this.expenses = response.data})
          },
          loadClassifications: function() {
            axios.get("https://localhost:8000/expense_classifications?from=" + this.from + "&to=" + this.to)
              .then(response => {this.raw_classifications= response.data})
          },
          groupTotal: function(expenses) {
            var total = 0;
            for (var expense, i = 0; (expense = expenses[i++]);) {
              if (expense.Currency === this.displayCCY) {
                total += expense.Amount
              }
            }
            return total
          },
          change_date: function(delta) {
            if (delta === 'monthBack') {
                var fromDelta = -1
                var toDelta = 0
            } else if (delta === 'monthForward') {
                fromDelta = 1
                toDelta = 2
            }
            var d = new Date(this.from);
            var month = d.getMonth()
            var newFrom = new Date(d.getFullYear(), month+fromDelta , 1)
            newFrom = new Date(newFrom.getTime() - newFrom.getTimezoneOffset() * 60 *1000)
            this.from = newFrom.toISOString().split('T')[0]
            var newTo = new Date(d.getFullYear(), month+toDelta, 0)
            newTo = new Date(newTo.getTime() - newTo.getTimezoneOffset() * 60 * 1000)
            this.to = newTo.toISOString().split('T')[0]
            document.getElementById('dateFrom').dispatchEvent(new Event('change'))
          }
        },
        computed: {
          groupedExpenses: function() {
            var lookup = {};
            var key

            for (var expense, i = 0; (expense = this.expenses[i++]);) {
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
          var date = new Date();
          var firstDay = new Date(date.getFullYear(), date.getMonth(), 1);
          firstDay = new Date(firstDay.getTime() - firstDay.getTimezoneOffset() * 60 *1000)
          var lastDay = new Date(date.getFullYear(), date.getMonth() + 1, 0);
          lastDay = new Date(lastDay.getTime() - lastDay.getTimezoneOffset() * 60 * 1000)
          this.to=lastDay.toISOString().split('T')[0]
          this.from=firstDay.toISOString().split('T')[0]

          this.loadExpenses()
          this.loadClassifications()
        }
}
</script>
<style>
</style>
