<template>
  <v-card>
    <v-card-text>
      <v-layout row>
        <v-flex class='text-xs-center mb-3'>
          <sn-player-btn
            :id="`panel-button-${player.user.id}`"
            :player='player'
          ></sn-player-btn>
        </v-flex>
        <v-flex class='text-xs-left'>
          <div
            @click="$router.push({ name: 'user', params: { id: player.user.id } })"
            class='mt-3'
          >
          <v-icon :color="iconColor">{{icon}}</v-icon>
          {{player.user.name}}
        </div>
          <div><strong>Score:</strong> {{player.score}}</div>
        </v-flex>
      </v-layout>
      <v-divider></v-divider>
      <v-layout row class='mt-3'>
        <v-flex>
          <deck :id='`hand-${player.id}`' label='Hand' :deck='player.hand' :show='false'></deck>
        </v-flex>
        <v-flex>
          <deck :id='`draw-${player.id}`' label='Draw' :deck='player.draw' :show='false'></deck>
        </v-flex>
      </v-layout>
      <v-layout row class='mt-3'>
        <v-flex>
          <div>
            <v-btn small v-if='isPlayerFor(player, cu)' color="green" dark @click.stop="$emit('show')">Hand</v-btn>
          </div>
          <div>
            <v-btn small v-if='isPlayerFor(player, cu) && canPass' color="green" dark @click.stop="$emit('pass')">Pass</v-btn>
          </div>
        </v-flex>
        <v-flex>
          <deck :id='`discard-${player.id}`' label='Discard' :deck='player.discard' :show='true'></deck>
        </v-flex>
      </v-layout>
    </v-card-text>
  </v-card>
</template>

<script>
  import Button from '@/components/player/Button'
  import Deck from '@/components/deck/Deck'
  import Player from '@/components/mixins/Player'
  import CurrentUser from '@/components/mixins/CurrentUser'

  export default {
    mixins: [ Player, CurrentUser ],
    name: 'sn-player-panel',
    components: {
      'sn-player-btn': Button,
      'deck': Deck
    },
    props: [ 'player' ],
    computed: {
      playIcon: function () {
        var self = this
        return self.cp && (self.cp.id === self.player.id)
      },
      icon: function () {
        var self = this
        return self.playIcon ? 'play_arrow' : 'stop'
      },
      iconColor: function () {
        var self = this
        return self.playIcon ? 'green' : 'red'
      },
      canPass: function () {
        var self = this
        return self.cp && self.cu && !self.cp.performedAction && (self.cp.user.id === self.cu.id)
      }
    }
  }
</script>
