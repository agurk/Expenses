<template>
  <div class="container">
    <div class="row"><div class="col-sm-12">
        <input id="dateFrom" style="width: 100px" v-model="from" v-on:change="loadAnalysis()">
        —
        <input id="dateTo" style="width: 100px" v-model="to" v-on:change="loadAnalysis()">
        <div class="float-right"><input type="text" id="ccy" style="text-align: center; width: 80px" v-model="ccy" v-on:change="loadAnalysis()"></div>
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
            <td><div class="float-right">{{ analysis[year]['salary'] | currency }}</div></td>
            <td><div class="float-right">{{ analysis[year]['expenses'] | currency }}</div></td>
            <td><div class="float-right">{{ analysis[year]['espp'] | currency }}</div></td>
            <td><div class="float-right">{{ analysis[year]['fullIncome'] | currency }}</div></td>
            <td><div class="float-right">{{ analysis[year]['spend'] | currency }}</div></td>
            <td><div class="float-right">{{ analysis[year]['saved'] | currency }}</div></td>
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
      for (const year in this.rawAnalysis) {
        result[year] = {'salary': 0, 'expenses': 0, 'espp':0, 'fullIncome':0, 'spend':0}
        result[year]['salary'] = this.zeroOrValue(this.rawAnalysis[year].classifications, 17)
        result[year]['expenses'] = this.zeroOrValue(this.rawAnalysis[year].classifications, 12) + this.zeroOrValue(this.rawAnalysis[year].classifications, 18)
        result[year]['espp'] = this.zeroOrValue(this.rawAnalysis[year].classifications, 27)
        result[year]['fullIncome'] =  result[year]['salary'] + result[year]['expenses'] + result[year]['espp']
        result[year]['spend'] = this.rawAnalysis[year].allSpend
        result[year]['saved'] = result[year]['fullIncome'] + result[year]['spend']
        result[year]['savedPercent'] = result[year]['saved'] / result[year]['fullIncome'] * 100
      }
      return result
    }
  },
  mounted() {
    this.loadAnalysis()
    this.rawAnalysis.sort((function(a, b){return b - a}))
  }
}
</script>
<style>
</style>
