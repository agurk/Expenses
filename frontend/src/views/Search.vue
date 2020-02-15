<template>
  <div class="container">
    <div class="row">
      <div class="col-sm-12 topbar">
        <div class="input-group">
          <input type="text" v-autofocus class="form-control date-box" id="query" v-model="query" v-on:change="loadExpenses()">
          <div class="input-group-prepend">
            <button type="button" class="btn btn-outline-secondary"> Search </button>
          </div>
        </div>
      </div>
    </div>
    <div class="row details-header">
      <b-col cols="12">
        <div class="float-right">
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
      v-bind:expanded="expanded"
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
import axios from 'axios'

export default {
  name: 'expenses',
  data: function() {
    return {
      expenses: [],
      raw_classifications: [],
      query: "",
      groups: {day: "0", month: "1", year: "2", classification: "3"},
      groupedBy: "1",
      showHidden: true,
      expanded: true,
      displayCCY: "GBP",
      reverseOrder: true,
      selectedId: "",
      modalDocument: {},
    }},
  components: {
    ExpenseSection
  },
  methods: {
    loadExpenses: function() {
      if (this.query === "" ) {
        this.expenses = []
      } else {
        axios.get(this.$backend + "/expenses/classifications")
          .then(response => {this.raw_classifications = response.data; 
            axios.get(this.$backend + "/expenses?search=" + this.query)
              .then(response => {this.expenses = response.data})
          })}
    },
    select: function(id) {
      if (this.selectedId === "" ) {
        this.selectedId = id
      } else if (this.selectedId === id) {
        this.selectedId = ""
      } else {
        this.selectedId = id
      }
    },
    sectionKeys: function() {
      if (this.reverseOrder === true ) {
        return Object.keys(this.groupedExpenses).sort().reverse()
      } else {
        return Object.keys(this.groupedExpenses).sort()
      }
    },
    showdoc: function(path) {
      axios.get(this.$backend + "/documents/" + path)
        .then(response => {this.modalDocument = response.data})
    },
    imageURL: function() {
      return '/resources/documents/' + this.modalDocument.filename
    },
    docURL: function() {
      return '/documents/' + this.modalDocument.id
    },
  },
  computed: {
    groupedExpenses: function() {
      var lookup = {};
      var key
      var expense
      var i = 0

      for (;(expense = this.expenses[i++]);) {
        if ( !this.showHidden && !this.classifications[expense.metadata.classification].hidden ) {
          continue 
        }
        if ( this.groupedBy === this.groups.classification ) {
          key = this.classifications[expense.metadata.classification].description;
        } else if (this.groupedBy === this.groups.day ) {
          key = expense.date;
        } else if (this.groupedBy === this.groups.month ) {
          key = expense.date.substr(0, 7);
        } else if (this.groupedBy === this.groups.year) {
          key = expense.date.substr(0, 4);
        } else {
          key = expense.date;
        }

        if (!(key in lookup)) {
          lookup[key]= []
        }
        lookup[key].push(expense)
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
    this.loadExpenses()
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
