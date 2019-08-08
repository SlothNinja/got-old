<template>
  <v-card
    color="green" >
    <v-card-text>
      <v-layout row >
        <div class="col-pad"></div>
        <div 
          v-for="i in numCols"
          :key="i"
          class="col-text"
        >{{i}}</div>
      </v-layout>
      <v-layout 
        row 
        v-for="(row, rindex) in value"
        :key="rindex"
      >
        <div class="row-text">
          {{rowText(row[rindex].row)}}
        </div>
        <sn-board-space
          v-for="(cell, cindex) in row"
          :id="'space-' + cell.row + '-' + cell.column"
          :key="`space-${cell.row}-${cell.column}`"
          :value="cell"
          @selected="$emit('selected', cell)"
        >
        </sn-board-space>
        <div class="row-text">
          {{rowText(row[rindex].row)}}
        </div>
      </v-layout>
      <v-layout row >
        <div class="col-pad"></div>
        <div 
          v-for="i in numCols"
          :key="i"
          class="col-text"
        >{{i}}</div>
      </v-layout>
    </v-card-text>
  </v-card>
</template>

<script>
  import Space from '@/components/board/Space'
  import Text from '@/components/mixins/Text'

  var _ = require('lodash')

  export default {
    mixins: [ Text ],
    name: 'sn-board',
    props: [ 'value' ],
    components: {
      'sn-board-space': Space
    },
    computed: {
      numCols: function () {
        var self = this
        return _.get(self.value, 0, []).length
      }
    }
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
  .col-text, .row-text {
    color:white;
    font-weight:bold
  }

  .row-text {
    position:relative;
    text-align:center;
    width:20px;
    top:40px;
  }

  .col-text {
    position:relative;
    text-align:center;
    width:90px;
    margin:4px;
  }

  .col-pad {
    width:20px;
  }

</style>
