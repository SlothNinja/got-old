<template>
  <span>Moved thief to {{value.to.card.type}} at {{row}}{{value.to.column}}.</span>
</template>

<script>
  import Text from '@/components/mixins/Text'

  export default {
    mixins: [ Text ],
    name: 'sn-log-move-thief-msg',
    props: [ 'value' ],
    computed: {
      row: function () {
        var self = this
        return self.rowText(self.value.to.row)
      }
    }
  }
</script>
