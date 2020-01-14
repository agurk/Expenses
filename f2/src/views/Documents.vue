<template>
<div class="container">
    <div class="row">
        <document-overview v-for="doc in documents"
                           v-bind:doc="doc"
                           v-bind:key="doc.id"></document-overview>
    </div>
</div>
</template>

<script>
import DocumentOverview from '@/components/DocumentOverview.vue'
import axios from 'axios'

export default {
  name: 'documents',
  data: function() {
      return {
          documents: [],
      }},
  components: {
      DocumentOverview
  },
  methods: {
      loadDocuments: function() {
        axios.get("https://localhost:8000/documents")
          .then(response => {this.documents = response.data})
      },
  },
  computed: {
      visibleDocuments: function() {
          var docs = []
          for (var doc, i = 0; ( doc = this.documents[i++]);) {
              if (doc.expenses === null) {
                  docs.unshift(doc)
              }
          }
          return docs
      }
  },
  mounted() {
      this.loadDocuments()
  }
}
</script>
<style>
</style>
