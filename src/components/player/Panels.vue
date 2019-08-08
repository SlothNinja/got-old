<template>
  <v-card class='mt-3'>
    <v-tabs
      id='player-tabs'
      v-model="tab"
      color="green"
      grow
      dark
      slider-color="yellow"
    >
      <v-tab
        v-for="player in players"
        :key="player.id"
        :href="`#player-${player.id}`"
        ripple
      >
        {{ player.user.name }}
      </v-tab>
      <v-tab-item
        v-for="player in players"
        :key="player.id"
        :value="`player-${player.id}`"
      >
        <sn-player-panel
          :player="player"
          @show="$emit('show', player.id)"
          @pass="$emit('pass')"
        >
        </sn-player-panel>
      </v-tab-item>
    </v-tabs>
  </v-card>
</template>

<script>
  import Panel from '@/components/player/Panel'
  import Player from '@/components/mixins/Player'

  export default {
    mixins: [ Player ],
    name: 'sn-player-panels',
    components: {
      'sn-player-panel': Panel
    },
    props: [ 'value', 'players' ],
    computed: {
      tab: {
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
