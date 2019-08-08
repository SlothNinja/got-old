<script>
  var _ = require('lodash')

  export default {
    methods: {
      gravatar: function (hash, size, t) {
        var self = this
        var sz = _.get(self.gravSizes(), size, '48')

        if (!(_.includes(this.gravTypes(), t))) {
          t = 'monsterid'
        }
        if (t == 'personal') {
          return `https://www.gravatar.com/avatar/${hash}?s=${sz}`
        }
        return `https://www.gravatar.com/avatar/${hash}?s=${sz}&d=${t}&f=y`
      },
      gravTypes: function () {
        return ['personal', 'identicon', 'monsterid', 'retro', 'robohash']
      },
      gravSizes: function () {
        return { 'small': '30', 'medium': '48' }
      }
    }
  }
</script>
