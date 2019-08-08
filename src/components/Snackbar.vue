<template>
  <v-snackbar
    color="info"
    auto-height
    vertical
    top
    v-model="open"
  >
    <slot></slot>
    <v-btn dark flat @click.native="open = false">Close</v-btn>
  </v-snackbar>
</template>

<script>
  export default {
    name: 'sn-snackbar',
    props: [ 'value' ],
    computed: {
      open: {
        get: function () {
          var self = this
          return self.value
        },
        set: function (value) {
          var self = this
          self.$emit('input', value)
        }
      }
    }
  }
</script>
