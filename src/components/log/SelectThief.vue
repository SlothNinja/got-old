<template>
  <li>
    Selected thief on {{value.area.card.kind}} at {{row}}{{value.area.column}}.
  </li>
</template>

<script>
  import Text from '@/components/mixins/Text'

  export default {
    mixins: [ Text ],
    name: 'sn-log-select-thief-msg',
    props: [ 'value' ],
    computed: {
      row: function () {
        var self = this
        return self.rowText(self.value.area.row)
      }
    }
  }
</script>
