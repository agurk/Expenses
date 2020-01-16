<template>
  <div class="exepense-summary">
      <div class="row"><div class="col-sm-4">
          <table class="table  table-sm">
        <thead>
            <tr>
                <th>Classification</th>
                <th><div class="float-right">Amount</div></th>
            </tr>
        </thead>
        <tr class="totalRow" v-for="(total, classif) in totals" v-bind:key="classif">
           <td v-if="classifications[classif].hidden" scope="row">{{ classifications[classif].description }}</td>
           <td v-if="classifications[classif].hidden"><div class="float-right">{{ total | currency(ccy) }}</div></td>
        </tr>
        <tfoot>
            <tr>
                <td>Total</td>
                <td><div class="float-right">{{ sumTotal | currency(ccy) }}</div></td>
            </tr>
        </tfoot>
        </table>
      </div>
      <div class="col-sm-8"><span v-html="this.graph"></span></div>
      </div>
  </div>
</template>

<script>


export default {
  name: 'expense-summary',
  props: ['ccy', 'totals', 'classifications', 'graph'],
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
