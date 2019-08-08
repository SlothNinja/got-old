<template>
  <div>
    <sn-toolbar v-model="nav"></sn-toolbar>
    <sn-snackbar v-model="snackbar.open">
      <div class="text-xs-center">
        {{snackbar.message}}
      </div>
    </sn-snackbar>
    <sn-nav-drawer v-model="nav" app></sn-nav-drawer>
    <v-content></v-content>
    <sn-footer></sn-footer>
  </div>
</template>

<script>
  import Toolbar from '@/components/Toolbar'
  // import GameToolbar from '@/components/game/Toolbar'
  import NavDrawer from '@/components/NavDrawer'
  import Footer from '@/components/Footer'
  import Snackbar from '@/components/Snackbar'

  export default {
    name: 'home',
    data () {
      return {
        u: { name: '', emailReminders: 'true', emailNotifications: 'true', gravType: 'monsterid' },
        auto: true,
        idToken: '',
        loading: true
      }
    },
    components: {
      'sn-snackbar': Snackbar,
      'sn-toolbar': Toolbar,
    //  'sn-game-toolbar': GameToolbar,
      'sn-nav-drawer': NavDrawer,
      'sn-footer': Footer
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
      nav: {
        get: function () {
          return this.$root.nav
        },
        set: function (value) {
          this.$root.nav = value
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
      tbComponent: function () {
        if (this.$route.name == "game") {
          return 'sn-game-toolbar'
        }
        return 'sn-toolbar'
      }
    }
  }
</script>
