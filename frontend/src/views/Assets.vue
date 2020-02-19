<template>
  <div class="container">
    <b-row>
    </b-row>
    <b-row>
      <b-table small :items="assets" :fields="fields">
        <template v-slot:cell(amount)="row">
          {{ latestSeries(row.item).amount }}
        </template>
        <template v-slot:cell(date)="row">
          {{ latestSeries(row.item).date }}
        </template>
        <template v-slot:cell(actions)="row">
          <b-button size="sm" @click="newSeriesModal(row.item)" class="mr-2 btn-secondary" >New Amount</b-button>
          <b-button size="sm" @click="deleteAsset(row.item.id)" class="mr-1 btn-danger" >Delete</b-button>
        </template>
      </b-table>
    </b-row>

    <b-col>
      <b-row>
        <b-button v-b-modal.asset-modal>New Asset</b-button>
      </b-row>
    </b-col>

    <b-modal id="asset-modal" title="Asset" @hide="cancelAssetModal" @ok="saveAssetModal">
      <div class="input-group-prepend">
        <span class="input-group-text field-desc">Name</span>
        <input class="form-control" v-model="assetModal.name">
      </div>
      <div class="input-group-prepend">
        <span class="input-group-text field-desc">Type</span>
        <input class="form-control" v-model="assetModal.type">
      </div>
      <div class="input-group-prepend">
        <span class="input-group-text field-desc">Symbol</span>
        <input class="form-control" v-model="assetModal.symbol">
      </div>
      <div class="input-group-prepend">
        <span class="input-group-text field-desc">Reference</span>
        <input class="form-control" v-model="assetModal.reference">
      </div>
    </b-modal>

    <b-modal id="series-modal" title="Series" @hide="cancelSeriesModal" @ok="saveSeriesModal">
      <div class="input-group-prepend">
        <span class="input-group-text field-desc">Asset ID</span>
        <input class="form-control" v-model="seriesModal.assetid">
      </div>
      <div class="input-group-prepend">
        <span class="input-group-text field-desc">Date</span>
        <input class="form-control" v-model="seriesModal.date">
      </div>
      <div class="input-group-prepend">
        <span class="input-group-text field-desc">Amount</span>
        <input class="form-control" v-model="seriesModal.amount">
      </div>
    </b-modal>

        <b-modal id="fail-modal" title="Error" ok-only>
      <p class="my-4">{{ failModalText }}</p>
    </b-modal>

  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'analysis',
  data: function() {
    return {
      assets: [],
      fields: ['name', 'type', 'amount', 'date', 'actions'],
      assetModal: {name: "", type: "", symbol: "", reference: ""},
      seriesModal: {assetid: "", date: "", amount: 0},
      failModalText: "",
    }},
  components: {
  },
  methods: {
    loadAssets: function() {
      axios.get(this.$backend + "/assets")
        .then(response => {this.assets = response.data})
    },
    latestSeries: function(asset) {
      var latestDate = new Date("1970-01-01")
      var series = "none"
      asset.series.forEach( element => {
        var d = new Date(element.date)
        if (d.getTime() > latestDate.getTime() ) {
          latestDate = d
          series = element
        }

      })
      return series 
    },
    cancelAssetModal: function() {
      this.assetModal = {name: "", type: "", symbol: "", reference: ""}
    },
    saveAssetModal: function() {
      axios.post(this.$backend + "/assets", this.assetModal)
        .then(response => { if (response.status === 200) { this.loadAssets() } })
        .catch( error => { this.requestFail(error) } )
    },
    deleteAsset: function(id) {
      axios.delete(this.$backend + "/assets/" + id)
    },
    newSeriesModal: function(asset) {
      this.seriesModal.assetid = asset.id
      this.$root.$emit('bv::show::modal', 'series-modal')
    },
    cancelSeriesModal: function() {
      this.seriesModal = {assetid: "", date: "", amount: ""}
    },
    saveSeriesModal: function() {
      axios.post(this.$backend + "/assets/series", this.seriesModal)
        .then(response => { if (response.status === 200) { this.loadAssets() } })
        .catch( error=> { this.requestFail(error) } )
    },
    requestFail: function(error) {
      this.failModalText = error.response.data
      this.$root.$emit('bv::show::modal', "fail-modal")
    },
  },
  computed: {
  },
  mounted() {
    this.loadAssets()
  }
}
</script>
<style>
</style>
