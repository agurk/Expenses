<template>
  <div>
    <div v-if="record != null">
      <div class="input-group">
        <input type="text" class="form-control" v-model="record.reference">
      </div>
    </div>
    <div v-else>
      nothing here
    </div>
  </div>
</template>

<script>

import axios from 'axios'

export default {
  name: 'external-record',
  data: function() {
    return {
      newRecord: {}
    }
  },
  props: ['record', 'eid'],
  methods: {
    saveNew: function() {
      axios.post(this.$backend + "/expenses/"+this.eid, this.newRecord)
    }
  }
}

</script>
