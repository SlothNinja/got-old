<template>
  <div class='pa-2 mb-2 sn-entry'>
    <v-system-bar color='green' dark >
      Turn: {{turn}}
    </v-system-bar>
    <div>
      <sn-log-message
        v-for="(entry, index) in value.log"
        :key="index"
        :value='entry'
      >
      </sn-log-message>
    </div>
    <v-divider></v-divider>
    <div class="caption">
      {{createdAt}}
    </div>
  </div>
</template>

<script>
  import Message from '@/components/log/Message'

  const _ = require('lodash')

  export default {
    name: 'sn-log-entry',
    props: [ 'value' ],
    components: {
      'sn-log-message': Message
    },
    computed: {
      turn: function () {
        var self = this
        return _.get(self.value, 'log[0].turn', 0)
      },
      createdAt: function () {
        var self = this
        var d = _.get(self.value, 'header.createdAt', false)
        if (d) {
          return new Date(d).toString()
        }
        return ""
      }
    }
  }
</script>

<style scoped lang="sass">
  ul 
    display: block
    list-style-type: disc
    margin-top: 0
    margin-bottom: 0
    margin-left: 0
    margin-right: 0
    padding-left: 40px

  .sn-entry
    border: 1px solid black
</style>
