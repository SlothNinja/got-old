<template>
  <div>
    <sn-player-btn :player='player' size='small'></sn-player-btn>
    {{player.user.name}} played {{kind}} card.
  </div>
</template>

<script>
  import Button from '@/components/player/Button'
  import Player from '@/components/mixins/Player'

  export default {
    mixins: [ Player ],
    name: 'sn-log-play-card-msg',
    props: [ 'value' ],
    components: {
      'sn-player-btn': Button
    },
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
      },
      player: function () {
        var self = this
        return self.playerByPID(self.value.pid)
      }
    }
  }
</script>
