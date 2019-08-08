<template>
  <v-container fluid>
    <v-layout row>
      <v-flex>
        <v-card color="green" dark>
          <v-card-title>
            <span class="title">Log</span>
            <v-spacer></v-spacer>{{log.length}} of {{count}}
          </v-card-title>
        </v-card>
      </v-flex>
    </v-layout>
    <v-layout row justify-center>
      <v-flex>
        <v-card>
          <v-card-text>
            <v-container style="border:2px solid black;height:600px;overflow:scroll">
              <sn-log-entry
                class='pt-2'
                v-for="(entry, index) in log"
                :key="index"
                :value='entry'
              >
              </sn-log-entry>
            </v-container>
          </v-card-text>
        </v-card>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
  import Entry from '@/components/log/Entry'

  const _ = require('lodash')
  const axios = require('axios')

  export default {
    data: function () {
      return {
        path: 'game/glog',
        offset: 0,
        loaded: false
      }
    },
    props: [ 'count', 'open' ],
    components: {
      'sn-log-entry': Entry
    },
    name: 'sn-game-log',
    created () {
      var self = this
      self.fetchData()
    },
    watch: {
      count: function (newCount, oldCount) {
        var self = this
        if (self.loaded) {
          self.loaded = false
        }
      },
      open: function (oldValue, newValue) {
        var self = this
        if ((self.open) && (!self.loaded)) {
          self.fetchData()
        }
      }
    },
    computed: {
      snackbar: {
        get: function () {
          return this.$root.snackbar
        },
        set: function (value) {
          this.$root.snackbar = value
        }
      },
      log: {
        get: function () {
          return this.$root.log
        },
        set: function (value) {
          this.$root.log = value
        }
      }
    },
    methods: {
      fetchData: _.debounce(
        function () {
          var self = this
          axios.get(`${self.path}/${self.$route.params.id}/${self.count}/${self.offset}`)
            .then(function (response) {
              console.log(`response: ${JSON.stringify(response)}`)
              var msg = _.get(response, 'data.message', false)
              if (msg) {
                self.snackbar.message = msg
                self.snackbar.open = true
              }
              var offset = _.get(response, 'data.offset', false)
              if (offset) {
                self.offset = offset
              }
              var logs = _.get(response, 'data.logs', false)
              if (logs) {
                console.log(`logs: ${JSON.stringify(logs)}`)
                var flogs = _.filter(logs, function(log) { return !(_.isNull(_.get(log, 'entries'))) })
                self.log = flogs
              }
              self.loaded = true
            })
            .catch(function () {
              self.loading = false
              self.snackbar.message = 'Server Error.  Try refreshing page.'
              self.snackbar.open = true
          })
        },
        500
      )
    }
  }
</script>
