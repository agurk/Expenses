<template>
  <div class="exepense-summary">
      <div class="row"><div class="col-sm-4">
          <table class="table table-hover table-sm">
        <thead>
            <tr>
                <th>Classification</th>
                <th><div class="float-right">Amount</div></th>
            </tr>
        </thead>
        <tr class="totalRow" v-for="total in Object.keys(totals)" v-bind:key="total">
           <td scope="row">{{ classifications[total].description }}</td>
           <td><div class="float-right">{{ totals[total] | currency(ccy) }}</div></td>
        </tr>
        </table>
      </div></div>
  </div>
</template>

<script>


export default {
  name: 'expense-summary',
  props: ['ccy', 'totals', 'classifications'],
  components: {},
  mounted() {
      console.log (this.totals)
  }
}
</script>

<style>
.totalRow {
    line-height: 12px;
}
</style>

