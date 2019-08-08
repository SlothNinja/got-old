<template>
  <v-btn 
    fab
    :small="size ? size === 'small' : false"
    :medium="size ? size === 'medium' : true"
    :color="color || 'black' "
    :href="showlink"
  >
    <v-avatar
      :size="size && size === 'small' ? '30px' : '48px'"
    >
      <img :src="gravatar(user.emailHash, size, user.gravType)" />
    </v-avatar>
  </v-btn>
</template>

<script>
  import Gravatar from '@/components/mixins/Gravatar'

  export default {
    mixins: [ Gravatar ],
    name: 'sn-user-btn',
    props: [ 'color', 'user', 'size' ],
    computed: {
      showlink: function() {
        var self = this
        return `https://user.slothninja.com/#/show/${self.user.id}`
      }
    }
  }
</script>
