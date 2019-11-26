<template>
  <div class="header">
    <div class="control">
      <input type=file style="display: none;" ref="upload-file">
      <Button type="primary" icon="md-cloud-upload" @click="upload">上传</Button>
    </div>
  </div>
</template>

<script>
import client from '../ajax-client'

export default {
  methods: {
    upload() {
      this.$refs['upload-file'].click()
    }
  },
  mounted() {
    this.$refs['upload-file'].onchange = e => {
      client.upload(e.target.files[0])
    }
  }
}
</script>

<style scoped>
.header {
  height: 100%;
  border-bottom: 1px #ddd solid;
  padding: 0 20px;
}
</style>