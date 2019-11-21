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
      <div class="title">已用 {{ usedText }}</div>
      <Progress status="normal" :percent="usedPercent" :stroke-width="6" hide-info />
    </div>

  </div>
</template>

<script>
import { mapState } from 'vuex'

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
      let usedText = this.humanSize(used)
      let sizeText = size <= 0 ? '无限制' : this.humanSize(size)
      return usedText + ' / ' + sizeText
    },
    ...mapState([
      'info'
    ])
  },
  methods: {
    humanSize(size) {
      const suffix = ['B', 'KB', 'MB', 'GB', 'TB']
      let i
      for (i = 0; size > 1024; ++i) {
        size /= 1024
      }
      if (i >= suffix.length) {
        i = suffix.length - 1
      }
      return size.toFixed(2) + ' ' + suffix[i]
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
  position: relative;
  top: 8px;
  font-size: 12px;
}
</style>