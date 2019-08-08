<template>
  <div>
    <v-tooltip 
      :disabled="tooltip(kind) === '' || !show"
      left
      max-width='150px'
      color='green'
      open-delay='400'
    >
      <sn-card-image
        slot='activator'
        :kind='kind'
        :show='show'
      >
      </sn-card-image>
      <span>{{tooltip(kind)}}</span>
    </v-tooltip>
  </div>
</template>

<script>
  import Image from '@/components/card/Image'
  import Tooltip from '@/components/mixins/Tooltip'

  export default {
    mixins: [ Tooltip ],
    name: 'sn-card',
    props: [ 'kind', 'show' ],
    components: {
      'sn-card-image': Image
    }
  }
</script>
