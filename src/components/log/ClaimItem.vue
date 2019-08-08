<template>
  <span>Claimed {{value.from.card.type}} card at {{row}}{{value.from.column}}.</span>
</template>

<script>
  import Text from '@/components/mixins/Text'

  export default {
    mixins: [ Text ],
    name: 'sn-log-claim-item-msg',
    props: [ 'value' ],
    computed: {
      row: function () {
        var self = this
        return self.rowText(self.value.from.row)
      }
    }
  }
</script>
