<script>
  export default {
    computed: {
      cu: {
        get: function () {
          return this.$root.cu
        },
        set: function (value) {
          this.$root.cu = value
        }
      },
      cuLoading: function () {
        return this.$root.cuLoading
      }
    }
  }
</script>
