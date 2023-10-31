<template>
  <div class="container">
    <b-row>
      <b-col cols="4" lg="7">
        <h3>Spending</h3>
      </b-col>
      <b-col cols="8" lg="3">
        <div class="input-group" >
          <input class="form-control date-box" label-align="right" v-model="from" v-on:change="loadAnalysis()">
          <input class="form-control date-box" label-align="right" v-model="to" v-on:change="loadAnalysis()">
        </div>
      </b-col>
      <b-col>
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
      </b-col>
    </b-row>

    <b-row>
      <b-col>
        <b-table small :items="yearlySpend" :fields="yearlySpendFields" sort-by="year" :sort-desc.sync="trueVal">
          <template v-slot:head()="data">
            <div v-if="data.label !== 'Year'" class="float-right">
              {{ data.label }}
            </div>
          </template>
          <template v-slot:cell()="row">
            <div v-if="row.field.key==='year'">
              {{ row.value }}
            </div>
            <div v-else-if="row.field.key==='savedPercent'" class="float-right">
              {{ row.value.toFixed(1) }}
            </div>
            <div v-else class="float-right">
              {{ row.value | currency }}
            </div>
          </template>
        </b-table>
      </b-col>
    </b-row>

    <b-row><b-col><h3>Assets</h3></b-col></b-row>
    <b-row>
      <b-col>
        <b-table small :items="assets" :fields="assetsFields" foot-clone>
          <template v-slot:cell(name)="row">
            {{ row.item.name }}
          </template>

          <template v-slot:cell()="row">
            <div class="float-right">
              <b-tr label-align="right">
                {{row.item.values[assetsFields.indexOf(row.field.key) - 1].amount | currency}}
              </b-tr>
              <b-tr>
                <div class="float-right">
                  <div v-if="assetDiff(row.field.key, row.index) >= 0" class="asset-gain">
                    {{ assetDiff(row.field.key, row.index) | currency }}
                  </div>
                  <div v-else-if="assetDiff(row.field.key, row.index) < 0" class="asset-loss">
                    {{ assetDiff(row.field.key, row.index) | currency }}
                  </div>
                  <div v-else>
                    -
                  </div>
                </div>
              </b-tr>
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
      </b-col>
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
      to: "2022-12-31",
      ccy: "DKK",
      customccy: "",
      classifications: [27, 17, 12, 18],
      assetsFields: ['name', 'today', 'last_week', 'last_month', 'last_year'],
      yearlySpendFields: [{key: 'year', sortable: true}, 'salary', 'expenses', 'espp', {key: 'fullIncome', label: 'Income'}, 'spend', 'saved', {key: 'savedPercent', label: '% Saved'}],
      trueVal: true,
    }},
  components: {
  },
  methods: {
    loadAnalysis: function() {
      axios.get(this.$backend + "/analysis/totals?from=" +this.from+"&to="+this.to+"&currency="+this.ccy+"&classifications="+this.classifications+"&allSpend=true&years=true")
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
    assetDiff: function(key, assetIndex) {
      if (key === 'last_year') {
        return "-"
      }
      var i = this.assetsFields.indexOf(key)
      return this.assets[assetIndex].values[i - 1].amount - this.assets[assetIndex].values[i].amount 
    },
  },
  computed: {
    yearlySpend: function() {
      var result = [] 
      for (const year in this.rawAnalysis) {
        var foo ={ 'year': year,
          'salary': this.zeroOrValue(this.rawAnalysis[year].classifications, 17),
          'expenses': this.zeroOrValue(this.rawAnalysis[year].classifications, 12) + this.zeroOrValue(this.rawAnalysis[year].classifications, 18),
          'espp': this.zeroOrValue(this.rawAnalysis[year].classifications, 27) + this.zeroOrValue(this.rawAnalysis[year].classifications, 36),
          'spend': this.rawAnalysis[year].allSpend}
        foo.fullIncome =  foo.salary + foo.expenses + foo.espp
        foo.saved =  foo.fullIncome + foo.spend
        foo.savedPercent =  foo.saved / foo.fullIncome * 100
        result.push(foo)
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
    this.to = new Date().getFullYear() + "-12-31"
    this.loadAnalysis()
    this.rawAnalysis.sort((function(a, b){return b - a}))
    this.loadAssets()
  }
}
</script>
<style>
.date-box {
  background-color: rgb(90, 98, 104);
  border-color: rgb(90, 98, 104);
  color: #FFFFFF;
  font-weight: bold;
  text-align: center;
}

.asset-loss{
  font-weight: bold;
  color: #DD0000;
  font-size: small;
}

.asset-gain {
  font-size: small;
  font-weight: bold;
  color: #00DD00;
}

</style>
