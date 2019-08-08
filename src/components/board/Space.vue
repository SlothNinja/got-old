<template>
  <v-card
    :color="value.clickable ? 'yellow' : 'green darken-4'"
    ripple
    raised
    class="ma-1"
    hover
    >
    <v-card-text class="pa-0">
      <v-tooltip
        left
        max-width='150px'
        color='green'
        open-delay='800'
        :disabled="!showCard"
      >
        <template #activator="{ on }">
          <div
            class="board-space"
            :class="value.clickable ? 'clickable' : null"
            @click="value.clickable ? $emit('selected') : null "
            v-on="on"
          >
            <transition
              name="grow"
              @leave='claimCard'
            >
              <sn-space-image
                v-if="showCard"
                :value='value.card'
              ></sn-space-image>
            </transition>
            <transition
              @enter="enter"
              @leave="leave"
              :css="false"
            >
              <sn-thief-image
                v-if='showThief'
                :id='id'
                :value='thiefColor'
              >
              </sn-thief-image>
            </transition>
          </div>
        </template>
        <span>{{tooltip(value.card.kind)}}</span>
      </v-tooltip>
    </v-card-text>
  </v-card>
</template>

<script>
  import SpaceImage from '@/components/board/SpaceImage'
  import Tooltip from '@/components/mixins/Tooltip'
  import Player from '@/components/mixins/Player'
  import CurrentUser from '@/components/mixins/CurrentUser'
  import Thief from '@/components/thief/Image'
  import Image from '@/components/card/Image'

  const Velocity = require('velocity-animate')
  const _ = require('lodash')

  export default {
    mixins: [ Tooltip, Player, CurrentUser ],
    name: 'sn-space',
    components: {
      'sn-space-image': SpaceImage,
      'sn-card-image': Image,
      'sn-thief-image': Thief
    },
    props: [ 'value' ],
    methods: {
      claimCard: function (el, done) {
        console.log('claimCard')
        var self = this
        if (!self.animate) {
          done()
          return
        }
        var origin = document.getElementById(`space-${self.value.row}-${self.value.column}`)
        var selector = `discard-${self.cp.id}`
        console.log(`selector: ${selector}`)
        var target = document.getElementById(selector)
        _.delay(function(el, target, origin, done) {
          self.moveFrom(el, target, origin, done)
        }, 750, el, target, origin, done)
      },
      enter: function (el, done) {
        var self = this
        if (!self.animate) {
          done()
          return
        }
        var origin = document.getElementById(`space-${self.value.row}-${self.value.column}`)
        if (self.value.thief.from.row == 0) {
          var panelSelector = `panel-button-${self.cp.user.id}`
          var panel = document.getElementById(panelSelector)
          self.moveTo(el, panel, origin, done)
          return
        }
        var fromSpace = document.getElementById(`space-${self.value.thief.from.row}-${self.value.thief.from.column}`)
        self.moveTo(el, fromSpace, origin, done)
      },
      leave: function (el, done) {
        var self = this
        if (!self.animate) {
          done()
          return
        }
        var origin = document.getElementById(`space-${self.value.row}-${self.value.column}`)
        if (self.value.thief.from.row == 0) {
          done()
          return
        }
        var panelSelector = `panel-button-${self.cp.user.id}`
        var panel = document.getElementById(panelSelector)
        self.moveFrom(el, panel, origin, done)
      },
      moveTo: function (from, to, origin, done) {
        var fromRect = from.getBoundingClientRect()
        var originRect = origin.getBoundingClientRect()
        var toRect = to.getBoundingClientRect()
        var stopTop = fromRect.top - originRect.top
        var stopLeft = fromRect.left - originRect.left
        var startTop = toRect.top + 25 - originRect.top
        var startLeft = toRect.left + 25- originRect.left
        var midLeft = (startLeft-stopLeft)/2
        var midTop = (startTop-stopTop)/2
        console.log(`stop: left ${stopLeft} top ${stopTop}`)
        console.log(`start: left ${startLeft} top ${startTop}`)
        console.log(`mid: left ${midLeft} top ${midTop}`)
        from.style.left = `${startLeft}px`
        from.style.top = `${startTop}px`
        Velocity(from, { top: midTop, left: midLeft, scale: 2 }, { easing: "easeInSine", duration: 500 })
        Velocity(from, { top: stopTop, left: stopLeft, scale: 1 }, { easing: "easeOutSine", duration: 500, complete: done })
      },
      moveFrom: function (from, to, origin, done) {
        var fromRect = from.getBoundingClientRect()
        console.log(`moveFrom from.top: ${fromRect.top}`)
        console.log(`moveFrom from.left: ${fromRect.left}`)
        var originRect = origin.getBoundingClientRect()
        console.log(`moveFrom origin.top: ${originRect.top}`)
        console.log(`moveFrom origin.left: ${originRect.left}`)
        var toRect = to.getBoundingClientRect()
        console.log(`moveFrom to.top: ${toRect.top}`)
        console.log(`moveFrom to.left: ${toRect.left}`)
        var startTop = fromRect.top - originRect.top
        var startLeft = fromRect.left - originRect.left
        console.log(`moveFrom start.top: ${startTop}`)
        console.log(`moveFrom start.left: ${startLeft}`)
        var stopTop = toRect.top - originRect.top
        var stopLeft = toRect.left - originRect.left
        console.log(`moveFrom stop.top: ${stopTop}`)
        console.log(`moveFrom stop.left: ${stopLeft}`)
        var midLeft = (stopLeft-startLeft)/2
        var midTop = (stopTop-startTop)/2
        console.log(`moveFrom mid.top: ${midTop}`)
        console.log(`moveFrom mid.left: ${midLeft}`)
        from.style.left = `${startLeft}px`
        from.style.top = `${startTop}px`
        Velocity(from, { top: midTop, left: midLeft, scale: 2 }, { easing: "easeInSine", duration: 500 })
        Velocity(from, { top: stopTop, left: stopLeft, scale: 1 }, { easing: "easeOutSine", duration: 500, complete: done })
      }
    },
    computed: {
      showCard: function () {
        var self = this
        return self.value.card.kind != 'none'
      },
      showThief: function () {
        var self = this
        return self.thiefColor != 'none'
      },
      id: function () {
        var self = this
        return `thief-${self.value.row}-${self.value.column}`
      },
      thiefColor: function () {
        var self = this
        if (self.cuLoading) {
          return 'none'
        }
        var cuid = _.get(self.cu, 'id', false)
        if (cuid) {
          var p = self.playerByUID(self.cu.id)
          return _.get(p.colors, self.value.thief.pid - 1, 'none')
        }
        return _.get(self.cp.colors, self.value.thief.pid - 1, 'none')
      },
      animate: {
        get: function () {
          return this.$root.animate
        },
        set: function (value) {
          this.$root.animate = value
        }
      }
    }
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">

  v-card {
    border-color: yellow;
    border-width: 3px;
    border-style: solid;
  }

  .board-space {
    height:90px;
    width:90px;
  }

</style>
