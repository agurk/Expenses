<template>
  <div class="document-overview">
    <div class="delete" v-on:click="deleteImage">❌</div>
    <a :href="docURL()">
      <img class="img-fluid document" hspace="20" vspace="20" alt="image" :src="imageURL()">
    </a>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'document-overview',
  props: ['doc'],
  methods: {
    imageURL: function() {
      return '/resources/documents/thumbs/' + this.doc.filename
    },
    docURL: function() {
      return '/documents/' + this.doc.id
    },
    deleteImage: function() {
        axios.delete(this.$backend + "/documents/"+this.doc.id)
        .then(response => { if (response.status === 200) {
          this.document.deleted = true
        }})
    }
  }
}
</script>

<style>
.document {
  box-shadow: 10px 10px 5px #888888;
}
.document-overview {
    position: relative
}
.delete {
    position: absolute;
    top: 0;
    right: 0;
    cursor: pointer;
}
.delete:hover {
    visibility: visible;
}
</style>
