<template>
  <div>
    <div v-if="record != null">
      <div class="row">
      <div class="input-group">
        <h2>Splitwise expense: {{ record.reference  }}</h2>
      </div>
      </div>
    </div>
    <div v-else >
      <div class="row">
        <div class="col-sm-12">
          <b-button v-b-toggle.newSW>New Splitwise Expense</b-button>
        </div>
      </div>
      <b-collapse id="newSW">
        <div class="row">
          <div class="col-sm-4">
            <div class="input-group">
              <div class="input-group-prepend">
                <span class="input-group-text field-desc">Group</span>
              </div>
              <select class="form-control" v-model="group">
                <option v-bind:key="key" v-bind:value="parseInt(key)" v-for="key in Object.keys(groups).reverse()" v-on:click="selectedMembers = []" >{{ groups[key].name }}</option>
              </select>
            </div>
          </div>
          <div class="col-sm-6">
            <div v-if="group >= 0">
              <b-form-checkbox-group v-model="selectedMembers">
                <div v-bind:key="key2" v-for="key2 in Object.keys(groups[group].members)">
                  <b-form-checkbox v-on:click="saveStarred()" v-bind:value="parseInt(key2)">
                    {{ groups[group].members[key2] }}</b-form-checkbox>
                </div>
              </b-form-checkbox-group>
            </div>
          </div>
          <div class="col-sm-2">
            <div class="form-check">
              <button class="btn btn-secondary" v-on:click="saveNew()">Save</button>
            </div>
          </div>
        </div>
      </b-collapse>
    </div>
  </div>
</template>

<script>

import axios from 'axios'

export default {
  name: 'external-record',
  data: function() {
    return {
      group: -1,
      selectedMembers: [],
      groups: {} 
    }
  },
  props: ['record', 'eid'],
  methods: {
    saveNew: function() {
      axios.post(this.$backend + "/expenses/externalrecords/", {eid: this.eid, group: this.group, members: this.selectedMembers})
    },
    getGroups: function() {
      axios.get(this.$backend + "/expenses/externalrecords/")
        .then(response => {this.groups = response.data})
    }
  },
  computed: {

  },
  mounted() {
    this.getGroups()
  }
}

</script>
