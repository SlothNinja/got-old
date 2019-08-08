<template>
  <v-layout row>
  <v-flex xs8 class="text-xs-center">
    <v-tooltip bottom color='green'>
      <v-btn
        slot='activator'
        icon
        :disabled="!canReset"
              @click.native="$emit('action', { action: 'reset' })"
      >
        <v-icon>clear</v-icon>
      </v-btn>
      <span>Reset</span>
    </v-tooltip>
    <v-tooltip bottom color='green'>
      <v-btn
        slot='activator'
        icon
        :disabled="!canUndo"
        @click.native="$emit('action', { action: 'undo' })"
      >
        <v-icon>undo</v-icon>
      </v-btn>
      <span>Undo</span>
    </v-tooltip>
    <v-tooltip bottom color='green'>
      <v-btn
        slot='activator'
        icon
        :disabled="!canRedo"
        @click.native="$emit('action', { action: 'redo' })"
      >
        <v-icon>redo</v-icon>
      </v-btn>
      <span>Redo</span>
    </v-tooltip>
    <v-tooltip bottom color='green'>
      <v-btn
        slot='activator'
        icon
        :disabled="!canFinish"
        @click.native="$emit('action', { action : 'finish' })"
      >
        <v-icon>done</v-icon>
      </v-btn>
      <span>Finish</span>
    </v-tooltip>
  </v-flex>
  </v-layout>
</template>

<script>
  import Player from '@/components/mixins/Player'
  import CurrentUser from '@/components/mixins/CurrentUser'

  var _ = require('lodash')

  export default {
    name: 'sn-controlbar',
    mixins: [ Player, CurrentUser ],
    props: [ 'value' ],
    computed: {
      cp: {
        get: function () {
          return this.$root.cp
        },
        set: function (value) {
          this.$root.cp = value
        }
      },
      canUndo: function () {
        var self = this
        return (self.isCPorAdmin) && (self.value.undoStack.current > self.value.undoStack.committed)
      },
      canRedo: function () {
        var self = this
        return (self.isCPorAdmin) && (self.value.undoStack.current < self.value.undoStack.updated)
      },
      canReset: function () {
        var self = this
        return self.isCPorAdmin
      },
      canFinish: function () {
        var self = this
        return self.isCPorAdmin ? (_.get(self.cp, 'performedAction', true)) : false
      }
    }
  }
</script>
