<template>
  <sn-user-btn 
    :color="color"
    :user="player.user"
    :size="size"
  >
  </sn-user-btn>
</template>

<script>
  import Button from '@/components/user/Button'
  import Player from '@/components/mixins/Player'
  import CurrentUser from '@/components/mixins/CurrentUser'

  export default {
    mixins: [ Player, CurrentUser ],
    name: 'sn-player-btn',
    props: [ 'player', 'size' ],
    components: {
      'sn-user-btn': Button
    },
    computed: {
      color: function () {
        var self = this
        if (self.cuLoading) {
          return 'none'
        }
        var cuid = _.get(self.cu, 'id', false)
        if (cuid) {
          var p = self.playerByUID(self.cu.id)
          return _.get(p.colors, self.player.id - 1, 'none')
        }
        return _.get(self.cp.colors, self.player.id - 1, 'none')
      }
    }
  }
</script>
