<template>
  <v-navigation-drawer clipped width="150" v-model="drawer" light app >
     <v-list dense>
       <v-list-tile @click="drawer = false; $router.push({ name: 'home' })">
         <v-list-tile-action>
           <v-icon>home</v-icon>
         </v-list-tile-action>
         <v-list-tile-content>
           <v-list-tile-title>Home</v-list-tile-title>
         </v-list-tile-content>
       </v-list-tile>
       <v-list-tile @click="drawer = false; $router.push({ name: 'new' })">
         <v-list-tile-action>
           <v-icon>create</v-icon>
         </v-list-tile-action>
         <v-list-tile-content>
           <v-list-tile-title>Create</v-list-tile-title>
         </v-list-tile-content>
       </v-list-tile>
       <v-list-tile @click="drawer = false; $router.push({ name: 'index', params: {status: 'recruiting'} })">
         <v-list-tile-action>
           <v-icon>playlist_add</v-icon>
         </v-list-tile-action>
         <v-list-tile-content>
           <v-list-tile-title>Join</v-list-tile-title>
         </v-list-tile-content>
       </v-list-tile>
       <v-list-tile @click="drawer = false; $router.push({ name: 'index', params: {status: 'running' } })">
         <v-list-tile-action>
           <v-icon>playlist_play</v-icon>
         </v-list-tile-action>
         <v-list-tile-content>
           <v-list-tile-title>Play</v-list-tile-title>
         </v-list-tile-content>
       </v-list-tile>
     </v-list>
  </v-navigation-drawer>
</template>

<script>
  export default {
    name: 'nav-drawer',
    props: [ 'value' ],
    computed: {
      drawer: {
        get: function () {
          var self = this
          return self.value
        },
        set: function (value) {
          var self = this
          self.$emit('input', value)
        }
      }
    }
  }
</script>
