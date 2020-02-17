<template>
  <div class="container">
    <div class="row">
      <div class="col-sm-12 topbar">
        <div class="input-group">
          <div class="input-group-prepend">
            <button type="button" class="btn btn-outline-secondary" v-on:click="change_date('monthBack')" > &lt; </button>
          </div>
          <input type="text" class="form-control date-box" id="dateFrom" v-model="display.from" v-on:change="loadExpenses()">
          <input type="text" id="dateTo" class="form-control date-box" v-model="display.to" v-on:change="loadExpenses()">
          <div class="input-group-append">
            <button type="button" class="btn btn-outline-secondary"  v-on:click="change_date('monthForward')"> &gt; </button>
          </div>
        </div>
      </div>
    </div>
    <expense-summary v-bind:ccy="display.ccy"
      v-bind:classifications="classifications"
      v-bind:graph="svg"
      v-on:classification-select="classSelect"
      v-on:classification-deselect="classDeselect"
      v-bind:totals="rawTotals.total.classifications"></expense-summary>
    <div class="row details-header">
      <b-col cols="4">
        <b-dropdown v-bind:text="display.ccy">
          <b-dropdown-item-button @click="display.ccy='DKK'; loadExpenses()">DKK</b-dropdown-item-button>
          <b-dropdown-item-button @click="display.ccy='GBP'; loadExpenses()">GBP</b-dropdown-item-button>
          <b-dropdown-divider></b-dropdown-divider>
          <b-form-group label="Other">
            <b-form-input
              id="customccy"
              v-model="display.customccy"
              @change="changeCCY()"
            ></b-form-input>
          </b-form-group>
        </b-dropdown>
      </b-col>
      <b-col cols="8">
        <div class="float-right">
          <button class="btn btn-secondary" v-if="connection.connected" v-on:click="loadExpenses()">Refresh</button>
          <button class="btn btn-outline-danger" v-else v-on:click="connect()">Connect</button>
          &nbsp;
          <b-dropdown text="Show" right>

            <b-dropdown-item-button v-bind:active="display.reverseOrder" @click="display.reverseOrder = true" >Newest First</b-dropdown-item-button>
            <b-dropdown-item-button v-bind:active="!display.reverseOrder" @click="display.reverseOrder = false" >Oldest First</b-dropdown-item-button>

            <b-dropdown-divider></b-dropdown-divider>
            <b-dropdown-group id="dropdown-group-1" header="Grouped by">
              <b-dropdown-item-button v-bind:active="display.groupedBy === groups.day" @click="display.groupedBy = groups.day" >Day</b-dropdown-item-button>
              <b-dropdown-item-button v-bind:active="display.groupedBy === groups.month" @click="display.groupedBy = groups.month" >Month</b-dropdown-item-button>
              <b-dropdown-item-button v-bind:active="display.groupedBy === groups.year" @click="display.groupedBy = groups.year" >Year</b-dropdown-item-button>
              <b-dropdown-item-button v-bind:active="display.groupedBy === groups.classification" @click="display.groupedBy = groups.classification" >Classification</b-dropdown-item-button>

              <b-dropdown-divider></b-dropdown-divider>
              <b-dropdown-item-button v-bind:active="display.showHidden" @click="display.showHidden = !display.showHidden" >Show Hidden</b-dropdown-item-button>
            </b-dropdown-group>

          </b-dropdown>
        </div>
      </b-col>
    </div>

    <expense-section v-for="key in sectionKeys()"
      v-bind:expenses="groupedExpenses[key]"
      v-bind:label="key"
      v-bind:groupedby="display.groupedBy"
      v-bind:groups="groups"
      v-bind:classifications="classifications"
      v-bind:selectedId="selectedId"
      v-on:select="select"
      v-on:showdocument="showdoc"
      v-bind:key="key"></expense-section>

    <b-modal id="document" title="Receipt" ok-only>
      <template v-slot:modal-header>
        <h5>
          <router-link v-bind:to="docURL()">Document {{modalDocument.id}}</router-link>
        </h5>
      </template>
      <img class="img-fluid" alt="Receipt image missing" :src="imageURL()">
    </b-modal>

    <b-modal id="fail-modal" title="Error" ok-only>
      <p class="my-4">{{ failModalText }}</p>
    </b-modal>

  </div>

</template>

<script>
import ExpenseSection from '@/components/ExpenseSection.vue'
import ExpenseSummary from '@/components/ExpenseSummary.vue'
import axios from 'axios'

export default {
  name: 'expenses',
  data: function() {
    return {
      failModalText: "",
      expenses: [],
      raw_classifications: [],
      raw_fx_rates: [],
      rawTotals: {total:{totals:{}, classifications: []}},
      svg: "",
      groups: {day: "0", month: "1", year: "2", classification: "3"},
      display: {from: "", to: "", groupedBy: "0", showHidden: false, ccy: "GBP", customccy: "", reverseOrder: true},
      selectedId: "",
      selectedClassifications: {},
      selectedClassCount: 0,
      modalDocument: {},
      connection: {socket: 0, connected: false},
    }},
  components: {
    ExpenseSection, ExpenseSummary
  },
  methods: {
    loadExpenses: function() {
      axios.get(this.$backend + "/expenses/classifications?from=" + this.display.from + "&to=" + this.display.to)
        .then(response => {this.raw_classifications = response.data; 
          axios.get(this.$backend + "/expenses?from=" + this.display.from + "&to=" + this.display.to)
            .then(response => {this.expenses = response.data})
          axios.get(this.$backend + "/analysis/totals?from=" + this.display.from + "&to=" + this.display.to + "&currency=" + this.display.ccy + "&grouping=together&classifications=" + Object.keys(this.classifications) )
            .then(response => {this.rawTotals = response.data})
          axios.get(this.$backend + "/analysis/graph?from=" + this.display.from + "&to=" + this.display.to + "&currency=" + this.display.ccy )
            .then(response => {this.svg = response.data})
        })
        .catch( error=> { this.requestFail(error) } )
    },
    change_date: function(delta) {
      if (delta === 'monthBack') {
        var fromDelta = -1
        var toDelta = 0
      } else if (delta === 'monthForward') {
        fromDelta = 1
        toDelta = 2
      }
      var d = new Date(this.display.from);
      var month = d.getMonth()
      var newFrom = new Date(d.getFullYear(), month+fromDelta , 1)
      newFrom = new Date(newFrom.getTime() - newFrom.getTimezoneOffset() * 60 *1000)
      this.display.from = newFrom.toISOString().split('T')[0]
      var newTo = new Date(d.getFullYear(), month+toDelta, 0)
      newTo = new Date(newTo.getTime() - newTo.getTimezoneOffset() * 60 * 1000)
      this.display.to = newTo.toISOString().split('T')[0]
      document.getElementById('dateFrom').dispatchEvent(new Event('change'))
    },
    select: function(id) {
      if (id === "MERGED") {
        this.selectedId = ""
        this.loadExpenses()
      } else if (this.selectedId === "" ) {
        this.selectedId = id
      } else if (this.selectedId === id) {
        this.selectedId = ""
      } else {
        this.selectedId = id
      }
    },
    classSelect: function(id) {
      this.$set(this.selectedClassifications, id, true)
      this.selectedClassCount++
    },
    classDeselect: function(id) {
      this.$set(this.selectedClassifications, id, false)
      this.selectedClassCount--
    },
    showdoc: function(path) {
      axios.get(this.$backend + "/documents/" + path)
        .then(response => {this.modalDocument = response.data})
    },
    sectionKeys: function() {
      if (this.display.reverseOrder === true ) {
        return Object.keys(this.groupedExpenses).sort().reverse()
      } else {
        return Object.keys(this.groupedExpenses).sort()
      }
    },
    imageURL: function() {
      return '/resources/documents/' + this.modalDocument.filename
    },
    docURL: function() {
      return '/documents/' + this.modalDocument.id
    },
    connect: function() {
      if (this.connection.connected === false) {
        this.newConnect()
        this.loadExpenses()
      } else {
        this.connection.socket.send('ping')
      }
    },
    changeCCY: function() {
      if (this.display.customccy.length === 3 ) {
        this.display.ccy = this.display.customccy
        this.loadExpenses()
      }
    },
    newConnect: function() {
      this.connection.socket = new WebSocket(this.$wsBackend + "/changes/expenses");
      this.connection.socket.onopen = () => {
        this.connection.connected = true;
        this.connection.socket.onmessage = ({data}) => {
          if (data == "check") {
            this.connection.socket.send("alive")
          } else {
            this.loadExpenses();
          }
        };
      }
      this.connection.socket.onclose = () => {
        this.connection.connected = false
      }
    },
    requestFail: function(error) {
      this.failModalText = error.response.data
      this.$root.$emit('bv::show::modal', "fail-modal")
    },
  },
  computed: {
    groupedExpenses: function() {
      var lookup = {};
      var key

      for (var i = 0; i < this.expenses.length; i++) {
        if ( !this.display.showHidden && !this.classifications[this.expenses[i].metadata.classification].hidden ) {
          continue 
        }
        if (this.selectedClassCount > 0 &&(!(this.expenses[i].metadata.classification in this.selectedClassifications) || this.selectedClassifications[this.expenses[i].metadata.classification] === false )) {
          continue
        }
        if ( this.display.groupedBy === this.groups.classification ) {
          key = this.classifications[this.expenses[i].metadata.classification].description;
        } else if (this.display.groupedBy === this.groups.day ) {
          key = this.expenses[i].date;
        } else if (this.display.groupedBy === this.groups.month ) {
          key = this.expenses[i].date.substr(0, 7);
        } else if (this.display.groupedBy === this.groups.year) {
          key = this.expenses[i].date.substr(0, 4);
        } else {
          key = this.expenses[i].date;
        }

        if (!(key in lookup)) {
          lookup[key]= []
        }
        lookup[key].push(this.expenses[i])
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
    this.display.to=lastDay.toISOString().split('T')[0]
    this.display.from=firstDay.toISOString().split('T')[0]

    this.connect()
    window.setInterval(() => {
      this.connect()
    }, 300000)
  }
}
</script>
<style>
.topbar {
  font-weight: bold;
  border-top: 2px solid #404040;
  padding-top: 5px;
  margin-top: 5px;
  padding-bottom: 5px;
  margin-bottom: 5px;
}

.details-header{
  border-bottom: 1px solid #404040;
  padding-top: 10px;
  margin-top: 10px;
  padding-bottom: 10px;
  margin-bottom: 10px;
}

.date-box {
  background-color: rgb(90, 98, 104);
  border-color: rgb(90, 98, 104);
  color: #FFFFFF;
  font-weight: bold;
  text-align: center;
}
</style>
