import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    info: {
      size: 0,
      used: 0
    },
    user: {
      logged: true,
      username: 'sdjdd',
      password: ''
    },
    files: []
  },
  mutations: {
    user(state, userInfo) {
      state.user.logged = true
      state.user.username = userInfo.username
      state.user.password = userInfo.password
    },
    setFiles(state, files) {
      state.files = files
    },
    setInfo(state, info) {
      state.info.size = info.size
      state.info.used = info.used
    }
  },
  actions: {
  },
  modules: {
  }
})
