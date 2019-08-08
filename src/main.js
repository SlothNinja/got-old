import Vue from 'vue'
import './plugins/vuetify'
import App from './App.vue'
import router from './router/router'
import axios from 'axios'

Vue.config.productionTip = false

const _ = require('lodash')

new Vue({
  data () {
    return {
      game: {
        header: { title: '', id: 0, turn: 0, phase: 0, colorMaps: [], options: {} },
        state: { glog: [], jewels: {} }
      },
      tab: 'player-1',
      cu: null,
      cuLoading: true,
      idToken: '',
      nav: false,
      animate: true,
      snackbar: { open: false, message: '' },
      tbModel: {},
      log: [],
      cp: {}
    }
  },
  created () {
      var self = this
      self.fetchData()
  },
  methods: {
    fetchData () {
      var self = this
      axios.get('/current')
        .then(function (response) {
          var cu = _.get(response, 'data.cu', false)
          if (cu) {
            self.cu = cu
          }
          self.cuLoading = false
        })
        .catch(function () {
          self.snackbar.message = 'Server Error.  Try again.'
          self.snackbar.open = true
          self.$router.push({ name: 'show', params: { id: self.$route.params.id}})
        })
    },
  },
  router,
  render: h => h(App),
}).$mount('#app')
