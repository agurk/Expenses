<template>
  <div class="container">
    <div class="row">
      <div class="col-sm-12 topbar">
        <div class="input-group">
          <div class="input-group-prepend">
            <button type="button" class="btn btn-outline-secondary" v-on:click="change_date('monthBack')" > &lt; </button>
          </div>
          <input type="text" class="form-control date-box" id="dateFrom" v-model="from" v-on:change="loadExpenses()">
          <input type="text" id="dateTo" class="form-control date-box" v-model="to" v-on:change="loadExpenses()">
          <div class="input-group-append">
            <button type="button" class="btn btn-outline-secondary"  v-on:click="change_date('monthForward')"> &gt; </button>
          </div>
        </div>
      </div>
    </div>
    <expense-summary v-bind:ccy="displayCCY"
      v-bind:classifications="classifications"
      v-bind:graph="svg"
      v-bind:totals="rawTotals.total.classifications"></expense-summary>
    <div class="row details-header">
      <b-col cols="4">
        <b-dropdown v-bind:text="displayCCY">
          <b-dropdown-item-button @click="displayCCY='DKK'; loadExpenses()">DKK</b-dropdown-item-button>
          <b-dropdown-item-button @click="displayCCY='GBP'; loadExpenses()">GBP</b-dropdown-item-button>
          <b-dropdown-divider></b-dropdown-divider>
          <b-form-group label="Other">
            <b-form-input
              id="customCCY"
              v-model="customCCY"
              @change="changeCCY()"
            ></b-form-input>
          </b-form-group>
        </b-dropdown>
      </b-col>
      <b-col cols="8">
        <div class="float-right">
          <button class="btn btn-secondary" v-if="connected" v-on:click="loadExpenses()">Reload</button>
          <button class="btn btn-secondary" v-else v-on:click="connect()">Connect</button>
          &nbsp;
          <b-dropdown text="Show" right>

            <b-dropdown-item-button v-bind:active="reverseOrder" @click="reverseOrder = true" >Newest First</b-dropdown-item-button>
            <b-dropdown-item-button v-bind:active="!reverseOrder" @click="reverseOrder = false" >Oldest First</b-dropdown-item-button>

            <b-dropdown-divider></b-dropdown-divider>
            <b-dropdown-group id="dropdown-group-1" header="Grouped by">
              <b-dropdown-item-button v-bind:active="groupedBy === groups.day" @click="groupedBy = groups.day" >Day</b-dropdown-item-button>
              <b-dropdown-item-button v-bind:active="groupedBy === groups.month" @click="groupedBy = groups.month" >Month</b-dropdown-item-button>
              <b-dropdown-item-button v-bind:active="groupedBy === groups.year" @click="groupedBy = groups.year" >Year</b-dropdown-item-button>
              <b-dropdown-item-button v-bind:active="groupedBy === groups.classification" @click="groupedBy = groups.classification" >Classification</b-dropdown-item-button>

              <b-dropdown-divider></b-dropdown-divider>
              <b-dropdown-item-button v-bind:active="showHidden" @click="showHidden = !showHidden" >Show Hidden</b-dropdown-item-button>
            </b-dropdown-group>

          </b-dropdown>
        </div>
      </b-col>
    </div>

    <expense-section v-for="key in sectionKeys()"
      v-bind:expenses="groupedExpenses[key]"
      v-bind:label="key"
      v-bind:groupedby="groupedBy"
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
      expenses: [],
      raw_classifications: [],
      raw_fx_rates: [],
      rawTotals: {total:{totals:{}, classifications: []}},
      svg: "",
      from: "",
      to: "",
      groups: {day: "0", month: "1", year: "2", classification: "3"},
      groupedBy: "0",
      showHidden: false,
      displayCCY: "GBP",
      reverseOrder: true,
      selectedId: "",
      modalDocument: {},
      connected: false,
      socket: 0,
      customCCY: "",
    }},
  components: {
    ExpenseSection, ExpenseSummary
  },
  methods: {
    loadExpenses: function() {
      axios.get(this.$backend + "/expenses/classifications?from=" + this.from + "&to=" + this.to)
        .then(response => {this.raw_classifications = response.data; 
          axios.get(this.$backend + "/expenses?from=" + this.from + "&to=" + this.to)
            .then(response => {this.expenses = response.data})
          axios.get(this.$backend + "/analysis/totals?from=" + this.from + "&to=" + this.to + "&currency=" + this.displayCCY + "&grouping=together&classifications=" + Object.keys(this.classifications) )
            .then(response => {this.rawTotals = response.data})
          axios.get(this.$backend + "/analysis/graph?from=" + this.from + "&to=" + this.to + "&currency=" + this.displayCCY )
            .then(response => {this.svg = response.data})
        })
    },
    change_date: function(delta) {
      if (delta === 'monthBack') {
        var fromDelta = -1
        var toDelta = 0
      } else if (delta === 'monthForward') {
        fromDelta = 1
        toDelta = 2
      }
      var d = new Date(this.from);
      var month = d.getMonth()
      var newFrom = new Date(d.getFullYear(), month+fromDelta , 1)
      newFrom = new Date(newFrom.getTime() - newFrom.getTimezoneOffset() * 60 *1000)
      this.from = newFrom.toISOString().split('T')[0]
      var newTo = new Date(d.getFullYear(), month+toDelta, 0)
      newTo = new Date(newTo.getTime() - newTo.getTimezoneOffset() * 60 * 1000)
      this.to = newTo.toISOString().split('T')[0]
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
    showdoc: function(path) {
      axios.get(this.$backend + "/documents/" + path)
        .then(response => {this.modalDocument = response.data})
    },
    sectionKeys: function() {
      if (this.reverseOrder === true ) {
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
      if (this.connected === false) {
        this.newConnect()
        this.loadExpenses()
      } else {
        this.socket.send('ping')
      }
    },
    changeCCY: function() {
      if (this.customCCY.length === 3 ) {
        this.displayCCY = this.customCCY
        this.loadExpenses()
      }
    },
    newConnect: function() {
      this.socket = new WebSocket(this.$wsBackend + "/changes/expenses");
      this.socket.onopen = () => {
        this.connected = true;
        this.socket.onmessage = ({data}) => {
          if (data == "check") {
            this.socket.send("alive")
          } else {
            this.loadExpenses();
          }
        };
      }
      this.socket.onclose = () => {
        this.connected = false
      }
    },
  },
  computed: {
    groupedExpenses: function() {
      var lookup = {};
      var key

      for (var i = 0; i < this.expenses.length; i++) {
        if ( !this.showHidden && !this.classifications[this.expenses[i].metadata.classification].hidden ) {
          continue 
        }
        if ( this.groupedBy === this.groups.classification ) {
          key = this.classifications[this.expenses[i].metadata.classification].description;
        } else if (this.groupedBy === this.groups.day ) {
          key = this.expenses[i].date;
        } else if (this.groupedBy === this.groups.month ) {
          key = this.expenses[i].date.substr(0, 7);
        } else if (this.groupedBy === this.groups.year) {
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
    this.to=lastDay.toISOString().split('T')[0]
    this.from=firstDay.toISOString().split('T')[0]

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
