<template>
  <div class="container">
    <div class="row">
      <div class="col-sm-12 section-header">
        Starred
      </div>
    </div>
    <div class="row">
      <document-overview v-for="doc in starredDocuments"
        v-bind:doc="doc"
        v-bind:key="doc.id"></document-overview>
      <div class="col-sm-12"  v-if="starredDocuments.length === 0">- none -</div>
    </div>
    <div class="row">
      <div class="col-sm-12 section-header">
        Unmatched
      </div>
    </div>
    <div class="row">
      <document-overview v-for="doc in unmatchedDocuments"
        v-bind:doc="doc"
        v-bind:key="doc.id"></document-overview>
      <div class="col-sm-12"  v-if="unmatchedDocuments.length === 0">- none -</div>
    </div>
    <div class="row">
      <div class="col-sm-12 section-header">
        <b-button v-b-toggle.archivedDocs class="float-right m-1">+</b-button>
        <span class="align-bottom">Archived</span>
      </div>
    </div>
    <b-collapse id="archivedDocs">
      <div class="row">
        <document-overview v-for="doc in archivedDocuments"
          v-bind:doc="doc"
          v-bind:key="doc.id"></document-overview>
        <div class="col-sm-12"  v-if="archivedDocuments.length === 0">- none -</div>
      </div>
    </b-collapse>
  </div>
</template>

<script>
import DocumentOverview from '@/components/DocumentOverview.vue'
import axios from 'axios'

export default {
  name: 'unmatchedDocuments',
  data: function() {
    return {
      unmatchedDocuments: [],
      starredDocuments: [],
      archivedDocuments: [],
    }},
  components: {
    DocumentOverview
  },
  methods: {
    loadDocuments: function() {
      axios.get(this.$backend + "/documents?starred=true")
        .then(response => {this.starredDocuments = response.data})
      axios.get(this.$backend + "/documents?unmatched=true")
        .then(response => {this.unmatchedDocuments = response.data})
      axios.get(this.$backend + "/documents?archived=true")
        .then(response => {this.archivedDocuments= response.data})
    },
    connect() {
      this.socket = new WebSocket(this.$wsBackend + "/changes/documents");
      this.socket.onopen = () => {
        this.socket.onmessage = ({data}) => {
          if (data == "check") {
            this.socket.send("alive")
          } else {
            this.loadDocuments();
          }
        };
      };
    },
  },
  computed: {
    visibleDocuments: function() {
      var docs = []
      for (var doc, i = 0; ( doc = this.unmatchedDocuments[i++]);) {
        if (doc.expenses === null) {
          docs.unshift(doc)
        }
      }
      return docs
    }
  },
  mounted() {
    this.loadDocuments()
    this.connect()
  }
}
</script>
<style>
.section-header {
  font-weight: bold;
  border-bottom: 2px solid #404040;
}
</style>
