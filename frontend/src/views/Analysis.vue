<template>
  <div class="container">
    <div class="row"><div class="col-sm-12">
        <input id="dateFrom" style="width: 100px" v-model="from" v-on:change="loadAnalysis()">
        —
        <input id="dateTo" style="width: 100px" v-model="to" v-on:change="loadAnalysis()">
        <div class="float-right">
          <b-dropdown v-bind:text="ccy">
            <b-dropdown-item-button @click="ccy='DKK'; loadAnalysis(); loadAssets()">DKK</b-dropdown-item-button>
            <b-dropdown-item-button @click="ccy='GBP'; loadAnalysis(); loadAssets()">GBP</b-dropdown-item-button>
            <b-dropdown-divider></b-dropdown-divider>
            <b-form-group label="Other">
              <b-form-input
                id="customccy"
                v-model="customccy"
                @change="changeCCY()"
              ></b-form-input>
            </b-form-group>
          </b-dropdown>

        </div>
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
    <b-row>
      <b-table small :items="assets" :fields="assetsFields" foot-clone>
        <template v-slot:cell(today)="row">
          <div class="float-right">
            {{row.item.values[0].amount | currency}}
          </div>
        </template>
        <template v-slot:cell(last_week)="row">
          <div class="float-right">
            {{row.item.values[1].amount | currency}}
          </div>
        </template>
        <template v-slot:cell(last_month)="row">
          <div class="float-right">
            {{row.item.values[2].amount | currency}}
          </div>
        </template>
        <template v-slot:cell(last_year)="row">
          <div class="float-right">
            {{row.item.values[3].amount | currency}}
          </div>
        </template>
        <template v-slot:foot()="data">
          <div v-if="data.label !== 'Name'" class="float-right">
            {{ assetTotal[assetsFields.indexOf(data.column)] | currency }}
          </div>
          <div v-else>
            Total
          </div>
        </template>
        <template v-slot:head()="data">
          <div v-if="data.label !== 'Name'" class="float-right">
            {{ data.label }}
          </div>
        </template>
      </b-table>

    </b-row>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'analysis',
  data: function() {
    return {
      rawAnalysis: [],
      assets: [],
      from: "2015-01-01",
      to: "2020-12-31",
      ccy: "DKK",
      customccy: "",
      classifications: [27, 17, 12, 18],
      assetsFields: ['name', 'today', 'last_week', 'last_month', 'last_year'],
    }},
  components: {
  },
  methods: {
    loadAnalysis: function() {
      axios.get(this.$backend + "/analysis/totals?from=" +this.from+"&to="+this.to+"&currency="+this.ccy+"&classifications="+this.classifications+"&allSpend=true")
        .then(response => {this.rawAnalysis = response.data})
    },
    loadAssets: function() {
      axios.get(this.$backend + "/analysis/assets?currency=" + this.ccy )
        .then(response => {this.assets = response.data})
    },
    zeroOrValue: function(map, key) {
      if (key in map) {
        return map[key]
      }
      return 0
    },
    changeCCY: function() {
      if (this.customccy.length === 3 ) {
        this.ccy = this.customccy
        this.loadAnalysis()
        this.loadAssets()
      }
    },
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
    },
    assetTotal: function() {
      var today = 0
      var week = 0
      var month = 0
      var year = 0
      this.assets.forEach( element => { 
        today += element.values[0].amount
        week += element.values[1].amount
        month += element.values[2].amount
        year += element.values[3].amount

      })
      return ['Total', ""+today, ""+week, ""+month, ""+year]
    }
  },
  mounted() {
    this.loadAnalysis()
    this.rawAnalysis.sort((function(a, b){return b - a}))
    this.loadAssets()
  }
}
</script>
<style>
</style>
