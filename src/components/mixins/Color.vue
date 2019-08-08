<script>
  var _ = require('lodash')

  export default {
    methods: {
      colorByPID: function (colorMap, pid) {
        var colorStrings = [ 'none', 'red', 'yellow', 'purple', 'black', 'brown', 'white', 'green', 'orange' ]
        var index = pid - 1
        return _.get(colorStrings, _.get(colorMap, index, 0), 'none')
      }
//      },
//      thiefColor: function (colorMap, pid) {
//        var self = this
//        return self.colorByPID(colorMap, pid)
//      }
    }
  }
</script>
