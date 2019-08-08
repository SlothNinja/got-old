<template>
  <div>
    <div>
      Good luck:
    </div>
    <div
      v-for='(player, index) in players'
      :player='player'
      :key='index'
    >
      <sn-player-btn :player='player' size='small'></sn-player-btn>
      {{player.user.name}}
    </div>
    <div>
      Have fun.
    </div>
  </div>
</template>

<script>
  import Button from '@/components/player/Button'
  import Player from '@/components/mixins/Player'

  export default {
    name: 'sn-log-start-game-msg',
    mixins: [ Player ],
    props: [ 'value' ],
    components: {
      'sn-player-btn': Button
    },
    computed: {
      players: function () {
        var self = this
        return self.playersByPIDS(self.value.pids)
      }
    }
  }
</script>
