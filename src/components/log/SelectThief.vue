<template>
  <span>Selected thief on {{data.area.card.type}} at {{row}}{{data.area.column}}.</span>
</template>

<script>
  import Text from '@/components/mixins/Text'

  export default {
    mixins: [ Text ],
    name: 'sn-log-select-thief-msg',
    props: [ 'data' ],
    computed: {
      row: function () {
        var self = this
        return self.rowText(self.data.area.row)
      }
    }
  }
</script>
