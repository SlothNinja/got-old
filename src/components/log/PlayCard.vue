<template>
  <li>
    Played {{kind}} card.
  </li>
</template>

<script>
  export default {
    name: 'sn-log-play-card-msg',
    props: [ 'value' ],
    computed: {
      kind: function () {
        var self = this
        switch (self.value.card.kind) {
          case 'start-lamp':
            return 'lamp'
          case 'start-camel':
            return 'camel'
          default:
            return self.value.card.kind
        }
      }
    }
  }
</script>
