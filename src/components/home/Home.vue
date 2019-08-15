<template>
  <div>
    <sn-toolbar v-model="nav"></sn-toolbar>
    <sn-snackbar v-model="snackbar.open">
      <div class="text-xs-center"> {{snackbar.message}} </div>
    </sn-snackbar>
    <sn-nav-drawer v-model="nav" app></sn-nav-drawer>
    <v-content></v-content>
    <sn-footer></sn-footer>
  </div>
</template>

<script>
  import Toolbar from '@/components/Toolbar'
  import NavDrawer from '@/components/NavDrawer'
  import Footer from '@/components/Footer'
  import Snackbar from '@/components/Snackbar'

  export default {
    name: 'home',
    components: {
      'sn-snackbar': Snackbar,
      'sn-toolbar': Toolbar,
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
      }
    }
  }
</script>
