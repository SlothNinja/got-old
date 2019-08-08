<template>
  <div>
    <sn-player-btn :player='player' size='small'></sn-player-btn>
    {{player.user.name}} placed thief on {{value.area.card.kind}} at {{row}}{{value.area.column}}.
  </div>
</template>

<script>
  import Text from '@/components/mixins/Text'
  import Player from '@/components/mixins/Player'
  import Button from '@/components/player/Button'

  export default {
    mixins: [ Text, Player ],
    name: 'sn-log-place-thief-msg',
    components: {
      'sn-player-btn': Button
    },
    props: [ 'value' ],
    computed: {
      row: function () {
        var self = this
        return self.rowText(self.value.area.row)
      },
      player: function () {
        var self = this
        return self.playerByPID(self.value.pid)
      }
    }
  }
</script>
