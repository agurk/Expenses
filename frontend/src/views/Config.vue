<template>
  <div class="container">
    <div class="row">
      <div class="col-sm-8 section-header">
        <h3>Classifications</h3>
      </div>
      <div class="col-sm-2">
        <b-button v-b-toggle.classifications class="float-right">Show All</b-button>
      </div>
      <div class="col-sm-2 section-header">
        <b-button v-b-modal.class-modal class="float-right">New Classification</b-button>
      </div>
    </div>

    <b-collapse id="classifications">
      <div class="row">
        <b-table small hover :fields="classFields" :items="rawClassifications" :sort-by.sync="classSort">
          <template v-slot:cell(actions)="row">
            <b-button size="sm" @click="classmodal(row.item, row.index, $event.target)" class="mr-1">
              Edit
            </b-button>
          </template>
        </b-table>

        <b-modal :id="classModal.id" :title="classModal.title"  @hide="resetClassModal" @ok="classModalOk">
          <div class="input-group-prepend">
            <span class="input-group-text field-desc">Description</span>
            <input class="form-control" v-model="classModal.classification.description" >
          </div>
          <div class="input-group-prepend">
            <span class="input-group-text field-desc">Valid From</span>
            <input class="form-control" v-model="classModal.classification.from"  >
          </div>
          <div class="input-group-prepend">
            <span class="input-group-text field-desc">Valid To</span>
            <input class="form-control"  v-model="classModal.classification.to"  >
          </div>
          Expense
          <input class="form-control" type="checkbox" v-model="classModal.classification.hidden">  
        </b-modal>

      </div>
    </b-collapse>



    <div class="row">
      <div class="col-sm-8 section-header">
        <h3>Accounts</h3>
      </div>
      <div class="col-sm-2">
        <b-button v-b-toggle.accounts class="float-right">Show All</b-button>
      </div>
      <div class="col-sm-2 section-header">
        <b-button v-b-modal.account-modal class="float-right">New Account</b-button>
      </div>
    </div>

    <b-collapse id="accounts">
      <div class="row">
        <b-table small hover :fields="accountFields" :items="accounts" :sort-by.sync="accountSort">
          <template v-slot:cell(actions)="row">
            <b-button size="sm" @click="accmodal(row.item, row.index, $event.target)" class="mr-1">
              Edit
            </b-button>
          </template>

        </b-table>
        <div>&nbsp;</div>

        <b-modal :id="accountModal.id" :title="accountModal.title"  @hide="resetAccModal" @ok="accModalOk">
          <div class="input-group-prepend">
            <span class="input-group-text field-desc">Name</span>
            <input class="form-control" v-model="accountModal.account.name" >
          </div>
        </b-modal>

      </div>
    </b-collapse>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'config',
  data: function() {
    return {
      rawClassifications: [],
      newClassification: {},
      accounts: {},
      classFields: [{key:'description', sortable:true},'from',{key: 'to', sortable:true},'actions'],
      accountFields: ['name','actions'],
      classSort: 'to',
      accountSort: 'name',
      classModal: {
        id: 'class-modal',
        title: '',
        classification: {}
      },
      accountModal: {
        id: 'account-modal',
        title: '',
        account: {}
      }
    }},
  components: {
  },
  methods: {
    loadClassifications: function() {
      axios.get(this.$backend + "/expenses/classifications")
        .then(response => {this.rawClassifications = response.data;
        })
    },
    saveClassification: function(classification) {
      axios.put(this.$backend + "/expenses/classifications/"+classification.id, classification)
    },
    addClassification: function(classification) {
      axios.post(this.$backend + "/expenses/classifications/", classification)
        .then(response => { if (response.status === 200) { this.loadClassifications() } })
    },
    classmodal(item, index, button) {
      this.classModal.title = `Edit Classification`
      this.classModal.classification = item
      this.$root.$emit('bv::show::modal', this.classModal.id, button)
    },
    classModalOk: function() {
      if (this.classModal.classification.id != null ) {
        this.saveClassification(this.classModal.classification)
      } else {
        this.addClassification(this.classModal.classification)
      }
    },
    resetClassModal() {
      this.classModal.title = ''
      this.classModal.classification= {}
    },
    loadAccounts: function() {
      axios.get(this.$backend + "/expenses/accounts")
        .then(response => {this.accounts= response.data;
        })
    },
    saveAccount: function(account) {
      axios.put(this.$backend + "/expenses/accounts/"+account.id, account)
    },
    addAccount: function(account) {
      axios.post(this.$backend + "/expenses/accounts/", account)
        .then(response => { if (response.status === 200) { this.loadAccounts() } })
    },
    accmodal(item, index, button) {
      this.accountModal.title = `Edit Account`
      this.accountModal.account = item
      this.$root.$emit('bv::show::modal', this.accountModal.id, button)
    },
    resetAccModal() {
      this.accountModal.title = ''
      this.accountModal.account = {}
    },
    accModalOk: function() {
      if (this.accountModal.account.id != null) {
        this.saveAccount(this.accountModal.account)
      } else {
        this.addAccount(this.accountModal.account)
      }
    }
  },
  mounted() {
    this.loadClassifications()
    this.loadAccounts()
  }
}
</script>
<style>
</style>
