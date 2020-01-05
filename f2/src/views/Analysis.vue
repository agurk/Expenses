<template>
<div class="container">
    <div class="row"><div class="col-sm-12">
        <input id="dateFrom" style="width: 100px" v-model="from" v-on:change="loadAnalysis()">
        —
        <input id="dateTo" style="width: 100px" v-model="to" v-on:change="loadAnalysis()">
        <div class="float-right"><input id="ccy" style="width: 80px" v-model="ccy" v-on:change="loadAnalysis()"></div>
    </div></div>
    <div class="row">
    <div class="col-sm-12">
        <table id="overall_expenses" class="table table-hover table-sm">
        <thead>
            <tr>
                <th>Year</th>
                <th><div class="float-right">Salary</div></th>
                <th><div class="float-right">Expenses</div></th>
                <th><div class="float-right">ESPP</div></th>
                <th><div class="float-right">Total Income</div></th>
                <th><div class="float-right">Spend</div></th>
                <th><div class="float-right">Saved</div></th>
                <th><div class="float-right">% Saved</div></th>
            </tr>
        </thead>
            <tr v-for="year in Object.keys(analysis).sort().reverse()" v-bind:key="year">
               <th scope="row">{{ year  }}</th>  
               <td><div class="float-right">{{ analysis[year]['salary'].toFixed(0) }}</div></td>
               <td><div class="float-right">{{ analysis[year]['expenses'].toFixed(0) }}</div></td>
               <td><div class="float-right">{{ analysis[year]['espp'].toFixed(0) }}</div></td>
               <td><div class="float-right">{{ analysis[year]['fullIncome'].toFixed(0) }}</div></td>
               <td><div class="float-right">{{ analysis[year]['spend'].toFixed(0) }}</div></td>
               <td><div class="float-right">{{ analysis[year]['saved'].toFixed(0) }}</div></td>
               <td><div class="float-right">{{ analysis[year]['savedPercent'].toFixed(1) }}</div></td>
            </tr>
        </table>
    </div>
    </div>
</div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'analysis',
  data: function() {
      return {
          rawAnalysis: [],
          from: "2015-01-01",
          to: "2020-12-31",
          ccy: "DKK",
          classifications: [27, 17, 12, 18]
      }},
  components: {
  },
  methods: {
      loadAnalysis: function() {
        axios.get("https://localhost:8000/analysis/totals?from=" +this.from+"&to="+this.to+"&currency="+this.ccy+"&classifications="+this.classifications+"&allSpend=true")
          .then(response => {this.rawAnalysis = response.data})
      },
      zeroOrValue: function(map, key) {
        if (key in map) {
            return map[key]
        }
        return 0
      }
  },
  computed: {
      analysis: function() {
          var result = {}
          for (var year, i = 0; (year = this.rawAnalysis[i++]);) {
              result[year.date] = {'salary': 0, 'expenses': 0, 'espp':0, 'fullIncome':0, 'spend':0}
              result[year.date]['salary'] = this.zeroOrValue(year.totals, 17)
              result[year.date]['expenses'] = this.zeroOrValue(year.totals, 12) + this.zeroOrValue(year.totals, 18)
              result[year.date]['espp'] = this.zeroOrValue(year.totals, 27)
              result[year.date]['fullIncome'] =  result[year.date]['salary'] + result[year.date]['expenses'] + result[year.date]['espp']
              result[year.date]['spend'] = year.allSpend
              result[year.date]['saved'] = result[year.date]['fullIncome'] + result[year.date]['spend']
              result[year.date]['savedPercent'] = result[year.date]['saved'] / result[year.date]['fullIncome'] * 100
          }
          return result
      }
  },
  mounted() {
      this.loadAnalysis()
      this.rawAnalysis.sort((function(a, b){return b.Date - a.Date}))
  }
}
</script>
<style>
</style>
