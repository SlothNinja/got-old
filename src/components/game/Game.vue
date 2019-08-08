<template>
  <div>
    <sn-toolbar v-model="nav">
      <v-layout row>
        <v-flex xs3>
          <v-tooltip bottom color='green'>
            <v-btn slot='activator' icon @click.stop="history = !history">
              <v-icon>history</v-icon>
            </v-btn>
            <span>History</span>
          </v-tooltip>
          <v-tooltip bottom color='green'>
            <v-btn slot='activator' icon @click.stop="chat = !chat">
              <v-icon>chat</v-icon>
            </v-btn>
            <span>Chat</span>
          </v-tooltip>
        </v-flex>
        <v-flex xs8>
        </v-flex>
        <v-flex xs8>
          <sn-control-bar
            v-model="game"
            @action="action($event)"
          >
          </sn-control-bar>
        </v-flex>
      </v-layout>
    </sn-toolbar>
    <sn-snackbar v-model="snackbar.open">
      <div class="text-xs-center">
        {{snackbar.message}}
      </div>
    </sn-snackbar>
    <sn-nav-drawer v-model="nav" app></sn-nav-drawer>
    <sn-rdrawer v-model="history" app>
      <sn-game-log v-if="history" :count="current" :open="history"></sn-game-log>
    </sn-rdrawer>
    <sn-rdrawer v-model="chat" app>
      <sn-chat-box v-if="chat" :user="cu"></sn-chat-box>
    </sn-rdrawer>
    <v-content>
    <v-container id="game" grid-list-md >
      <v-layout row>
        <v-flex xs3>
          <sn-status-panel class="mt-1" :game="game"></sn-status-panel>
          <sn-player-panels
            v-model="tab"
            @show="cardbar = $event"
            @pass="action({action: 'pass', data: {} })"
            :players="game.state.players"
          >
          </sn-player-panels>
        </v-flex>
        <v-flex xs9>
          <sn-messagebar class="mt-1 mb-2">
            {{message}}
          </sn-messagebar>
          <sn-board
            id="board"
            :value="game.state.grid"
            @selected="selected($event)"
          >
          </sn-board>
        </v-flex>
      </v-layout>
      <v-layout row>
        <sn-card-bar
          v-if="selectedPlayer"
          :phase="game.header.phase"
          :player="selectedPlayer"
          v-model="cardbar"
          @selected="selected($event)"
        >
        </sn-card-bar>
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
    mixins: [ Game, CurrentUser, Player, Color ],
    name: 'game',
    data () {
      return {
        path: '/game',
        cardbar: false,
        nav: false,
        history: false,
        chat: false
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
      'sn-footer': Footer,
    },
    created () {
      var self = this
      self.fetchData()
    },
    watch: {
      '$route': 'fetchData'
    },
    methods: {
      myUpdate: function (data) {
        var self = this
        if (_.has(data, 'game')) {
          self.game = data.game
          document.title = data.game.header.title + ' #' + data.game.header.id
          // self.cp = self.currentPlayer(data.game.header, data.game.state.players)
        }
        var msg = _.get(data, 'message', '')
        if (msg !== '') {
          self.snackbar.message = msg
          self.snackbar.open = true
        }
        if (_.has(data, 'cu')) {
          self.cu = data.cu
        }
        self.cardbar = false
      },
      fetchData: function () {
        var self = this
        axios.get(`${self.path}/show/${self.$route.params.id}`)
          .then(function (response) {
            var msg = _.get(response, 'data.message', false)
            if (msg) {
              self.snackbar.message = msg
              self.snackbar.open = true
            }
            var game = _.get(response, 'data.game', false)
            if (game) {
              self.game = game
              var tabid = 1
              var cuid = _.get(self.cu, 'id', false)
              if (cuid) {
                self.tab = `player-${self.pidByUID(cuid)}`
              } 
            }
            self.loading = false
          })
          .catch(function (error) {
            self.loading = false
            self.snackbar.message = 'Server Error.  Try refreshing page.'
            self.snackbar.open = true
        })
        // $.ajax({
        //   url: 'http://localhost:8081/got/game/show/' + self.$route.params.id,
        //   dataType: 'json',
        //   success: function (data) {
        //     self.myUpdate(data)
        //     self.tab = `player-${self.pidByUID(self.game.state.players, self.cu.id)}`
        //   }
        // })
      },
      selected: function (data) {
        var self = this
        console.log(`selected: ${JSON.stringify(data)}`)
        console.log(`phase: ${JSON.stringify(self.game.header.phase)}`)
        switch (self.game.header.phase) {
          case 3:  // place-thief
            self.action({
              action: 'place-thief',
              data: {
                row: data.row,
                column: data.column
              }
            })
            break
          case 4: // play-card
            self.action({
              action: 'play-card',
              data: {
                kind: data
              }
            })
            break
          case 5: // select-thief
            self.action({
              action: 'select-thief',
              data: {
                row: data.row,
                column: data.column
              }
            })
            break
          case 6: // move-thief
            self.action({
              action: 'move-thief',
              data: {
                row: data.row,
                column: data.column
              }
            })
            break
        }
      },
      action: function (data) {
        var self = this
        var action = data.action
        self.tab = `player-${self.pidByUID(self.cu.id)}`
        axios.put(`${self.path}/${action}/${self.$route.params.id}`, data.data)
          .then(function (response) {
            console.log(`response: ${JSON.stringify(response)}`)
            var game = _.get(response, 'data.game', false)
            if (game) {
              self.game = game
              // var cp = self.currentPlayer(game.header, game.state.players)
              // console.log(`game cp: ${JSON.stringify(cp)}`)
              // self.cp = cp
            }
            var data = _.get(response, 'data', false)
            console.log(`data: ${JSON.stringify(data)}`)
            var msg = _.get(response, 'data.message', false)
            console.log(`message: ${JSON.stringify(msg)}`)
            if (msg) {
              self.loading = false
              self.snackbar.message = msg
              self.snackbar.open = true
            }
                  
            console.log(`action: ${action}`)
            if (action == 'undo' || action == 'reset') {
              self.animate = false
            } else {
              self.animate = true
            }
            var animations = _.get(data, 'game.state.animations', [])
            var play = _.get(data, 'animate', false)
            if (play && (_.size(animations) > 0)) {
              self.animations(data.game.state.animations, function () { self.myUpdate(data) })
            } else {
              self.myUpdate(data)
            }
          })
          .catch(function () {
            self.loading = false
            self.snackbar.message = 'Server Error.  Try refreshing page.'
            self.snackbar.open = true
        })
//        $.ajax({
//          url: 'http://localhost:8081/got/game/' + data.action + '/' + self.$route.params.id,
//          dataType: 'json',
//          data: JSON.stringify(data.data),
//          contentType: 'application/json; charset=utf-8',
//          type: 'PUT',
//          success: function (data) {
//            console.log('data: ' + JSON.stringify(data))
//            var animations = _.get(data, 'game.state.animations', [])
//            var play = _.get(data, 'animate', false)
//            if (play && (_.size(animations) > 0)) {
//              self.animations(data.game.state.animations, function () { self.myUpdate(data) })
//            } else {
//              self.myUpdate(data)
//            }
//          }
//        })
      },
      animations: function (animations, completed) {
        var self = this
        var head = _.head(animations)
        var tail = _.tail(animations)
        switch (head.name) {
          case 'place-thief':
            self.animatePlaceThief(head, tail, completed)
            break
          case 'play-card':
            self.animatePlayCard(head, tail, completed)
            break
          case 'move-thief':
            self.animateMoveThief(head, tail, completed)
            break
          case 'claim-item':
            self.animateClaimItem(head, tail, completed)
            break
          case 'bumped-thief':
            self.animateBumpedThief(head, tail, completed)
            break
          case 'draw-card':
            self.animateDrawCard(head, tail, completed)
            break
          case 'shuffle-card':
            self.animateShuffle(head, tail, completed)
            break
          default:
            console.log(`unrecognized animation ${head.name} received`)
        }
      },
      completeAnimation: function (animations, completed) {
        var self = this
        if (_.size(animations) === 0) {
          completed()
        } else {
          self.animations(animations, completed)
        }
      },
      animatePlaceThief: function (animation, animations, completed) {
        // var self = this
        // var data = animation.data
        // var thief = $('#movable-thief')
        //   .attr('class', `thief-${data.player.color} thief`)
        //   .offset($(`#${self.tab} .avatar`).offset())
        // var to = $(`#space-${data.area.row}-${data.area.column} .thief`)
        // self.animateMoveTo(thief, to.offset(), function () {
        //   self.completeAnimation(animations, function () {
        //     completed()
        //     $('#movable-thief')
        //       .attr('class', 'thief-none thief')
        //   })
        // })
      },
      animatePlayCard: function (animation, animations, completed) {
        // var self = this
        // var data = animation.data
        // var card = $('#playable-card')
        //   .attr('class', `got-card ${data.card.kind}`)
        // card.offset($(`#cardbar .${data.card.kind}`).offset())
        // self.cardbar = false
        // self.animateMoveTo(card, $(`#discard-${self.cp.id} .got-card`).offset(), function () {
        //   self.completeAnimation(animations, function () {
        //     completed()
        //     $('#playable-card')
        //       .attr('class', 'got-card none')
        //   })
        // })
      },
      animateMoveThief: function (animation, animations, completed) {
        // var self = this
        // var data = animation.data
        // var to = data.to
        // var thief = $(`#space-${to.row}-${to.column} .thief`)
        // var from = data.from
        // var area = self.game.state.grid[from.row - 1][from.column - 1]
        // var newThief = $('#movable-thief')
        //   .attr('class', `thief-${data.color} thief`)
        //   .offset($(`#space-${from.row}-${from.column} .thief`).offset())
        // area.thief.pid = 0
        // area.thief.color = 'none'
        // self.animateMoveTo(newThief, thief.offset(), function () {
        //   self.completeAnimation(animations, function () {
        //     completed()
        //     $('#movable-thief')
        //       .attr('class', 'thief-none thief')
        //   })
        // })
      },
      animateBumpedThief: function (animation, animations, completed) {
        // var self = this
        // var data = animation.data
        // var to = data.to
        // var thief = $(`#space-${to.row}-${to.column} .thief`)
        // var from = data.from
        // var area = self.game.state.grid[from.row - 1][from.column - 1]
        // var newThief = $('#bumpable-thief')
        //   .attr('class', `thief-${data.color} thief`)
        //   .offset($(`#space-${from.row}-${from.column} .thief`).offset())
        // area.thief.pid = 0
        // area.thief.color = 'none'
        // self.animateMoveTo(newThief, thief.offset(), function () {
        //   self.completeAnimation(animations, function () {
        //     completed()
        //     $('#bumpable-thief')
        //       .attr('class', `thief-none thief`)
        //   })
        // })
      },
      animateClaimItem: function (animation, animations, completed) {
        console.log('animation.data: ' + JSON.stringify(animation.data))
        // var self = this
        // var data = animation.data
        // var pid = data.player.id
        // var to
        // if (data.toHand) {
        //   to = $(`#hand-${pid} .got-card`)
        // } else {
        //   to = $(`#discard-${pid} .got-card`)
        // }
        // var from = data.from
        // var area = self.game.state.grid[from.row - 1][from.column - 1]
        // area.card.kind = 'none'
        // var newCard = $('#claimable-card')
        //   .attr('class', `got-card ${data.card.kind}`)
        //   .offset($(`#space-${from.row}-${from.column} .thief`).offset())
        // self.animateMoveTo(newCard, to.offset(), function () {
        //   self.completeAnimation(animations, function () {
        //     completed()
        //     $('#claimable-card')
        //       .attr('class', 'got-card none')
        //     $('#movable-thief')
        //       .attr('class', 'thief-none thief')
        //   })
        // })
      },
      animateDrawCard: function (data, animations, completed) {
        // var self = this
        // var card = $('#drawable-card')
        // card.offset($(`#draw-${data.pid} .got-card`).offset())
        //   .attr('class', 'got-card card-back')
        // self.animateMoveTo(card, $(`#hand-${data.pid} .got-card`).offset(), function () {
        //   self.completeAnimation(animations, function () {
        //     completed()
        //     $('#drawable-card')
        //       .attr('class', 'got-card none')
        //   })
        // })
      },
      animateShuffle: function (data, animations, completed) {
        //  var self = this
        // var card = $('#shufflable-card')
        // card.offset($(`#discard-${data.pid} .got-card`).offset())
        //   .attr('class', 'got-card card-back')
        // $('#claimable-card').attr('class', 'got-card card-back')
        // self.animateMoveTo(card, $(`#draw-${data.pid} .got-card`).offset(), function () {
        //   self.completeAnimation(animations, function () {
        //     completed()
        //     $('#shufflable-card')
        //       .attr('class', 'got-card none')
        //   })
        // })
      },
      animateMoveTo: function (obj, to, complete) {
        var height = obj.height()
        var width = obj.width()
        var from = obj.offset()
        var midpoint = {
          top: (from.top + to.top) / 2,
          left: (from.left + to.left) / 2
        }
        obj.velocity({
          top: midpoint.top,
          left: midpoint.left,
          height: height * 2,
          width: width * 2
        }, { duration: 200 })
        .velocity({
          top: to.top,
          left: to.left,
          height: height,
          width: width
        }, {
          duration: 200,
          complete: function () {
            if (complete) {
              complete()
            }
          }
        })
      }
    },
    computed: {
      animate: {
        get: function () {
          return this.$root.animate
        },
        set: function (value) {
          this.$root.animate = value
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
      tbModel: {
        get: function () {
          return this.$root.tbModel
        },
        set: function (value) {
          this.$root.tbModel = value
        }
      },
      tbComponent: {
        get: function () {
          return this.$root.tbComponent
        },
        set: function (value) {
          this.$root.tbComponent = value
        }
      },
      cu: {
        get: function () {
          return this.$root.cu
        },
        set: function (value) {
          this.$root.cu = value
        }
      },
      selectedPlayer: function () {
        var self = this
        console.log('selectedPlayer')
        return self.playerByPID(self.cardbar)
      },
      current: function () {
        var self = this
        return _.get(self, 'game.undoStack.current', 0)
      },
      message: function () {
        var self = this
        if (!self.isCP) {
          var name = _.get(self.cp, 'user.name', 'current player')
          return 'Please wait for ' + name + ' to take a turn.'
        }
        switch (self.game.header.phase) {
          case 0: { // None
            break
          }
          case 1: { // Setup
            break
          }
          case 2: { // Start Game
            break
          }
          case 3: { // Place Thieves
            if (self.cp.performedAction) {
              return 'Finish Turn.'
            } else {
              return 'Select empty space in grid to place thief.'
            }
          }
          case 4: { // Play Card
            if (!self.cardbar) {
              return 'Select hand or pass'
            } else {
              return 'Select card from hand'
            }
          }
          case 5: { // Select Thief
            return 'Select thief in grid'
          }
          case 6: { // Move Thief
            return 'Select highlighted spot in grid to move thief'
          }
          case 7: { // Claim Magical Item
            return 'Finish turn by selecting above check mark.'
          }
        }
      }
    }
  }
</script>
