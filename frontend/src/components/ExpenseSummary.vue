<template>
  <div class="exepense-summary">
    <div class="row">
    </div>
    <div class="row"><div class="col-sm-3">
        <b-table small :items="displayTotals" :fields="totalsFields" :sort-by.sync="sortBy">
          <template v-slot:cell(amount)="data">
            <div class="float-right">
              {{ data.item.amount | currency(ccy) }}
            </div>
          </template>
        </b-table>
        <b-table small :fields="totalFields">
          <template v-slot:head()="data">
            <div v-if="data.label !== 'Total'" class="float-right">
              {{ data.label | currency(ccy) }}
            </div>
          </template>
        </b-table>
      </div>
      <div class="col-sm-9"><span v-html="this.graph"></span></div>
    </div>
  </div>
</template>

<script>


export default {
  name: 'expense-summary',
  props: ['ccy', 'totals', 'classifications', 'graph'],
  data: function() {
    return {
      totalsFields: [{ key: 'classification', sortable: true},
        {key: 'amount', sortable: true}],
      sortBy: 'amount',
    }
  },
  components: {},
  computed: {
    sumTotal: function() {
      var totes =  0
      for (var key in this.totals) {
        if (this.classifications[key].hidden === true  ) {
          totes += this.totals[key]
        }
      }
      return totes
    },
    displayTotals: function() {
      var result = []
      for (var key in this.totals) {
        if (this.classifications[key].hidden === true  ) {
          var line = { 'classification': this.classifications[key].description, 'amount': this.totals[key] }
          result.push(line)
        }
      }
      return result
    },
    totalFields: function() {
      var totes =  0
      for (var key in this.totals) {
        if (this.classifications[key].hidden === true  ) {
          totes += this.totals[key]
        }
      }
      return ['Total', ""+totes ]
    }
  },
  mounted() {
  }
}
</script>

<style>
.totalRow {
  line-height: 12px;
}
