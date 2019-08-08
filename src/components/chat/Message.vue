<template>
  <div class="mb-2">
    <v-toolbar dense color='green' dark flat >
      <sn-user-btn :user="message.creator" size="small" ></sn-user-btn>
      {{message.creator.name}}
      <v-spacer></v-spacer>
      <v-tooltip bottom color='green'>
        <v-btn slot='activator' icon>
          <v-icon>help</v-icon>
        </v-btn>
        <span>Help</span>
      </v-tooltip>
    </v-toolbar>
    <div class="ma-2">
      {{message.message}}
    </div>
    <v-divider></v-divider>
    <div class="caption">
      Created: {{Date(message.createdAt)}}
    </div>
  </div>
</template>

<script>
  import UserButton from '@/components/user/Button'
  import Color from '@/components/mixins/Color'
  // import { mapGetters } from 'vuex'

  export default {
    mixins: [ Color ],
    name: 'sn-message',
    props: [ 'message' ],
    components: {
      'sn-user-btn': UserButton
    },
    computed: {
      // ...mapGetters([
      //   'pidByUID'
      // ]),
      pid: function () {
        var self = this
        var pid = self.pidByUID(self.message.creatorID)
        return pid
      }
    }
  }
</script>
