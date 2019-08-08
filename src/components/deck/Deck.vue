<template>
  <div>
    <div class='text-xs-center'>
      <v-badge color='green' dark overlap bottom >
        <span slot="badge">{{deck ? deck.length : 0}}</span>
        <sn-card :show='show' :type='topCard(deck)'></sn-card>
      </v-badge>
    </div>
    <div class='text-xs-center'><strong>{{label}}</strong></div>
  </div>
</template>

<script>
  import Card from '@/components/card/Card'

  const _ = require('lodash')

  export default {
    name: 'sn-deck',
    components: {
      'sn-card': Card
    },
    props: [ 'deck', 'show', 'label' ],
    methods: {
      topCard: function (deck) {
        var card = _.last(deck)
        return card ? card.type : 'card-back'
      }
    }
  }
</script>
