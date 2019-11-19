import Vue from 'vue'
import Vuex from 'vuex'
import axios from 'axios'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    user: {
      logged: true,
      username: '',
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
    }
  },
  actions: {
    async updateFiles({ state, commit }) {
      let auth = {
        username: state.user.username,
        password: state.user.password
      }
      let resp = await axios.get('/fs/', { auth })
      commit('setFiles', resp.data)
    }
  },
  modules: {
  }
})
