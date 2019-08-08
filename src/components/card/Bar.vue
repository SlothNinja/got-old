<template>
  <v-bottom-sheet v-model="cardbar" inset hide-overlay lazy >
    <v-card>
      <v-card-text>
        <v-layout row justify-center id="cardbar">
          <div>Player: {{player.user.name}}</div>
          <div  
            class="mr-2"
            v-for="(count, kind) in cards"
            :key="kind"
            @click="canClick ? $emit('selected', kind) : null"
          >
            <sn-card-with-count
              :kind="kind"
              :count="count"
            >
            </sn-card-with-count>
          </div>
        </v-layout>
      </v-card-text>
    </v-card>
  </v-bottom-sheet>
</template>

<script>
  import WithCount from '@/components/card/WithCount'
  import Player from '@/components/mixins/Player'

  var _ = require('lodash')

  export default {
    mixins: [ Player ],
    name: 'sn-card-bar',
    components: {
      'sn-card-with-count': WithCount
    },
    props: [ 'value', 'phase', 'player' ],
    computed: {
      cards: function () {
        return _.countBy(this.player.hand, function (card) {
          if (card.facing) {
            return card.kind
          }
          return 'card-back'
        })
      },
      canClick: function () {
        const playCardPhase = 4
        return (this.phase === playCardPhase) && (this.isPlayerFor(this.player, this.cu))
      },
      cu: {
        get: function () {
          return this.$root.cu
        },
        set: function (value) {
          this.$root.cu = value
        }
      },
      cardbar: {
        get: function () {
          return this.value
        },
        set: function (value) {
          this.$emit('input', value)
        }
      }
    }
  }
</script>
