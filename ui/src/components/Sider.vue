<template>
  <div>
    <div class="logo">Banana</div>
    <Menu active-name="1-1" width="auto">
      <MenuItem name="1-1">
        <Icon type="ios-navigate"></Icon>
        <span>全部</span>
      </MenuItem>
    </Menu>
    <div class="used-info">
      <div class="title">
        <span>已用 {{ usedText }}</span>
        <span class="refresh" @click="refreshUsed">刷新</span>
      </div>
      <Progress status="normal" :percent="usedPercent" :stroke-width="6" hide-info />
    </div>

  </div>
</template>

<script>
import client from '../ajax-client'
import { mapState } from 'vuex'
import { humanSize } from '../utils'

export default {
  computed: {
    usedPercent() {
      if (this.info.size <= 0) {
        return 0
      }
      return this.info.used / this.info.size * 100
    },
    usedText() {
      let { used, size } = this.info
      let usedText = humanSize(used)
      let sizeText = size <= 0 ? '无限制' : humanSize(size)
      return usedText + ' / ' + sizeText
    },
    ...mapState([
      'info'
    ])
  },
  methods: {
    refreshUsed() {
      client.info()
    }
  }
}
</script>

<style scoped>
.used-info {
  position: absolute;
  bottom: 10px;
  padding: 0 10px;
  width: 100%;
}
.used-info .title {
  cursor: default;
  position: relative;
  top: 8px;
  font-size: 12px;
}
.used-info .refresh {
  display: inline-block;
  position: absolute;
  right: 0;
}
.used-info .refresh:hover {
  cursor: pointer;
  text-decoration: underline;
}
</style>