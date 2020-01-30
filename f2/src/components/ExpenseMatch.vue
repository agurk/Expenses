<template>
  <div class="row expense-item">
    <div class="col-sm-4">
      <button type="button" class="btn btn-outline-danger btn-sm"  v-on:click="deleteMapping()">delete</button>
      &nbsp;
      <button v-if="!confirmed" type="button" class="btn btn-outline-secondary btn-sm"  v-on:click="confirmMapping()">confirm</button>
    </div>
    <div class="col-sm-7"> <router-link v-bind:to="linkURL()" >{{ expense.description }}</router-link></div>
    <div class="col-sm-1">{{ expense.date }}</div>
  </div>
</template>

<script>
import axios from 'axios'
export default {
  name: 'expense-match',
  props: ['id', 'mapId', 'confirmed'],
  data: function() { return {
    expense: []
  }},
  methods: {
    confirmMapping: function() {
      axios.patch(this.$backend + "/mappings/"+this.mapId, {"confirmed":true})
        .then(response => { if (response.status === 200) {
          this.confirmed = true
        }})
    },
    deleteMapping: function() {
      axios.delete(this.$backend + "/mappings/"+this.mapId)
        .then(response => { if (response.status === 200) {
          this.$emit('del')
        }})
    },
    linkURL: function() {
      return '/expenses/' + this.id
    },
    getExpense: function() {
      axios.get(this.$backend + "/expenses/"+this.id)
        .then(response => {this.expense= response.data})
    }
  },
  mounted() {
    this.getExpense()
  }

}
</script>
<style>
.expense-item {
  border-bottom: 1px dashed #404040;
}

</style>
