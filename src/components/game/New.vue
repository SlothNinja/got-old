<template>
  <div>
    <sn-toolbar v-model="nav"></sn-toolbar>
    <sn-snackbar v-model="snackbar.open">
      <div class="text-xs-center">
        {{snackbar.message}}
      </div>
    </sn-snackbar>
    <sn-nav-drawer v-model="nav" app></sn-nav-drawer>
    <v-content>  
      <v-container grid-list-md >
        <v-layout row wrap>
          <v-flex xs6>
            <v-card height="31em">
              <v-card-title primary-title>
                <h3>New Game</h3>
              </v-card-title>
              <v-card-text>
                <v-form action="/got/game" method="post">
                  <v-text-field
                    name="title"
                    label="Title"
                    v-model="game.header.title"
                    id="title"
                  >
                  </v-text-field>
                  <v-select
                    label="Number Players"
                    name="num-players"
                    v-bind:items="npItems"
                    v-model="game.header.numPlayers"
                  >
                  </v-select> 
                  <v-select 
                    id="two-thief-variant"
                    label="Two Thief Variant"
                    name="two-thief-variant"
                    v-bind:items="optItems"
                    v-model="game.header.twoThief"
                  >
                  </v-select> 
                  <v-text-field
                    label="Password"
                    name="password"
                    id="password"
                    v-model="game.header.password"
                    placeholder="Enter Password for Private Game"
                    type="password"
                    autocomplete="new-password"
                  >
                  </v-text-field>
                  <v-btn color='green' dark @click="putData">Submit</v-btn>
                </v-form>
              </v-card-text>
            </v-card>
          </v-flex>
          <v-flex xs6>
            <v-card height="31em">
              <v-img height='200px' :src="require('../../assets/got-box.jpg')" />
              <v-card-text>
                <v-layout row>
                  <v-flex xs5>Designer</v-flex>
                  <v-flex>Adam E. Daulton</v-flex>
                </v-layout>
                <v-layout row>
                  <v-flex xs5>Artists</v-flex> 
                  <v-flex>Jeremy Montz</v-flex> 
                </v-layout>
                <v-layout row> 
                  <v-flex xs5>Publisher</v-flex> 
                  <v-flex><a href="http://www.thegamecrafter.com/">The Game Crafter, LLC</a></v-flex>
                </v-layout>
                <v-layout row>
                  <v-flex xs5>Year Published</v-flex>
                  <v-flex>2012</v-flex>
                </v-layout>
                <v-layout row> 
                  <v-flex xs5>On-Line Developer</v-flex> 
                  <v-flex>Jeff Huter</v-flex> 
                </v-layout> 
                <v-layout row> 
                  <v-flex xs5>Permission Provided By</v-flex> 
                  <v-flex>Adam E Daulton</v-flex> 
                </v-layout> 
                <v-layout row> 
                  <v-flex xs5>Rules (pdf)</v-flex> 
                  <v-flex><a href="/static/rules/got.pdf">Guild Of Thieves (English)</a></v-flex> 
                </v-layout> 
              </v-card-text>
            </v-card>
          </v-flex>
        </v-layout>
      </v-container>
    </v-content>
    <sn-footer></sn-footer>
  </div>
</template>

<script>
  import Controlbar from '@/components/game/Controlbar'
  import Toolbar from '@/components/Toolbar'
  import Snackbar from '@/components/Snackbar'
  import Footer from '@/components/Footer'
  import NavDrawer from '@/components/NavDrawer'
  import RDrawer from '@/components/rdrawer/Drawer'
  import Board from '@/components/board/Board'
  import Bar from '@/components/card/Bar'
  import StatusPanel from '@/components/game/StatusPanel'
  import Panels from '@/components/player/Panels'
  import Messagebar from '@/components/game/Messagebar'
  import ChatBox from '@/components/chat/Box'
  import GameLog from '@/components/log/Box'
  import Thief from '@/components/thief/Image'
  import Player from '@/components/mixins/Player'
  import CurrentUser from '@/components/mixins/CurrentUser'
  import Game from '@/components/mixins/Game'
  import CardImage from '@/components/card/Image'
  import Color from '@/components/mixins/Color'


  const _ = require('lodash')
  const axios = require('axios')

  export default {
    name: 'newGame',
    data () {
      return {
        game: {
          header: { title: '', id: 0, turn: 0, phase: 0, colorMaps: [], options: {} },
          state: { glog: [], jewels: {} }
        },
        path: '/game/new',
        nav: false,
        npItems: [
          { text: '2', value: 2 },
          { text: '3', value: 3 },
          { text: '4', value: 4 }
        ],
        optItems: [
          { text: 'Yes', value: true },
          { text: 'No', value: false }
        ]
      }
    },
    components: {
      'sn-control-bar': Controlbar,
      'sn-toolbar': Toolbar,
      'sn-snackbar': Snackbar,
      'sn-nav-drawer': NavDrawer,
      'sn-rdrawer': RDrawer,
      'sn-board': Board,
      'sn-card-bar': Bar,
      'sn-status-panel': StatusPanel,
      'sn-player-panels': Panels,
      'sn-chat-box': ChatBox,
      'sn-game-log': GameLog,
      'sn-messagebar': Messagebar,
      'sn-card-image': CardImage,
      'sn-thief-image': Thief,
      'sn-footer': Footer
    },
    created () {
      var self = this
      self.fetchData()
    },
    watch: {
      '$route': 'fetchData'
    },
    computed: {
      cu: {
        get: function () {
          return this.$root.cu
        },
        set: function (value) {
          this.$root.cu = value
        }
      },
      snackbar: {
        get: function () {
          return this.$root.snackbar
        },
        set: function (value) {
          this.$root.snackbar = value
        }
      },
    },
    methods: {
      fetchData () {
        var self = this
        axios.get(self.path)
          .then(function (response) {
            var msg = _.get(response, 'data.msg', false)
            if (msg) {
              self.snackbar.message = msg
              self.snackbar.open = true
            }
            var header = _.get(response, 'data.header', false)
            if (header) {
              self.game.header = header
            }
            var state = _.get(response, 'data.state', false)
            if (header) {
              self.game.state = state
            }
            self.loading = false
          })
          .catch(function () {
            self.loading = false
            self.snackbar.message = 'Server Error.  Try refreshing page.'
            self.snackbar.open = true
        })
      },
      putData () {
        var self = this
        self.loading = true
        console.log(`sent: ${JSON.stringify(self.game)}`)
        axios.put(self.path, self.game)
          .then(function (response) {
            console.log(`response: ${JSON.stringify(response)}`)
            self.loading = false
            self.$router.push({ name: 'index', params: { status: 'recruiting' }})
          })
          .catch(function () {
            self.loading = false
            self.snackbar.message = 'Server Error.  Try again.'
            self.snackbar.open = true
            self.$router.push({ name: 'home'})
          })
        }
      }
    }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
  h1, h2, h3 {
    font-weight: normal;
  }
</style>
