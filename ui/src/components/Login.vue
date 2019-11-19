<template>
  <div>
    <div>
      <input type="text" placeholder="用户名" v-model="username">
    </div>
    <div>
      <input type="password" placeholder="密码" v-model="password">
    </div>
    <div>
      <button @click="login">登录</button>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      username: '',
      password: ''
    }
  },
  methods: {
    login() {
      this.$store.commit('user', {
        username: this.username,
        password: this.password
      })
      this.$router.push({
        name: 'home'
      })
    }
  }
}
</script>

<style>

</style>