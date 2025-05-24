import { createStore } from 'vuex'

const store = createStore({
  state: {
    user: {
      name: '',
      avatar: ''
    }
  },
  mutations: {
    setUser(state, user) {
      state.user.name = user.name
      state.user.avatar = user.avatar
    }
  },
  actions: {
    login({ commit }, user) {
      // 模拟登录过程
      commit('setUser', user)
    }
  }
})

export default store