<template>
  <v-container fluid>
    <v-layout row>
      <v-flex>
        <v-card color="green" dark>
          <v-card-title>
            <span class="title">Chat</span>
            <v-spacer></v-spacer>{{messages.length}} of {{count}}
          </v-card-title>
        </v-card>
      </v-flex>
    </v-layout>
    <v-layout row justify-center>
      <v-flex>
        <v-card>
          <v-card-text>
            <v-container ref="chatbox" id="chatbox" @scroll="handleScroll" >
              <sn-message
                class='pt-2'
                v-for="(message, index) in messages"
                :key="message.id"
                :message="message"
              >
              </sn-message>
            </v-container>
          </v-card-text>
        </v-card>
      </v-flex>
    </v-layout>
    <v-layout row justify-center class="pt-3">
      <v-flex>
        <v-card >
          <v-card-text>
            <v-text-field
              color="green"
              label="Message"
              name="message"
              placeholder="Type Message.  Press 'Send' button."
              textarea
              rows="2"
              required
              v-model="message"
            ></v-text-field>
            <v-flex class="text-xs-center">
              <v-btn color="green" v-if="!(message === '')" dark @click="clear">Clear</v-btn>
              <v-btn color="green" v-if="!(message === '')" dark @click="send">Send</v-btn>
            </v-flex>
          </v-card-text>
        </v-card>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
  import Message from '@/components/chat/Message'

  export default {
    data: function () {
      return {
        count: 0,
        updateScrollHeight: false,
        scrollHeight: 0,
        messages: [],
        message: ''
      }
    },
    components: {
      'sn-message': Message
    },
    name: 'sn-chat-box',
    props: [ 'user' ],
    created () {
      this.fetchData()
    },
    watch: {
      '$route': 'fetchData'
    },
    updated: function () {
      var self = this
      if (self.updateScrollHeight) {
        self.scrollToHeight()
      } else {
        self.scroll()
      }
    },
    methods: {
      fetchData () {
        var self = this
        var offset = self.messages.length
        $.ajax({
          url: 'http://localhost:8081/got/game/messages/' + self.$route.params.id + '/' + offset,
          dataType: 'json',
          success: function (data) {
            if (data.messages) {
              self.messages = data.messages.concat(self.messages)
              self.count = data.count
            }
          }
        })
      },
      clear: function () {
        var self = this
        self.message = ''
      },
      send: function () {
        var self = this
        var obj = {
          message: self.message,
          creator: self.user
        }

        $.ajax({
          url: 'http://localhost:8081/got/game/add-message/' + self.$route.params.id,
          dataType: 'json',
          data: JSON.stringify(obj),
          contentType: 'application/json; charset=utf-8',
          type: 'PUT',
          success: function (data) {
            self.add(data)
            self.clear()
          }
        })
      },
      add: function (message) {
        var self = this
        self.updateScrollHeight = false
        self.messages.push(message)
        self.count += 1
      },
      scroll: function () {
        var self = this
        var chatbox = self.$refs.chatbox
        self.scrollHeight = chatbox.scrollHeight
        chatbox.scrollTop = self.scrollHeight
      },
      scrollToHeight: function () {
        var self = this
        var chatbox = self.$refs.chatbox
        chatbox.scrollTop = chatbox.scrollHeight - self.scrollHeight
        self.scrollHeight = chatbox.scrollHeight
      },
      handleScroll: function () {
        var self = this
        var chatbox = self.$refs.chatbox
        if ((chatbox.scrollTop === 0) && (self.messages.length < self.count)) {
          self.updateScrollHeight = true
          self.fetchData()
        }
      }
    }
  }
</script>

<style scoped lang="scss">

  #chatbox {
    border:2px solid black;
    height:400px;
    overflow:scroll;
  }
</style>
