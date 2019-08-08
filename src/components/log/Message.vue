<template>
  <component :is='component' :value='value'></component>
</template>

<script>
  import StartGame from '@/components/log/StartGame'
  import PlaceThief from '@/components/log/PlaceThief'
  import SelectThief from '@/components/log/SelectThief'
  import MoveThief from '@/components/log/MoveThief'
  import ClaimItem from '@/components/log/ClaimItem'
  import PlayCard from '@/components/log/PlayCard'
  import DrawCard from '@/components/log/DrawCard'

  export default {
    name: 'sn-log-message',
    props: [ 'value' ],
    components: {
      'sn-log-start-game-msg': StartGame,
      'sn-log-place-thief-msg': PlaceThief,
      'sn-log-select-thief-msg': SelectThief,
      'sn-log-move-thief-msg': MoveThief,
      'sn-log-claim-item-msg': ClaimItem,
      'sn-log-play-card-msg': PlayCard,
      'sn-log-draw-card-msg': DrawCard
    },
    computed: {
      component: function () {
        var self = this
        console.log(`component: ${JSON.stringify(self.value)}`)
        return `sn-log-${self.value.template}-msg`
      }
    }
  }
</script>
