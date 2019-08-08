<template>
  <v-card>
    <v-card-text>
      <v-layout>
        <v-flex>
          <div><strong>Title:</strong> {{game.header.title}}</div>
          <div><strong>ID:</strong> {{game.id}}</div>
          <div><strong>Turn:</strong> {{game.header.turn}}</div>
        </v-flex>
        <v-flex class='text-xs-center'>
          <div><strong>Jewels</strong></div>
          <div class='jewel-card green'>
            <v-tooltip content-class="tooltip" v-if='game.state.jewels' left max-width='150px' color='green'>
              <space-image slot='activator' :value='game.state.jewels'></space-image>
              <span>{{tooltip(game.state.jewels)}}</span>
            </v-tooltip>
          </div>
        </v-flex>
      </v-layout>
    </v-card-text>
  </v-card>
</template>

<script>
  import SpaceImage from '@/components/board/SpaceImage'
  import Button from '@/components/player/Button'
  import Tooltip from '@/components/mixins/Tooltip'

  export default {
    mixins: [ Tooltip ],
    name: 'sn-status-panel',
    components: {
      'sn-player-btn': Button,
      'space-image': SpaceImage
    },
    props: [ 'game' ]
  }
</script>

<style scoped lang="scss">
  .jewel-card {
    height:90px;
    width:90px;
    margin: 0 auto;
    border-radius:10px;
  }
</style>
