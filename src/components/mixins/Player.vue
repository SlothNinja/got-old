<script>
  import Game from '@/components/mixins/Game'

  const _ = require('lodash')

  export default {
    mixins: [ Game ],
    computed: {
      cp: function () {
        var self = this
        var cpid = _.get(self.game.header.cpUserIndices, 0, -1)
        self.$root.cp = self.playerByPID(cpid)
        return self.$root.cp
      },
      isCP: function () {
        var self = this
        return self.isPlayerFor(self.cp, self.$root.cu)
      },
       isCPorAdmin: function () {
         var self = this
         return (self.$root.cu && self.$root.cu.admin) || self.isCP
       },
       tab: {
         get: function () {
           return this.$root.tab
         },
         set: function (value) {
           this.$root.tab = value
         }
       },
    },
    methods: {
      playerByPID: function (pid) {
        var self = this
        return _.find(self.game.state.players, ['id', pid])
      },
      playersByPIDS: function (pids) {
        var self = this
        return _.map(pids, function (pid) {
          return self.playerByPID(pid)
        })
      },
      playerByUID: function (uid) {
        var self = this
        return _.find(self.game.state.players, ['user.id', uid])
      },
      pidByUID: function (uid) {
        var self = this
        return _.get(self.playerByUID(uid), 'id', -1)
      },
      isPlayerFor: function (player, user) {
        var admin = _.get(user, 'admin', false)
        if (admin) {
          return true
        }
        var pid = _.get(player, 'user.id', -1)
        var uid = _.get(user, 'id', -2)
        return pid === uid
      }
    }
  }
</script>
