<template>
  <div class="exepense-summary">
    <div class="d-none d-lg-block">
      <b-row align-h="center" >
        <b-col cols="10" lg="3">

          <b-table small :items="displayTotals" :fields="totalsFields" :sort-by.sync="sortBy" @row-clicked="totalsClicked">
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
        </b-col>

        <b-col cols="12" lg="9">
          <span v-html="this.graph"></span>
        </b-col>
      </b-row>
    </div>

    <div class="d-block d-lg-none">
      <div role="tablist">
        <b-card no-body class="mb-1">
          <b-card-header header-tag="header" class="p-1" role="tab">
            <b-button block v-b-toggle.accordion-1>Summary</b-button>
          </b-card-header>
          <b-collapse id="accordion-1" visible accordion="summary" role="tabpanel">
            <b-card-body>

              <b-table small :items="displayTotals" :fields="totalsFields" :sort-by.sync="sortBy" @row-clicked="totalsClicked">
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
            </b-card-body>
          </b-collapse>
        </b-card>


        <b-card no-body class="mb-1">
          <b-card-header header-tag="header" class="p-1" role="tab">
            <b-button block v-b-toggle.accordion-2>Graph</b-button>
          </b-card-header>
          <b-collapse id="accordion-2" visible accordion="summary" role="tabpanel">
            <b-card-body>
              <span v-html="this.graph"></span>
            </b-card-body>
          </b-collapse>
        </b-card>
      </div>
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
      selectedClassifications: {},
      selectedClassCount: 0,
    }
  },
  methods: {
    totalsClicked: function(record) {
      if (record.selected) {
        record._rowVariant=""
        this.$emit('classification-deselect', record.cid)
        this.$set(this.selectedClassifications, record.cid, false)
        this.selectedClassCount--
      } else {
        record._rowVariant="dark"
        this.$emit('classification-select', record.cid)
        this.$set(this.selectedClassifications, record.cid, true)
        this.selectedClassCount++
      }
      record.selected = !record.selected
    }
  },
  computed: {
    displayTotals: function() {
      var result = []
      for (var key in this.totals) {
        if (this.classifications[key].hidden === true ) {
          var line = { 'classification': this.classifications[key].description, 'amount': this.totals[key], '_rowVariant': '', 'cid':key, 'selected': false }
          result.push(line)
        }
      }
      return result
    },
    totalFields: function() {
      var totes =  0
      this.displayTotals.forEach( element => {
        if (this.classifications[element.cid].hidden === true) {
          if (this.selectedClassCount === 0 || this.selectedClassifications[element.cid] === true) {
            totes += this.totals[element.cid]
          }}})
      return ['Total', ""+totes ]
    }
  },
}
</script>

<style>
.totalRow {
  line-height: 12px;
}
