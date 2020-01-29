<template>
  <div class="document-overview">
    <a :href="docURL()">
      <img class="img-fluid document" hspace="20" vspace="20" alt="image" :src="imageURL()">
    </a>
  </div>
</template>

<script>


export default {
  name: 'document-overview',
  props: ['doc'],
  methods: {
    imageURL: function() {
      return '/resources/documents/thumbs/' + this.doc.filename
    },
    docURL: function() {
      return '/documents/' + this.doc.id
    }
  }
}
</script>

<style>
.document {
  box-shadow: 10px 10px 5px #888888;
}
</style>
